package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"strconv"
	"time"
)

const (
	TwentyFourHours = 86400
	FourtyEightHours = 172800
)

func (game *Game) isPlotTwistActive(name string)(active bool, appErr *ApplicationError) {
	property, appErr := game.GetGameProperty(name)
	if appErr != nil {
		return false, appErr
	}
	active, err := strconv.ParseBool(property)
	if err != nil {
		return false, NewApplicationError("Internal Error", err, ErrCodeInvalidPlotTwist)
	}
	return active, nil
}

func (user *User) handleSuccessiveKills(tx *sql.Tx, gameId uuid.UUID) (appErr *ApplicationError) {
	lastKilledProperty, appErr := user.GetUserProperty(`last_killed`)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}
	if lastKilledProperty == "" {
		return nil
	}

	nowTime := time.Now()
	now := nowTime.Unix()

	// Parse last killed to int
	lastKilled, err := strconv.ParseInt(lastKilledProperty, 10, 64)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

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

func (user *User) HandlePlotTwistOnKill(tx *sql.Tx, gameId uuid.UUID) (appErr *ApplicationError) {

	// Get the game in question	
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return appErr
	}

	// Check for successive kills plot twist
	successiveKillsIsActive, appErr := game.isPlotTwistActive(`successive_kills`)
	if appErr != nil {
		return appErr
	}

	if successiveKillsIsActive {
		appErr = user.handleSuccessiveKills(tx, gameId)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}
