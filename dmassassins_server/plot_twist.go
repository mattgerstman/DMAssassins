package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"strconv"
	"time"
)

type PlotTwist struct {
	AssignTargets string `json:"assign_targets"`
	KillMode      string `json:"kill_mode"`
	KillTimer     int64  `json:"kill_timer"`
	RequireTeams  bool   `json:"require_teams"`
	Revive        string `json:"revive"`
}

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

// Revive a group of players
func (game *Game) RevivePlayers(revive string) (appErr *ApplicationError) {

	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	switch revive {
	case `revive_strongest`:
		appErr = game.ReviveStrongestPlayers(tx)
	case `revive_captains`:
		appErr = game.ReviveCaptains(tx)
	}
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Activates a plot twist
func (game *Game) ActivatePlotTwist(twistName string) (appErr *ApplicationError) {
	if !game.Started {
		msg := "The game must be started to execute a plot twist!"
		err := errors.New("User tried to start a plot twist for an unstarted game")
		return NewApplicationError(msg, err, ErrCodeGameNotStarted)
	}

	twist, appErr := GetPlotTwist(twistName)
	if appErr != nil {
		return appErr
	}

	// Mae sure we have teams if we need them
	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return appErr
	}
	if teamsEnabled != `true` && twist.RequireTeams {
		msg := "Incompatible Plot Twist"
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidPlotTwist)
	}

	// Revive anyone we need to revive
	revive := twist.Revive
	if revive != `` {
		appErr = game.RevivePlayers(revive)
		if appErr != nil {
			return appErr
		}
	}

	// Assign targets
	assignMode := twist.AssignTargets
	if assignMode != `` {
		appErr = game.AssignTargetsBy(assignMode)
		if appErr != nil {
			return appErr
		}
	}

	// Set kill mode
	appErr = game.SetGameProperty(`kill_mode`, twist.KillMode)
	if appErr != nil {
		return appErr
	}

	// Set kill tiemr
	killTimer := twist.KillTimer
	if killTimer != 0 {
		_, appErr = game.NewKillTimer(killTimer)
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

// Kills the next weakest player on a team
func killNextWeakestPlayerForTeam(tx *sql.Tx, gameId, teamId, userId uuid.UUID) (appErr *ApplicationError) {
	// Get the id for the weakest player's assasin
	var assassingIdBuffer, secret string
	err := db.QueryRow(`SELECT targets.user_id, map.secret FROM dm_user_targets as targets, dm_user_game_mapping as map WHERE targets.target_id = map.user_id AND map.game_id = $1 AND map.team_id = $2 AND map.user_id != $3 AND alive = true ORDER BY map.kills ASC LIMIT 1`, gameId.String(), teamId.String(), userId.String()).Scan(&assassingIdBuffer, &secret)
	if err == sql.ErrNoRows {
		return nil
	}
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
	_, oldTargetId, appErr := assassin.KillTargetTransactional(tx, gameId, secret, false)
	if appErr != nil {
		return appErr
	}

	gameName := ``
	extra := map[string]interface{}{"assassin_id": assassinId, "game_id": gameId.String(), "old_target_id": oldTargetId.String()}

	game, appErr := GetGameById(gameId)
	if appErr == nil {
		gameName = game.GameName
	} else {
		LogWithSentry(appErr, nil, raven.WARNING, extra)
	}

	// Inform the assassin they have a new target
	_, appErr = assassin.SendDefendWeakNewTargetEmail(gameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra)
	}

	oldTarget, appErr := GetUserById(oldTargetId)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra)
		return nil
	}
	_, appErr = oldTarget.SendDefendWeakKilledEmail(gameName)
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.WARNING, extra)
	}

	return nil

}

// Check if the user killed is the weakest player for their team, if so kill the weakest player for that team
func (user *User) handleDefendWeak(tx *sql.Tx, oldTargetId, gameId, teamId uuid.UUID) (appErr *ApplicationError) {

	// Get the weakest player's id
	rows, err := db.Query(`SELECT user_id FROM (SELECT *, RANK() OVER (PARTITION BY team_id ORDER BY kills ASC) AS rnum FROM dm_user_game_mapping WHERE  (alive = true OR user_id = $1)) s WHERE s.team_id = $2 AND s.rnum = 1;`, oldTargetId.String(), teamId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Iterate through weak players
	for rows.Next() {
		var weakUserIdBuffer string
		err = rows.Scan(&weakUserIdBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Compare the current weakest player and the given last target, if they match kill the next weakest player
		weakUserId := uuid.Parse(weakUserIdBuffer)
		if uuid.Equal(oldTargetId, weakUserId) {
			fmt.Println(`Kill weakest`)
			return killNextWeakestPlayerForTeam(tx, gameId, teamId, oldTargetId)
		}
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
