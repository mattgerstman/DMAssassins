package main

import (
	"code.google.com/p/go-uuid/uuid"
	"crypto/sha1"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	UserId     uuid.UUID            `json:"user_id"`
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

	res, err := db.Exec(`INSERT INTO dm_users (user_id, username, email, secret, facebook_id) VALUES ($1,$2,$3,$4,$5)`, userId, username, email, secret, facebookId)
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
	var secret, email, facebookId string
	err := db.QueryRow(`SELECT user_id, email, secret, facebook_id FROM dm_users WHERE username = $1`, username).Scan(&userId, &email, &secret, &facebookId)
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

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(userId uuid.UUID) (*User, *ApplicationError) {
	var username, secret, email, facebookId string
	err := db.QueryRow(`SELECT username, email, secret, facebook_id FROM dm_users WHERE user_id = $1`, userId).Scan(&username, &email, &secret, &facebookId)
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
	var userId uuid.UUID
	var username, email, facebookId string
	err := db.QueryRow(`SELECT user_id, username, email, facebook_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets WHERE user_id = $1)`, user.UserId).Scan(&userId, &username, &email, &facebookId)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	target := &User{userId, user.Username, username, email, "", facebookId, nil}
	_, appErr := target.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return target, nil
}

func (user *User) GetArbitraryGame() (*Game, *ApplicationError) {
	var gameId uuid.UUID
	err := db.QueryRow(`SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1 ORDER BY alive DESC LIMIT 1`, user.UserId).Scan(&gameId)
	switch {
	case err == sql.ErrNoRows:
		msg := "User: " + user.Username + " is not mapped to any games"
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return GetGameById(gameId)

}

func (user *User) UpdateToken(facebook_token string) *ApplicationError {

	res, err := db.Exec(`UPDATE dm_users SET facebook_token = $1 WHERE user_id = $2`, facebook_token, user.UserId)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	return nil
}

func Sha1Hash(plaintext string) string {
	bv := []byte(plaintext)
	sha := sha1.Sum(bv)
	return string(sha[:sha1.Size])
}

func (user *User) GetHashedToken() (string, *ApplicationError) {
	var facebook_token string
	err := db.QueryRow(`SELECT facebook_token FROM dm_users WHERE user_id = $1`, user.UserId).Scan(&facebook_token)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return Sha1Hash(facebook_token), nil

}

//Kills an Assassin's target, user must be logged in
func (user *User) KillTarget(gameId uuid.UUID, secret string) (uuid.UUID, *ApplicationError) {

	var oldTargetId, newTargetId  uuid.UUID

	var targetSecret string
	// Grab the target's secret and user_id for comparison/use below
	err := db.QueryRow(`SELECT secret, user_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1 AND game_id = $2)`, user.UserId, gameId).Scan(&targetSecret, &oldTargetId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

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
	_, err = tx.Stmt(setDead).Exec(oldTargetId, gameId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the old target's target to assign to the Assassin
	err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1 AND game_id = $2)`, user.UserId, gameId).Scan(&newTargetId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Delete the row for the dead user's target
	removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1 AND game_id = $2)`)
	_, err = tx.Stmt(removeOldTarget).Exec(user.UserId, gameId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set up the Assassin's new target
	setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2 AND game_id = $3`)
	_, err = tx.Stmt(setNewTarget).Exec(newTargetId, user.UserId, gameId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	tx.Commit()

	return newTargetId, nil
}
