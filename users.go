package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
)

type User struct {
	User_id string
	Email   string
	Secret  string
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
	returnUser := &User{id, email, secret}
	return returnUser
}

func (user *User) KillTarget(secret string) {

	logged_in_user := user.User_id

	row := db.QueryRow(`SELECT secret FROM dm_users WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user)

	var target_secret string

	row.Scan(&target_secret)

	if secret == target_secret {

		row = db.QueryRow(`SELECT target_id FROM dm_user_targets WHERE user_id = (SELECT target_id FROM dm_user_targets where user_id = $1)`, logged_in_user)

		var new_target_id string
		row.Scan(&new_target_id)
		row = db.QueryRow(`DELETE FROM dm_user_targets WHERE user_id = (SELECT target_id from dm_user_targets WHERE user_id = $1)`, logged_in_user)
		row = db.QueryRow(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2`, new_target_id, logged_in_user)

	} else {
		fmt.Println("Invalid secret\n", target_secret, "\n")
	}
}
