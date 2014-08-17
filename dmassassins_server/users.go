package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	UserId     uuid.UUID         `json:"user_id"`
	Assassin   string            `json:"assassin"`
	Username   string            `json:"username"`
	Email      string            `json:"email"`
	Secret     string            `json:"secret"`
	FacebookId string            `json:"facebook_id"`
	Properties map[string]string `json:"properties"`
}

// Add a user to the DB and return it as a user object
func NewUser(username, email, secret, facebookId string, properties map[string]string) (*User, *ApplicationError) {
	userId := uuid.NewUUID()

	res, err := db.Exec(`INSERT INTO dm_users (user_id, username, email, secret, facebook_id) VALUES ($1,$2,$3,$4,$5)`, userId.String(), username, email, secret, facebookId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	if rowsAffected == 0 {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	user := &User{userId, "", username, email, secret, facebookId, properties}
	for key, value := range properties {
		user.SetUserProperty(key, value)
	}
	return user, nil
}

// Select a User from the DB by username and return it as a user object
func GetUserByUsername(username string) (*User, *ApplicationError) {
	var userId uuid.UUID
	var secret, email, facebookId, userIdBuffer string
	err := db.QueryRow(`SELECT user_id, email, secret, facebook_id FROM dm_users WHERE username = $1`, username).Scan(&userIdBuffer, &email, &secret, &facebookId)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	userId = uuid.Parse(userIdBuffer)

	user := &User{userId, "", username, email, secret, facebookId, nil}
	_, appErr := user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(userId uuid.UUID) (*User, *ApplicationError) {
	var username, secret, email, facebookId string
	err := db.QueryRow(`SELECT username, email, secret, facebook_id FROM dm_users WHERE user_id = $1`, userId.String()).Scan(&username, &email, &secret, &facebookId)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	user := &User{userId, "", username, email, secret, facebookId, nil}
	_, appErr := user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

func (user *User) GetTarget() (*User, *ApplicationError) {
	var targetId uuid.UUID
	var username, email, facebookId, targetIdBuffer string
	err := db.QueryRow(`SELECT user_id, username, email, facebook_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets WHERE user_id = $1)`, user.UserId.String()).Scan(&targetIdBuffer, &username, &email, &facebookId)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	targetId = uuid.Parse(targetIdBuffer)
	target := &User{targetId, user.Username, username, email, "", facebookId, nil}
	_, appErr := target.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return target, nil
}

func (user *User) GetArbitraryGame() (*Game, *ApplicationError) {
	var gameId uuid.UUID
	var gameIdBuffer string
	err := db.QueryRow(`SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1 ORDER BY alive DESC LIMIT 1`, user.UserId.String()).Scan(&gameIdBuffer)
	switch {
	case err == sql.ErrNoRows:
		msg := "User: " + user.Username + " is not mapped to any games"
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	gameId = uuid.Parse(gameIdBuffer)
	return GetGameById(gameId)

}

func (user *User) UpdateToken(facebook_token string) *ApplicationError {

	res, err := db.Exec(`UPDATE dm_users SET facebook_token = $1 WHERE user_id = $2`, facebook_token, user.UserId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	return nil
}

func (user *User) GetToken() (string, *ApplicationError) {
	var facebookToken string
	err := db.QueryRow(`SELECT facebook_token FROM dm_users WHERE user_id = $1`, user.UserId.String()).Scan(&facebookToken)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return facebookToken, nil

}

//Kills an Assassin's target, user must be logged in
func (user *User) KillTarget(gameId uuid.UUID, secret string) (uuid.UUID, *ApplicationError) {

	var oldTargetId, newTargetId uuid.UUID
	var targetSecret, oldTargetIdBuffer, newTargetIdBuffer string
	// Grab the target's secret and user_id for comparison/use below
	err := db.QueryRow(`SELECT secret, user_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1 AND game_id = $2)`, user.UserId.String(), gameId.String()).Scan(&targetSecret, &oldTargetIdBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	oldTargetId = uuid.Parse(oldTargetIdBuffer)

	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Confirm the user entered the right secret
	if secret != targetSecret {
		// If secret is invalid throw an error
		msg := fmt.Sprintf("Invalid secret: %s", secret)
		err := errors.New("Invalid Secret")
		return nil, NewApplicationError(msg, err, ErrCodeInvalidSecret)

	}
	// Prepare the statement to kill the old target
	setDead, err := db.Prepare(`UPDATE dm_user_game_mapping SET alive = false WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setDead).Exec(oldTargetId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the old target's target to assign to the Assassin
	err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1 AND game_id = $2)`, user.UserId.String(), gameId.String()).Scan(&newTargetIdBuffer)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	newTargetId = uuid.Parse(newTargetIdBuffer)
	// Delete the row for the dead user's target
	removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1 AND game_id = $2)`)
	_, err = tx.Stmt(removeOldTarget).Exec(user.UserId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set up the Assassin's new target
	setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2 AND game_id = $3`)
	_, err = tx.Stmt(setNewTarget).Exec(newTargetId, user.UserId.String(), gameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	tx.Commit()

	return newTargetId, nil
}
