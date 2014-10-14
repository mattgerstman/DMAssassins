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
	TwentyFourHours  = 86400
	FourtyEightHours = 172800
)

// Handles a Plot Twist that has an immediate action
func (game *Game) handleActionTwist(twistType, twistName string) (appErr *ApplicationError) {
	if twistType == `assign_targets` {
		return game.AssignTargetsBy(twistName)
	}
	if twistType == `kill_users` {
		// DROIDS write kill_users function
		return nil
	}
	return nil
}

// Activates a plot twist
func (game *Game) ActivatePlotTwist(twistType, twistName string) (appErr *ApplicationError) {

	// Determine what type of plot twist we're working with
	typeMap := map[string]string{`kill_mode`: `property`, `kill_users`: `action`, `assign_targets`: `action`}
	if _, ok := typeMap[twistType]; !ok {
		msg := "Invalid Plot Twist Type"
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidPlotTwist)
	}

	// set the game property for a property based twist
	handler := typeMap[twistType]
	if handler == `property` {
		appErr = game.SetGameProperty(twistType, twistName)
		if appErr != nil {
			return appErr
		}
	}

	// execute the action for an action based twist
	if handler == `action` {
		appErr = game.handleActionTwist(twistType, twistName)
		if appErr != nil {
			return appErr
		}
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
func killWeakestPlayerForTeam(gameId, teamId uuid.UUID) (appErr *ApplicationError) {
	// Get the id for the weakest player's assasin
	var assassingIdBuffer, secret string
	err := db.QueryRow(`SELECT targets.user_id, map.secret FROM dm_user_targets as targets, dm_user_game_mapping as map WHERE targets.target_id = map.user_id AND map.game_id = $1 AND map.team_id = $2 AND alive = true ORDER BY map.kills ASC LIMIT 1;=`, gameId.String(), teamId.String()).Scan(&assassingIdBuffer, &secret)
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
	_, _, appErr = assassin.KillTarget(gameId, secret, false)
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

	// Get the weakest player's id
	var weakUserIdBuffer string
	err := db.QueryRow(`SELECT user_id from dm_user_game_mapping WHERE game_id = $1 AND team_id = $2 AND (alive = true OR user_id = $3) ORDER BY kills ASC LIMIT 1`, gameId.String(), teamId.String(), oldTargetId.String()).Scan(&weakUserIdBuffer)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Compare the weakeste player and the given last target, if they match kill the next weakest player
	weakUserId := uuid.Parse(weakUserIdBuffer)
	if uuid.Equal(oldTargetId, weakUserId) {
		return killWeakestPlayerForTeam(gameId, teamId)
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
