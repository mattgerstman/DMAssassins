package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	User_id         string
	Email           string
	Secret          string
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
func NewUser(email string, plainPW string, secret string) (*User, *ApplicationError) {

	id := uuid.New()
	password, err := Crypt([]byte(plainPW))

	_, err = db.Exec(`INSERT INTO dm_users (user_id, email, password, secret) VALUES ($1,$2,$3,$4)`, id, email, password, secret)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return &User{id, email, secret, password}, nil
}

// Select a User from the DB by email and return it as a user object
func GetUserByEmail(email string) (*User, *ApplicationError) {
	var user_id, secret, hashed_password string
	err := db.QueryRow(`SELECT user_id, secret, password FROM dm_users WHERE email = $1`, email).Scan(&user_id, &secret, &hashed_password)
	fmt.Println(err)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + email
		return nil, NewApplicationError(msg, err, ErrCodeInvalidEmail)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	default:
		return &User{user_id, email, secret, []byte(hashed_password)}, nil
	}
}

// Select a user from the db by user_id (uuid) and return it as a user object
func GetUserById(user_id string) (*User, *ApplicationError) {
	var email, secret, hashed_password string
	err := db.QueryRow(`SELECT email, secret, password FROM dm_users WHERE user_id = $1`, user_id).Scan(&email, &secret, &hashed_password)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + user_id
		return nil, NewApplicationError(msg, err, ErrCodeInvalidUserId)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	default:
		return &User{user_id, email, secret, []byte(hashed_password)}, nil
	}
}

// Confirm a password is equal to it's hashed version
func (user *User) CheckPassword(bytePW []byte) bool {

	if bcrypt.CompareHashAndPassword(user.hashed_password, bytePW) == nil {
		return true
	}
	return false
}

//I need to handle what happens if each sql statement fails, right now I'm totally bypassing the error checking
func (user *User) KillTarget(secret string) (string, *ApplicationError) {

	tx, err := db.Begin()

	logged_in_user := user.User_id
	
	var target_secret string
	err = db.QueryRow(`SELECT secret FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&target_secret)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	new_target_id := ""
	if secret == target_secret {

		var old_target_id string

		err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = $1`, logged_in_user).Scan(&old_target_id)
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		setDead, err := db.Prepare(`UPDATE dm_users SET alive = false WHERE user_id = $1`)
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		_, err = tx.Stmt(setDead).Exec(old_target_id)
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&new_target_id)
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		removeOldTarget, err := db.Prepare(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`)
		_, err = tx.Stmt(removeOldTarget).Exec(logged_in_user)
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		setNewTarget, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`)
		_, err = tx.Stmt(setNewTarget).Exec(new_target_id, logged_in_user)
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

	} else {
		msg := fmt.Sprintf("Invalid secret: %s", secret)
		err := errors.New("Invalid Secret")
		return "", NewApplicationError(msg, err, ErrCodeInvalidSecret)
	}

	_ = err

	tx.Commit()

	return new_target_id, nil
}
