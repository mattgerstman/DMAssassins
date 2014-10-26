package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
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

func (game *Game) NewKillTimer(tx *sql.Tx, hours int64) (timer *time.Timer, appErr *ApplicationError) {

	nowTime := time.Now()
	now := nowTime.Unix()
	executeTs := (hours * SecInHour) + now

	insertTimer, err := db.Prepare(`INSERT INTO dm_kill_timers (game_id, create_ts, execute_ts) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(insertTimer).Exec(game.GameId.String(), now, executeTs)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	duration := time.Duration(now-executeTs) * time.Second

	return time.AfterFunc(duration, game.KillTimerHandler), nil
}

func (game *Game) KillTimerHandler() {
	appErr := game.ExecuteKillTimer()
	if appErr == nil {
		return
	}

	time.AfterFunc(10*time.Minute, game.KillTimerHandler)
	LogWithSentry(appErr, map[string]string{"game_id": game.GameId.String()}, raven.WARNING)
}

func (game *Game) ExecuteKillTimer() (appErr *ApplicationError) {
	var minKillTime int64
	err := db.QueryRow(`SELECT create_ts FROM dm_kill_timers where game_id = $1`, game.GameId.String()).Scan(&minKillTime)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr = game.KillPlayersWhoHaventKilledSince(minKillTime)
	if appErr != nil {
		return appErr
	}
	return nil
}

// Kill all the players who havent killed in the past x hours and randomize targets
func (game *Game) KillPlayersWhoHaventKilledSince(minKillTime int64) (appErr *ApplicationError) {

	// Get last_killed value for all users
	rows, err := db.Query(`SELECT DISTINCT ON (m.user_id) m.user_id, p.value FROM dm_user_game_mapping AS m LEFT OUTER JOIN dm_user_properties AS p ON m.user_id = p.user_id AND p.key='last_killed' WHERE m.game_id = $1 AND m.alive = true AND (m.user_role = 'dm_captain' OR m.user_role='dm_user')`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var toBeKilled []uuid.UUID

	for rows.Next() {
		var userIdBuffer string
		var lastKilledBuffer sql.NullString

		// Scan userId and lastKilled
		err = rows.Scan(&userIdBuffer, &lastKilledBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Set up lastKilled in advance
		var lastKilled int64
		lastKilled = 0

		// If the selected lastKilled is valid parse it to a float
		if lastKilledBuffer.Valid {
			lastKilled, err = strconv.ParseInt(lastKilledBuffer.String, 10, 64)
			if err != nil {
				return NewApplicationError("Internal Error", err, ErrCodeDatabase)
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
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Kill the users
	killUsers, err := db.Prepare(`UPDATE dm_user_game_mapping SET alive = false WHERE game_id = $1 AND user_id IN (` + params + `)`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(killUsers).Exec(toBeKilledUpdate...)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	appErr = game.AssignTargetsByTransactional(tx, `normal`)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	// Kill the users
	removeTimer, err := db.Prepare(`DELETE FROM dm_kill_timers WHERE game_id = $1`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(removeTimer).Exec(game.GameId.String())
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}
