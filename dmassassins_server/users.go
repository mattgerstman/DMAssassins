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

// Returns a dumb user object with just a userId to avoid DB Queries of setting up a full struct
// This should be used extremely cautiously and only for items like batch setting user properties
func GetDumbUser(userId uuid.UUID) (user *User) {
	properties := make(map[string]string)
	return &User{userId, "", "", "", properties}
}

// Add a user to the DB and return it as a user object
func NewUser(username, email, facebookId string, properties map[string]string) (user *User, appErr *ApplicationError) {
	// Generate the UUID and insert it
	userId := uuid.NewRandom()

	_, err := db.Exec(`INSERT INTO dm_users (user_id, username, email, facebook_id) VALUES ($1,$2,$3,$4)`, userId.String(), username, email, facebookId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
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
	// Get the user info from the db
	err := db.QueryRow(`SELECT user_id, email, facebook_id FROM dm_users WHERE username = $1`, username).Scan(&userIdBuffer, &email, &facebookId)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Parse the user_id into a uuid
	userId = uuid.Parse(userIdBuffer)

	// Build the user struct
	user = &User{userId, username, email, facebookId, nil}

	// Get the user's properties
	_, appErr = user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}

	// Return the user
	return user, nil
}

// Get a user From the DB by their email
func GetUserByEmail(email string) (user *User, appErr *ApplicationError) {
	var userId uuid.UUID
	var username, facebookId, userIdBuffer string
	// Get the user info from the db
	err := db.QueryRow(`SELECT user_id, username, facebook_id FROM dm_users WHERE email = $1`, email).Scan(&userIdBuffer, &username, &facebookId)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundEmail)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Parse the user_id into a uuid
	userId = uuid.Parse(userIdBuffer)

	// Build the user struct
	user = &User{userId, username, email, facebookId, nil}
	_, appErr = user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	// Return the user
	return user, nil
}

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(userId uuid.UUID) (user *User, appErr *ApplicationError) {
	var username, email, facebookId string
	// Get the user info from the db
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

	// Get the game mapping
	gameMapping, appErr := GetGameMapping(user.UserId, gameId)
	if appErr != nil {
		return appErr
	}

	// Convert al lthe game mapping columns to userProperties
	user.Properties["secret"] = gameMapping.Secret
	user.Properties["user_role"] = gameMapping.UserRole
	user.Properties["alive"] = strconv.FormatBool(gameMapping.Alive)
	user.Properties["team"] = ""

	if gameMapping.TeamId == nil {
		return nil
	}

	// Gets the user's team
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
	token, appErr := ExtendToken(facebook_token)
	if appErr != nil {
		return appErr
	}

	_, err := db.Exec(`UPDATE dm_users SET facebook_token = $1 WHERE user_id = $2`, token, user.UserId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
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
	_, err := db.Exec(`UPDATE dm_users SET email = $1 WHERE user_id = $2`, email, user.UserId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	user.Email = email
	return nil
}

func (user *User) KillTarget(gameId uuid.UUID, secret string, trueKill bool) (newTargetId, oldTargetId uuid.UUID, appErr *ApplicationError) {
	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	newTargetId, oldTargetId, appErr = user.KillTargetTransactional(tx, gameId, secret, trueKill)
	if appErr != nil {
		tx.Rollback()
		return nil, nil, appErr
	}

	// Check for errors on transaction commit
	err = tx.Commit()
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return newTargetId, oldTargetId, nil
}

//Kills an Assassin's target, user must be logged in
func (user *User) KillTargetTransactional(tx *sql.Tx, gameId uuid.UUID, secret string, trueKill bool) (newTargetId, oldTargetId uuid.UUID, appErr *ApplicationError) {

	var targetSecret string
	var oldTargetIdBuffer, newTargetIdBuffer, oldTargetTeamIdBuffer sql.NullString
	// Grab the target's secret and user_id for comparison/use below

	err := db.QueryRow(`SELECT map.secret, users.user_id, map.team_id FROM dm_users as users, dm_user_game_mapping as map WHERE users.user_id = (SELECT target_id FROM dm_user_targets where user_id = $1 AND game_id = $2) AND map.user_id = users.user_id AND map.game_id = $3 `, user.UserId.String(), gameId.String(), gameId.String()).Scan(&targetSecret, &oldTargetIdBuffer, &oldTargetTeamIdBuffer)
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	oldTargetTeamId := uuid.Parse(oldTargetTeamIdBuffer.String)
	oldTargetId = uuid.Parse(oldTargetIdBuffer.String)

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
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setDead).Exec(oldTargetId.String(), gameId.String())
	if err != nil {
		fmt.Println(err)
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the old target's target to assign to the Assassin
	err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = $1 AND game_id = $2`, oldTargetId.String(), gameId.String()).Scan(&newTargetIdBuffer)
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	newTargetId = uuid.Parse(newTargetIdBuffer.String)

	// Prepare the statement to delete the row for the dead user's target
	removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to delete the row for the dead user's target
	_, err = tx.Stmt(removeOldTarget).Exec(oldTargetId.String(), gameId.String())
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set up the Assassin's new target
	setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2 AND game_id = $3`)
	_, err = tx.Stmt(setNewTarget).Exec(newTargetId.String(), user.UserId.String(), gameId.String())
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	if !trueKill {
		return newTargetId, oldTargetId, nil
	}

	// Do anything necessary for a plot twist here
	appErr = user.HandlePlotTwistOnKill(tx, oldTargetId, gameId, oldTargetTeamId)
	if appErr != nil {
		return nil, nil, appErr
	}

	// Update kill count
	updateKills, err := db.Prepare(`UPDATE dm_user_game_mapping SET kills = kills + 1 WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(updateKills).Exec(user.UserId.String(), gameId.String())
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Update last kill timestamp
	nowTime := time.Now()
	now := nowTime.Unix()
	lastKill := strconv.FormatInt(now, 10)
	appErr = user.SetUserPropertyTransactional(tx, `last_killed`, lastKill)
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr = user.HandleKillPost(gameId, oldTargetId)
	fmt.Println(`Kill Post`)
	fmt.Println(appErr)
	// DROIDS LOG THIS ERROR

	return newTargetId, oldTargetId, nil
}

func (user *User) Equal(user2 *User) bool {
	if user2 == nil {
		return false
	}

	return uuid.Equal(user.UserId, user2.UserId)
}
