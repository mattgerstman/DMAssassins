package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	User_id         string
	Email           string
	Secret          string
	hashed_password []byte
}

func clear(b []byte) {
	for i := 0; i < len(b); i++ {
		b[i] = 0
	}
}

func Crypt(password []byte) ([]byte, error) {
	defer clear(password)
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

func NewUser(email string, plainPW string, secret string) *User {

	id := uuid.New()
	password, err := Crypt([]byte(plainPW))

	_, err = db.Query(`INSERT INTO dm_users (user_id, email, password, secret) VALUES ($1,$2,$3,$4)`, id, email, password, secret)
	fmt.Println(err)
	returnUser := &User{id, email, secret, password}
	return returnUser
}

func GetUserByEmail(email string) *User {
	var user_id, secret, hashed_password string
	err := db.QueryRow(`SELECT user_id, secret, password FROM dm_users WHERE email = $1`, email).Scan(&user_id, &secret, &hashed_password)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	default:
		return &User{user_id, email, secret, []byte(hashed_password)}
	}
	return nil
}

func GetUserById(user_id string) *User {
	var email, secret, hashed_password string
	err := db.QueryRow(`SELECT email, secret, password FROM dm_users WHERE user_id = $1`, user_id).Scan(&email, &secret, &hashed_password)
	switch {
	case err == sql.ErrNoRows:
		return nil
	case err != nil:
		log.Fatal(err)
	default:
		return &User{user_id, email, secret, []byte(hashed_password)}
	}
	return nil
}

func (user *User) CheckPassword(plainPW string) bool {

	bytePW := []byte(plainPW)
	if bcrypt.CompareHashAndPassword(user.hashed_password, bytePW) == nil {
		return true
	}
	return false
}

func (user *User) KillTarget(secret string) string {

	transaction, _ := db.Begin()
	defer transaction.Commit()
	logged_in_user := user.User_id
	var target_secret string
	db.QueryRow(`SELECT secret FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&target_secret)

	var new_target_id string
	new_target_id = ""
	if secret == target_secret {

		var new_target_id, old_target_id string

		db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = $1`, logged_in_user).Scan(&old_target_id)
		db.QueryRow(`UPDATE dm_users SET alive = false WHERE user_id = $1`, old_target_id)

		db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user).Scan(&new_target_id)

		db.QueryRow(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`, logged_in_user)
		db.QueryRow(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`, new_target_id, logged_in_user)

	} else {
		fmt.Println("Invalid secret: ", secret)
	}

	return new_target_id
}
