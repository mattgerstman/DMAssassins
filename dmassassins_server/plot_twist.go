package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"github.com/getsentry/raven-go"
	"strconv"
	"time"
)

const (
	SecInHour        = 3600
	TwentyFourHours  = 86400
	FourtyEightHours = 172800
)

// Gets plot twist from the plot twist config
func GetPlotTwist(twistName string) (twist *PlotTwist, appErr *ApplicationError) {

	if twist, ok := PlotTwistConfig[twistName]; ok {
		return twist, nil
	}

	msg := `Invalid Plot Twist: ` + twistName
	err := errors.New(msg)
	return nil, NewApplicationError(msg, err, ErrCodeInvalidPlotTwist)

}

// Revives a list of users, assumes they're all already dead
func (game *Game) ReviveUsers(tx *sql.Tx, toBeRevived []uuid.UUID) (appErr *ApplicationError) {

	// Create a interface slice to store the users to be killed and the
	toBeRevivedInterface := ConvertUUIDSliceToInterface(toBeRevived)

	var toBeRevivedUpdate []interface{}
	toBeRevivedUpdate = append(toBeRevivedUpdate, game.GameId.String())
	toBeRevivedUpdate = append(toBeRevivedUpdate, toBeRevivedInterface...)

	// Get the params string for the update
	params := GetParamsForSlice(1, toBeRevivedUpdate)

	// Prepare statement to revive captains
	reviveStmt, err := db.Prepare(`UPDATE dm_user_game_mapping SET alive = true WHERE game_id = $1 AND user_id IN (` + params + `)`)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute statement to revive captains
	_, err = tx.Stmt(reviveStmt).Exec(toBeRevivedUpdate...)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Mark users as revived
	for _, userId := range toBeRevived {
		user := GetDumbUser(userId)
		appErr = user.SetUserPropertyTransactional(tx, `revived`, `true`)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

// Revives all of the captains who are dead
func (game *Game) ReviveCaptains(tx *sql.Tx) (appErr *ApplicationError) {
	// Get all the dead team captains
	toBeRevived, appErr := game.GetDeadTeamCaptains()
	if appErr != nil {
		return appErr
	}
	// Revive them
	return game.ReviveUsers(tx, toBeRevived)

}

// Revives all of the strongest players who are dead
func (game *Game) ReviveStrongestPlayers(tx *sql.Tx) (appErr *ApplicationError) {
	// Get all the strongest players
	toBeRevived, appErr := game.getStrongPlayersWithState(false)
	if appErr != nil {
		return appErr
	}
	// Revive them
	return game.ReviveUsers(tx, toBeRevived)
}

// Kill all the players who havent killed in the past x hours and randomize targets
func (game *Game) KillPlayersWithNoRecentKills(hours float64) (appErr *ApplicationError) {

	// Get last_killed value for all users
	rows, err := db.Query(`SELECT DISTINCT ON (m.user_id) m.user_id, p.value FROM dm_user_game_mapping AS m LEFT OUTER JOIN dm_user_properties AS p ON m.user_id = p.user_id AND p.key='last_killed' WHERE m.game_id = $1 AND m.alive = true AND (m.user_role = 'dm_captain' OR m.user_role='dm_user')`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get current time in seconds
	now := time.Now()

	var toBeKilled []uuid.UUID
	minKillTime := float64(now.Unix()) - (hours * float64(SecInHour))

	for rows.Next() {
		var userIdBuffer string
		var lastKilledBuffer sql.NullString

		// Scan userId and lastKilled
		err = rows.Scan(&userIdBuffer, &lastKilledBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Set up lastKilled in advance
		var lastKilled float64
		lastKilled = 0

		// If the selected lastKilled is valid parse it to a float
		if lastKilledBuffer.Valid {
			lastKilled, err = strconv.ParseFloat(lastKilledBuffer.String, 64)
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

	// Kill the users
	_, err = db.Exec(`UPDATE dm_user_game_mapping SET alive = false WHERE game_id = $1 AND user_id IN (`+params+`)`, toBeKilledUpdate...)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return game.AssignTargetsBy(`normal`)
}

// Kill all the players with 0 kills and randomize targets
func (game *Game) KillPlayersWithNoKills() (appErr *ApplicationError) {
	_, err := db.Exec(`UPDATE dm_user_game_mapping SET alive = false WHERE kills = 0 AND game_id = $1 `, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return game.AssignTargetsBy(`normal`)
}

// Activates a plot twist
func (game *Game) ActivatePlotTwistPart(tx *sql.Tx, twistName, twistValue string) (appErr *ApplicationError) {

	switch twistName {
	case `assign_targets`:
		return
	case `kill_mode`:
		return
	case `kill_innocent`:
		return game.KillPlayersWithNoKills()
	case `kill_inactive`:
		numHours, err := strconv.ParseFloat(twistValue, 64)
		if err != nil || numHours == 0 {
			return NewApplicationError(`Invalid Number of Hours`, err, ErrCodeInvalidParameter)
		}
		return game.KillPlayersWithNoRecentKills(numHours)
	}

	msg := `Invalid Plot Twist: ` + twistName
	err := errors.New(msg)
	return NewApplicationError(msg, err, ErrCodeInvalidPlotTwist)
}

// Revive a group of players
func (game *Game) RevivePlayers(tx *sql.Tx, revive string) (appErr *ApplicationError) {
	switch revive {
	case `revive_strongest`:
		return game.ReviveStrongestPlayers(tx)
	case `revive_captains`:
		return game.ReviveCaptains(tx)
	}
	return nil
}

// Activates a plot twist
func (game *Game) ActivatePlotTwist(twistName string) (appErr *ApplicationError) {
	twist, appErr := GetPlotTwist(twistName)
	if appErr != nil {
		return appErr
	}

	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Revive anyone we need to revive
	revive := twist.Revive
	if revive != `` {
		appErr = game.RevivePlayers(tx, revive)
		if appErr != nil {
			tx.Rollback()
			return appErr
		}
	}

	// Assign targets
	assignMode := twist.AssignTargets
	if assignMode != `` {
		appErr = game.AssignTargetsByTransactional(tx, assignMode)
		if appErr != nil {
			tx.Rollback()
			return appErr
		}
	}

	// Set kill mode
	appErr = game.SetGamePropertyTransactional(tx, `kill_mode`, twist.KillMode)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	// Set kill tiemr
	killTimer := twist.KillTimer
	if killTimer != 0 {
		appErr = game.SetKillTimer(tx, killTimer)
		if appErr != nil {
			tx.Rollback()
			return appErr
		}
	}

	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Determine if the user has made another kill in the past 24 hours and if so gives them an extra point
func (user *User) handleSuccessiveKills(tx *sql.Tx, gameId uuid.UUID) (appErr *ApplicationError) {
	// Get last time they killed
	lastKilledProperty, appErr := user.GetUserProperty(`last_killed`)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}
	// If we don't have this property they haven't killed recently
	if lastKilledProperty == "" {
		return nil
	}

	// Get current time in seconds
	nowTime := time.Now()
	now := nowTime.Unix()

	// Parse last killed to int
	lastKilled, err := strconv.ParseInt(lastKilledProperty, 10, 64)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Check if the user has killed in the last 24 hours
	if (now - lastKilled) > TwentyFourHours {
		return nil
	}

	// Update kill count
	updateKills, err := db.Prepare(`UPDATE dm_user_game_mapping SET kills = kills + 1 WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(updateKills).Exec(user.UserId.String(), gameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil

}

// Kills the weakest player on a team
func killWeakestPlayerForTeam(tx *sql.Tx, gameId, teamId uuid.UUID) (appErr *ApplicationError) {
	// Get the id for the weakest player's assasin
	var assassingIdBuffer, secret string
	err := db.QueryRow(`SELECT targets.user_id, map.secret FROM dm_user_targets as targets, dm_user_game_mapping as map WHERE targets.target_id = map.user_id AND map.game_id = $1 AND map.team_id = $2 AND alive = true ORDER BY map.kills ASC LIMIT 1`, gameId.String(), teamId.String()).Scan(&assassingIdBuffer, &secret)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the assassin
	assassinId := uuid.Parse(assassingIdBuffer)
	assassin, appErr := GetUserById(assassinId)
	if appErr != nil {
		return appErr
	}

	// kill the assassins target (silently)
	_, _, appErr = assassin.KillTargetTransactional(tx, gameId, secret, false)
	if appErr != nil {
		return appErr
	}

	gameName := ``

	game, appErr := GetGameById(gameId)
	if appErr == nil {
		gameName = game.GameName
	} else {
		LogWithSentry(appErr, map[string]string{"assassin_id": assassin.UserId.String(), "game_id": gameId.String()}, raven.WARNING)
	}

	// DROIDS switch to plot twist email
	// Inform the assassin they have a new target
	_, appErr = assassin.SendNewTargetEmail(gameName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"assassin_id": assassin.UserId.String(), "game_id": gameId.String()}, raven.WARNING)

	}

	return nil

}

// Check if the user killed is the weakest player for their team, if so kill the weakest player for that team
func (user *User) handleDefendWeak(tx *sql.Tx, oldTargetId, gameId, teamId uuid.UUID) (appErr *ApplicationError) {

	// DROIDS HANDLE TIE FOR WEAKEST PLAYER
	// Get the weakest player's id
	var weakUserIdBuffer string
	err := db.QueryRow(`SELECT user_id from dm_user_game_mapping WHERE game_id = $1 AND team_id = $2 AND (alive = true OR user_id = $3) ORDER BY kills ASC LIMIT 1`, gameId.String(), teamId.String(), oldTargetId.String()).Scan(&weakUserIdBuffer)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Compare the weakeste player and the given last target, if they match kill the next weakest player
	weakUserId := uuid.Parse(weakUserIdBuffer)
	if uuid.Equal(oldTargetId, weakUserId) {
		return killWeakestPlayerForTeam(tx, gameId, teamId)
	}

	return nil

}

// Determine if there is an active plot twist that takes place on kill and if so execute it
func (user *User) HandlePlotTwistOnKill(tx *sql.Tx, oldTargetId, gameId, teamId uuid.UUID) (appErr *ApplicationError) {

	// Get the game in question
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return appErr
	}

	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return appErr
	}
	if teamsEnabled != `true` {
		return nil
	}

	// Check for successive kills plot twist
	killMode, appErr := game.GetGameProperty(`kill_mode`)
	if appErr != nil {
		return appErr
	}

	// Check plot twist mode
	switch killMode {
	case `successive_kills`:
		return user.handleSuccessiveKills(tx, gameId)
	case `defend_weak`:
		return user.handleDefendWeak(tx, oldTargetId, gameId, teamId)
	}

	return nil
}
