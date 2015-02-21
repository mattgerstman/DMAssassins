package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"fmt"
	"github.com/getsentry/raven-go"
	"strconv"
	"time"
)

const (
	SecInHour        = 3600
	TwentyFourHours  = 86400
	FourtyEightHours = 172800
)

type KillTimer struct {
	GameId    uuid.UUID
	CreateTs  int64
	ExecuteTs int64
}

var activeTimers map[string]*time.Timer

func getActiveTimers() map[string]*time.Timer {
	if activeTimers == nil {
		activeTimers = make(map[string]*time.Timer)
	}
	return activeTimers
}

func addActiveTimer(gameId uuid.UUID, timer *time.Timer) {
	if activeTimers == nil {
		activeTimers = make(map[string]*time.Timer)
	}
	activeTimers[gameId.String()] = timer
}

func getActiveTimer(gameId uuid.UUID) (timer *time.Timer) {
	if activeTimers == nil {
		return nil
	}
	return activeTimers[gameId.String()]
}

func stopTimer(gameId uuid.UUID) {
	if activeTimers == nil {
		return
	}
	timer := activeTimers[gameId.String()]
	if timer == nil {
		return
	}

	stopped := timer.Stop()
	if stopped {
		fmt.Println("Timer stopped for " + gameId.String())
	} else {
		fmt.Println("Failed to stop timer for " + gameId.String())
	}

	delete(activeTimers, gameId.String())
	return
}

// Gets the execute and min kill ts for a game
func (game *Game) GetTimesForGame() (executeTs, minKillTs int64, appErr *ApplicationError) {
	var executeTsBuffer, createTsBuffer time.Time
	err := db.QueryRow("SELECT execute_ts, create_ts FROM dm_kill_timers WHERE game_id = $1", game.GameId.String()).Scan(&executeTsBuffer, &createTsBuffer)
	if err != nil {
		return 0, 0, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return executeTsBuffer.Unix(), createTsBuffer.Unix(), nil
}

// reloads all timers in the database
func LoadAllTimers() (appErr *ApplicationError) {

	// Get all existing kill timers
	rows, err := db.Query(`SELECT game_id, execute_ts FROM dm_kill_timers`)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Loop through the kill timers
	for rows.Next() {
		var gameIdBuffer string
		var executeTsBuffer time.Time
		err := rows.Scan(&gameIdBuffer, &executeTsBuffer)
		// We almost never have scanning errors, but if we do this
		if err != nil {
			msg := `Error loading timer`
			appErr := NewApplicationError(msg, err, ErrCodeDatabase)
			LogWithSentry(appErr, nil, raven.ERROR, map[string]interface{}{"game_id": gameIdBuffer})
			continue
		}
		// Get game id
		gameId := uuid.Parse(gameIdBuffer)
		// Get game
		game, appErr := GetGameById(gameId)
		if appErr != nil {
			LogWithSentry(appErr, nil, raven.ERROR, map[string]interface{}{"game_id": gameId.String()})
			continue
		}

		executeTs := executeTsBuffer.Unix()

		game.LoadTimer(executeTs)

	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Loads a single timer and calls it after the set amount of time
func (game *Game) LoadTimer(executeTs int64) (timer *time.Timer) {

	timezone, appErr := game.GetGameProperty(`timezone`)
	if appErr != nil {
		timezone = Config.DefaultTimeZone
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		fmt.Println(err)
	}

	nowTime := time.Now()
	nowInLoc := nowTime.In(loc)
	now := nowInLoc.Unix()
	timeDiff := executeTs - now
	duration := time.Duration(timeDiff) * time.Second

	if timeDiff <= 0 {
		duration = 10 * time.Minute
	}
	fmt.Println(`Loading timer for ` + game.GameId.String())
	fmt.Print(`Executing in `)
	fmt.Println(duration)
	timer = time.AfterFunc(duration, game.KillTimerHandler)
	addActiveTimer(game.GameId, timer)

	return timer
}

func (game *Game) NewKillTimer(hours int64) (timer *time.Timer, appErr *ApplicationError) {
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	timer, appErr = game.NewKillTimerTransactional(tx, hours)
	if appErr != nil {
		tx.Rollback()
		return nil, appErr
	}
	err = tx.Commit()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return timer, nil
}

// Creates a new kill timer and inserts it into the database
func (game *Game) NewKillTimerTransactional(tx *sql.Tx, hours int64) (timer *time.Timer, appErr *ApplicationError) {

	nowTime := time.Now()
	now := nowTime.Unix()
	// calculate when we need to executee
	executeTs := (hours * SecInHour) + now

	// Insert into db
	insertTimer, err := db.Prepare(`INSERT INTO dm_kill_timers (game_id, create_ts, execute_ts) VALUES ($1, to_timestamp($2), to_timestamp($3))`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(insertTimer).Exec(game.GameId.String(), now, executeTs)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Load the actual timer
	return game.LoadTimer(executeTs), nil
}

// Handler for a timed kill
func (game *Game) KillTimerHandler() {
	appErr := game.ExecuteKillTimer()
	if appErr == nil {
		return
	}
	// if it fails try again in 10 minutes and log it
	fmt.Println(appErr)
	time.AfterFunc(10*time.Minute, game.KillTimerHandler)
	LogWithSentry(appErr, nil, raven.WARNING, map[string]interface{}{"game_id": game.GameId.String()})
}

// Gets the min kill time and executes the kill timer
func (game *Game) ExecuteKillTimer() (appErr *ApplicationError) {
	var minKillTimeBuffer time.Time
	err := db.QueryRow(`SELECT create_ts FROM dm_kill_timers where game_id = $1`, game.GameId.String()).Scan(&minKillTimeBuffer)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	minKillTime := minKillTimeBuffer.Unix()

	// Kill users
	killedUsers, appErr := game.KillPlayersWhoHaventKilledSince(minKillTime)
	if appErr != nil {
		return appErr
	}

	// Inform users the countdown is over
	appErr = game.SendTimerExpiredEmail(killedUsers)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.ERROR, map[string]interface{}{"game_id": game.GameId.String()})
	}
	return nil
}

// Kill all the players who havent killed in the past x hours and randomize targets
func (game *Game) KillPlayersWhoHaventKilledSince(minKillTime int64) (killedUsers []uuid.UUID, appErr *ApplicationError) {

	fmt.Println(`Killing for: ` + game.GameName)

	// Get last_killed value for all users
	rows, err := db.Query(`SELECT DISTINCT ON (m.user_id) m.user_id, p.value FROM dm_user_game_mapping AS m LEFT OUTER JOIN dm_user_properties AS p ON m.user_id = p.user_id AND p.key='last_killed' WHERE m.game_id = $1 AND m.alive = true AND (m.user_role = 'dm_captain' OR m.user_role='dm_user')`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var toBeKilled []uuid.UUID

	for rows.Next() {
		var userIdBuffer string
		var lastKilledBuffer sql.NullString

		// Scan userId and lastKilled
		err = rows.Scan(&userIdBuffer, &lastKilledBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Set up lastKilled in advance
		var lastKilled int64
		lastKilled = 0

		// If the selected lastKilled is valid parse it to a float
		if lastKilledBuffer.Valid {
			lastKilled, err = strconv.ParseInt(lastKilledBuffer.String, 10, 64)
			if err != nil {
				return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
			}
		}

		// If lastKilled is at least the minimum kill time continue
		if lastKilled >= minKillTime {
			continue
		}

		// append the userId to the kill list
		userId := uuid.Parse(userIdBuffer)
		toBeKilled = append(toBeKilled, userId)
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Create a interface slice to store the users to be killed and the gameId
	toBeKilledInterface := ConvertUUIDSliceToInterface(toBeKilled)
	var toBeKilledUpdate []interface{}
	toBeKilledUpdate = append(toBeKilledUpdate, game.GameId.String())
	toBeKilledUpdate = append(toBeKilledUpdate, toBeKilledInterface...)

	// Get the params string for the update
	params := GetParamsForSlice(1, toBeKilledUpdate)

	// Beging transaction for inserts
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Kill the users
	killUsers, err := tx.Prepare(`UPDATE dm_user_game_mapping SET alive = false WHERE game_id = $1 AND user_id IN (` + params + `)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(killUsers).Exec(toBeKilledUpdate...)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Kill the users
	removeTimer, err := tx.Prepare(`DELETE FROM dm_kill_timers WHERE game_id = $1`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(removeTimer).Exec(game.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Assign new targets
	appErr = game.AssignTargetsBy(`normal`)
	if appErr != nil {
		return nil, appErr
	}

	return toBeKilled, nil
}

func (game *Game) DeleteKillTimer() (appErr *ApplicationError) {
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr = game.DeleteKillTimerTransactional(tx)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	// Check transaction for errors
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil

}

func (game *Game) DeleteKillTimerTransactional(tx *sql.Tx) (appErr *ApplicationError) {
	stopTimer(game.GameId)
	// prepare the statement to delete related kill timers
	deleteTimer, err := db.Prepare("DELETE FROM dm_kill_timers WHERE game_id = $1")
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to delete related kill timers
	_, err = tx.Stmt(deleteTimer).Exec(game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return nil
}
