package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"fmt"
	//	"log"
)

//All of the functions in here need to switch to the ApplicationError Constructors

type User struct {
	User_id         string
	Email           string
	Secret          string
	hashed_password []byte
}

//Yea you like me clearing it from memory
func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func Crypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}
func NewUser(email string, plainPW string, secret string) (*User, *ApplicationError) {

	id := uuid.New()
	password, err := Crypt([]byte(plainPW))

	_, err = db.Exec(`INSERT INTO dm_users (user_id, email, password, secret) VALUES ($1,$2,$3,$4)`, id, email, password, secret)
	if err != nil {
		//log error
		return nil, &ApplicationError{"Internal Error", err, ERROR_DATABASE}
	}

	return &User{id, email, secret, password}, nil
}

func GetUserByEmail(email string) (*User, *ApplicationError) {
	var user_id, secret, hashed_password string
	err := db.QueryRow(`SELECT user_id, secret, password FROM dm_users WHERE email = $1`, email).Scan(&user_id, &secret, &hashed_password)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + email
		return nil, &ApplicationError{msg, err, ERROR_INVALID_EMAIL}
	default:
		return &User{user_id, email, secret, []byte(hashed_password)}, nil
	}
}

func GetUserById(user_id string) (*User, *ApplicationError) {
	var email, secret, hashed_password string
	err := db.QueryRow(`SELECT email, secret, password FROM dm_users WHERE user_id = $1`, user_id).Scan(&email, &secret, &hashed_password)
	switch {
	case err == sql.ErrNoRows:
		msg := "Invalid user: " + user_id
		return nil, &ApplicationError{msg, err, ERROR_INVALID_USER_ID}
	case err != nil:
		return nil, &ApplicationError{"Internal Error", err, ERROR_DATABASE}
	default:
		return &User{user_id, email, secret, []byte(hashed_password)}, nil
	}
}

func (user *User) CheckPassword(plainPW string) bool {

	bytePW := []byte(plainPW)
	if bcrypt.CompareHashAndPassword(user.hashed_password, bytePW) == nil {
		return true
	}
	return false
}

//I need to handle what happens if each sql statement fails, right now I'm totally bypassing the error checking
func (user *User) KillTarget(secret string) (string, *ApplicationError) {

	transaction, err := db.Begin()
	defer transaction.Commit()

	logged_in_user := user.User_id
	var target_secret string
	err = db.QueryRow(`SELECT secret FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&target_secret)

	new_target_id := ""
	if secret == target_secret {

		var old_target_id string

		err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = $1`, logged_in_user).Scan(&old_target_id)
		_, err = db.Exec(`UPDATE dm_users SET alive = false WHERE user_id = $1`, old_target_id)

		err = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&new_target_id)

		_, err = db.Exec(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`, logged_in_user)
		_, err = db.Exec(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`, new_target_id, logged_in_user)

	} else {
		msg := fmt.Sprintf("Invalid secret: %s", secret)
		return "", NewSimpleApplicationError(msg, ERROR_INVALID_SECRET)
	}

	_ = err

	return new_target_id, nil
}
