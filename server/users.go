package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	User_id         string            `json:"user_id"`
	Username        string            `json:"username"`
	Email           string            `json:"email"`
	Secret          string            `json:"secret"`
	Properties      map[string]string `json:"properties"`
	hashed_password []byte
}

// Function to clear user password from memory after login it
func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

// Encrypt a PW for storage
func Crypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// Add a user to the DB and return it as a user object
func NewUser(username string, email, plainPW string, secret string) (*User, *ApplicationError) {

	id := uuid.New()
	password, err := Crypt([]byte(plainPW))

	_, err = db.Exec(`INSERT INTO dm_users (user_id, username, email, password, secret) VALUES ($1,$2,$3,$4)`, id, username, email, password, secret)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return &User{id, username, email, secret, nil, password}, nil
}

// Select a User from the DB by username and return it as a user object
func GetUserByUsername(username string) (*User, *ApplicationError) {
	var user_id, secret, email, hashed_password string
	err := db.QueryRow(`SELECT user_id, secret, email, password FROM dm_users WHERE username = $1`, username).Scan(&user_id, &secret, &email, &hashed_password)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	default:
		user := &User{user_id, username, email, secret, nil, []byte(hashed_password)}
		properties, appErr := user.GetUserProperties()
		user.Properties = properties
		if appErr != nil {
			return nil, appErr
		}
		return user, nil
	}
}

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(user_id string) (*User, *ApplicationError) {
	var username, email, secret, hashed_password string
	err := db.QueryRow(`SELECT username, email, secret, password FROM dm_users WHERE user_id = $1`, user_id).Scan(&username, &email, &secret, &hashed_password)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + user_id
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUserId)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	default:
		user := &User{user_id, username, email, secret, nil, []byte(hashed_password)}
		properties, appErr := user.GetUserProperties()
		user.Properties = properties
		if appErr != nil {
			return nil, appErr
		}
		return user, nil
	}
}

func (user *User) GetTarget() (*User, *ApplicationError) {

	var user_id, username, email string
	err := db.QueryRow(`SELECT user_id, username, email FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets WHERE user_id = $1)`, user.User_id).Scan(&user_id, &username, &email)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + username
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUsername)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	default:
		user := &User{user_id, username, email, "", nil, nil}
		properties, appErr := user.GetUserProperties()
		user.Properties = properties
		if appErr != nil {
			return nil, appErr
		}
		return user, nil
	}

	return nil, nil
}

// Confirm a password is equal to it's hashed version
func (user *User) CheckPassword(bytePW []byte) bool {

	if bcrypt.CompareHashAndPassword(user.hashed_password, bytePW) == nil {
		return true
	}
	return false
}

//Kills an Assassin's target, user must be logged in
func (user *User) KillTarget(secret string) (string, *ApplicationError) {

	logged_in_user := user.User_id

	old_target_id := ""
	new_target_id := ""

	var target_secret string
	// Grab the target's secret and user_id for comparison/use below
	err := db.QueryRow(`SELECT secret, user_id FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&target_secret, old_target_id)
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
	err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&new_target_id)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Delete the row for the dead user's target
	removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`)
	_, err = tx.Stmt(removeOldTarget).Exec(logged_in_user)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set up the Assassin's new target
	setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`)
	_, err = tx.Stmt(setNewTarget).Exec(new_target_id, logged_in_user)
	if err != nil {
		tx.Rollback()
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	tx.Commit()

	return new_target_id, nil
}
