package main

import (
	"database/sql"
	"errors"
	"fmt"
	fb "github.com/huandu/facebook"
)

func getFacebookSession(token string) *fb.Session {

	fb.Version = "v2.0"
	var app = fb.New("643600385736129", "73cbc95ae6de7a6c26b16318330f796a")
	app.RedirectUri = "http://dmassassins.com"

	session := app.Session(token)
	return session
}

func facebook(path string) (interface{}, *ApplicationError) {

	username := "Matt"
	var facebook_id, facebook_token string
	_ = db.QueryRow(`SELECT facebook_id, facebook_token FROM dm_users WHERE username = $1`, username).Scan(&facebook_id, &facebook_token)

	//path := "/" + facebook_id + "/friends/"

	session := getFacebookSession(facebook_token)
	res, _ := session.Get(path, fb.Params{})

	return res, nil
}

func CreateUserFromFacebookToken(facebook_token string) (interface{}, *ApplicationError) {

	session := getFacebookSession(facebook_token)
	res, err := session.Get("/me/", fb.Params{})
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}
	var first_name, last_name, email string
	var facebook, facebook_id string
	_ = res.DecodeField("first_name", &first_name)
	_ = res.DecodeField("last_name", &last_name)
	_ = res.DecodeField("email", &email)
	_ = res.DecodeField("link", &facebook)
	_ = res.DecodeField("id", &facebook_id)

	username := first_name + last_name

	properties := make(map[string]string)
	properties["Facebook"] = facebook

	picture := "https://graph.facebook.com/" + facebook_id
	properties["photo"] = picture + "?width=1000"
	properties["photo_thumb"] = picture + "?width=300&height=300"

	user, appErr := NewUser(username, email, "muggle", facebook_id, properties)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil

}

func getUserFromFacebookId(facebook_id, facebook_token string) (interface{}, *ApplicationError) {
	var user_id string
	err := db.QueryRow(`SELECT user_id FROM dm_users WHERE facebook_id = $1`, facebook_id).Scan(&user_id)
	switch {
	case err == sql.ErrNoRows:
		return CreateUserFromFacebookToken(facebook_token)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	case err == nil:
		test_id, appErr := GetFacebookIdFromToken(facebook_token)
		if appErr != nil {
			return nil, appErr
		}
		if test_id == facebook_id {
			return GetUserById(user_id)
		}
		return nil, NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)

	}
	err = errors.New("Unknown error in getUserFromFacebookId")
	return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)

}

func GetUserFromFacebookData(facebook_id, facebook_token string) (interface{}, *ApplicationError) {

	var user_id string
	// See if we have a user with the given facebook_id/facebook_token in the db
	err := db.QueryRow(`SELECT user_id FROM dm_users WHERE facebook_id = $1 AND facebook_token = $2`, facebook_id, facebook_token).Scan(&user_id)
	switch {
	// If we have no user in the db check just the id and see if it's in the database
	case err == sql.ErrNoRows:
		return getUserFromFacebookId(facebook_id, facebook_token)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	//If we have a user_id from the db just return a user
	return GetUserById(user_id)
}

func GetFacebookIdFromToken(token string) (interface{}, *ApplicationError) {

	session := getFacebookSession(token)

	res, _ := session.Get("/debug_token", fb.Params{"input_token": token})
	var user_id string
	err := res.DecodeField("data.user_id", &user_id)
	fmt.Println(err)
	fmt.Println(user_id)
	return user_id, nil

}
