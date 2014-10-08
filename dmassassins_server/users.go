package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type User struct {
	UserId     uuid.UUID         `json:"user_id"`
	Username   string            `json:"username"`
	Email      string            `json:"email"`
	FacebookId string            `json:"facebook_id"`
	Properties map[string]string `json:"properties"`
}

// Add a user to the DB and return it as a user object
func NewUser(username, email, facebookId string, properties map[string]string) (user *User, appErr *ApplicationError) {
	// Generate the UUID and insert it
	userId := uuid.NewRandom()

	res, err := db.Exec(`INSERT INTO dm_users (user_id, username, email, facebook_id) VALUES ($1,$2,$3,$4)`, userId.String(), username, email, facebookId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Makes sure insert worked
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return nil, NoRowsAffectedAppErr
	}

	// Create user
	user = &User{userId, username, email, facebookId, properties}

	// Set properties
	for key, value := range properties {
		user.SetUserProperty(key, value)
	}
	// Return user
	return user, nil
}

// Select a User from the DB by username and return it as a user object
func GetUserByUsername(username string) (user *User, appErr *ApplicationError) {
	var userId uuid.UUID
	var email, facebookId, userIdBuffer string
	err := db.QueryRow(`SELECT user_id, email, facebook_id FROM dm_users WHERE username = $1`, username).Scan(&userIdBuffer, &email, &facebookId)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	userId = uuid.Parse(userIdBuffer)

	user = &User{userId, username, email, facebookId, nil}
	_, appErr = user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

// Get a user From the DB by their email
func GetUserByEmail(email string) (user *User, appErr *ApplicationError) {
	var userId uuid.UUID
	var username, facebookId, userIdBuffer string
	err := db.QueryRow(`SELECT user_id, username, facebook_id FROM dm_users WHERE email = $1`, email).Scan(&userIdBuffer, &username, &facebookId)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundEmail)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	userId = uuid.Parse(userIdBuffer)

	user = &User{userId, username, email, facebookId, nil}
	_, appErr = user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(userId uuid.UUID) (user *User, appErr *ApplicationError) {
	var username, email, facebookId string
	err := db.QueryRow(`SELECT username, email, facebook_id FROM dm_users WHERE user_id = $1`, userId.String()).Scan(&username, &email, &facebookId)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + userId.String()
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundUserId)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Create the user obj
	user = &User{userId, username, email, facebookId, nil}

	// Add user properties
	_, appErr = user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

// Gets game related properties for a user and loads them into the user.Properties map
func (user *User) GetUserGameProperties(gameId uuid.UUID) *ApplicationError {

	gameMapping, appErr := GetGameMapping(user.UserId, gameId)
	if appErr != nil {
		return appErr
	}

	user.Properties["secret"] = gameMapping.Secret
	user.Properties["user_role"] = gameMapping.UserRole
	user.Properties["alive"] = strconv.FormatBool(gameMapping.Alive)
	user.Properties["team"] = ""

	if gameMapping.TeamId == nil {
		return nil
	}

	team, appErr := GetTeamById(gameMapping.TeamId)
	if appErr != nil {
		return nil
	}
	user.Properties["team"] = team.TeamName
	return nil
}

// Gets a user in the context of a game with game related properties
func GetUserForGameById(userId, gameId uuid.UUID) (user *User, appErr *ApplicationError) {
	// Get the user
	user, appErr = GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	// Get the game related properties
	appErr = user.GetUserGameProperties(gameId)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil
}

// Gets a user's target for a game
func (user *User) GetTarget(gameId uuid.UUID) (target *User, appErr *ApplicationError) {
	var targetId uuid.UUID
	var username, email, facebookId, targetIdBuffer string

	// DB query
	err := db.QueryRow(`SELECT user_id, username, email, facebook_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets WHERE user_id = $1 AND game_id = $2)`, user.UserId.String(), gameId.String()).Scan(&targetIdBuffer, &username, &email, &facebookId)

	switch {
	case err == sql.ErrNoRows:
		msg := "No Target"
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundTarget)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Creates the target obj
	targetId = uuid.Parse(targetIdBuffer)
	target = &User{targetId, username, email, facebookId, nil}

	// gets properties for the target
	_, appErr = target.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return target, nil
}

// Updates a user's facebook token
func (user *User) UpdateToken(facebook_token string) (appErr *ApplicationError) {

	res, err := db.Exec(`UPDATE dm_users SET facebook_token = $1 WHERE user_id = $2`, facebook_token, user.UserId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Make sure update worked
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	return nil
}

// Gets a user's facebook token
func (user *User) GetToken() (fbToken string, appErr *ApplicationError) {
	var facebookToken sql.NullString
	err := db.QueryRow(`SELECT facebook_token FROM dm_users WHERE user_id = $1`, user.UserId.String()).Scan(&facebookToken)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return facebookToken.String, nil

}

// Change a user's email
func (user *User) ChangeEmail(email string) (appErr *ApplicationError) {
	if email == user.Email {
		return nil
	}
	res, err := db.Exec(`UPDATE dm_users SET email = $1 WHERE user_id = $2`, email, user.UserId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Make sure update worked
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}
	user.Email = email
	return nil
}

//Kills an Assassin's target, user must be logged in
func (user *User) KillTarget(gameId uuid.UUID, secret string, trueKill bool) (newTargetId, oldTargetId uuid.UUID, appErr *ApplicationError) {

	var targetSecret string
	var oldTargetIdBuffer, newTargetIdBuffer sql.NullString
	// Grab the target's secret and user_id for comparison/use below

	err := db.QueryRow(`SELECT map.secret, users.user_id FROM dm_users as users, dm_user_game_mapping as map WHERE users.user_id = (SELECT target_id FROM dm_user_targets where user_id = $1 AND game_id = $2) AND map.user_id = users.user_id AND map.game_id = $3 `, user.UserId.String(), gameId.String(), gameId.String()).Scan(&targetSecret, &oldTargetIdBuffer)
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	oldTargetId = uuid.Parse(oldTargetIdBuffer.String)

	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Confirm the user entered the right secret
	if !strings.EqualFold(secret, targetSecret) {
		// If secret is invalid throw an error
		msg := fmt.Sprintf("Invalid secret: %s", secret)
		err := errors.New("Invalid Secret")
		return nil, nil, NewApplicationError(msg, err, ErrCodeInvalidSecret)

	}
	// Prepare the statement to kill the old target
	setDead, err := db.Prepare(`UPDATE dm_user_game_mapping SET alive = false WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setDead).Exec(oldTargetId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the old target's target to assign to the Assassin
	err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = $1 AND game_id = $2`, oldTargetId.String(), gameId.String()).Scan(&newTargetIdBuffer)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	newTargetId = uuid.Parse(newTargetIdBuffer.String)

	// Prepare the statement to delete the row for the dead user's target
	removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to delete the row for the dead user's target
	_, err = tx.Stmt(removeOldTarget).Exec(oldTargetId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set up the Assassin's new target
	setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2 AND game_id = $3`)
	_, err = tx.Stmt(setNewTarget).Exec(newTargetId.String(), user.UserId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	if !trueKill {
		tx.Commit()
		return newTargetId, oldTargetId, nil
	}

	// Do anything necessary for a plot twist here
	appErr = user.HandlePlotTwistOnKill(tx, gameId)
	if appErr != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Update kill count
	updateKills, err := db.Prepare(`UPDATE dm_user_game_mapping SET kills = kills + 1 WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(updateKills).Exec(user.UserId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Update last kill timestamp
	nowTime := time.Now()
	now := nowTime.Unix()
	lastKill := strconv.FormatInt(now, 10)
	appErr = user.SetUserPropertyTransactional(tx, `last_killed`, lastKill)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	tx.Commit()

	return newTargetId, oldTargetId, nil
}
