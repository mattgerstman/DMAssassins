package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	User_id     string            `json:"user_id"`
	Assassin    string            `json:"assassin"`
	Username    string            `json:"username"`
	Email       string            `json:"email"`
	Secret      string            `json:"secret"`
	Facebook_id string            `json:"facebook_id"`
	Properties  map[string]string `json:"properties"`
}

// Add a user to the DB and return it as a user object
func NewUser(username, email, secret, facebook_id string, properties map[string]string) (*User, *ApplicationError) {
	user_id := uuid.New()

	res, err := db.Exec(`INSERT INTO dm_users (user_id, username, email, secret, facebook_id) VALUES ($1,$2,$3,$4,$5)`, user_id, username, email, secret, facebook_id)
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

	user := &User{user_id, "", username, email, secret, facebook_id, properties}
	for key, value := range properties {
		user.SetUserProperty(key, value)
	}
	return user, nil
}

// Select a User from the DB by username and return it as a user object
func GetUserByUsername(username string) (*User, *ApplicationError) {
	var user_id, secret, email, facebook_id string
	err := db.QueryRow(`SELECT user_id, email, secret, facebook_id FROM dm_users WHERE username = $1`, username).Scan(&user_id, &email, &secret, &facebook_id)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	user := &User{user_id, "", username, email, secret, facebook_id, nil}
	_, appErr := user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(user_id string) (*User, *ApplicationError) {
	var username, secret, email, facebook_id string
	err := db.QueryRow(`SELECT username, email, secret, facebook_id FROM dm_users WHERE user_id = $1`, user_id).Scan(&username, &email, &secret, &facebook_id)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	user := &User{user_id, "", username, email, secret, facebook_id, nil}
	_, appErr := user.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return user, nil
}

func (user *User) GetTarget() (*User, *ApplicationError) {

	var user_id, username, email, facebook_id string
	err := db.QueryRow(`SELECT user_id, username, email, facebook_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets WHERE user_id = $1)`, user.User_id).Scan(&user_id, &username, &email, &facebook_id)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	target := &User{user_id, user.Username, username, email, "", facebook_id, nil}
	_, appErr := target.GetUserProperties()
	if appErr != nil {
		return nil, appErr
	}
	return target, nil
}

//Kills an Assassin's target, user must be logged in
func (user *User) KillTarget(secret string) (string, *ApplicationError) {

	old_target_id := ""
	new_target_id := ""

	var target_secret string
	// Grab the target's secret and user_id for comparison/use below
	err := db.QueryRow(`SELECT secret, user_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, user.User_id).Scan(&target_secret, &old_target_id)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Confirm the user entered the right secret
	if secret != target_secret {
		// If secret is invalid throw an error
		msg := fmt.Sprintf("Invalid secret: %s", secret)
		err := errors.New("Invalid Secret")
		return "", NewApplicationError(msg, err, ErrCodeInvalidSecret)

	}
	// Prepare the statement to kill the old target
	setDead, err := db.Prepare(`UPDATE dm_users SET alive = false WHERE user_id = $1`)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setDead).Exec(old_target_id)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the old target's target to assign to the Assassin
	err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, user.User_id).Scan(&new_target_id)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Delete the row for the dead user's target
	removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`)
	_, err = tx.Stmt(removeOldTarget).Exec(user.User_id)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set up the Assassin's new target
	setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`)
	_, err = tx.Stmt(setNewTarget).Exec(new_target_id, user.User_id)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	tx.Commit()

	return new_target_id, nil
}
