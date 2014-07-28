package main

import (
	"database/sql"
	fb "github.com/huandu/facebook"
)

// Returns an authenticated facebook session with app id/secret
// Need to move app id/secret to config file
func getFacebookSession(token string) *fb.Session {

	fb.Version = "v2.0"
	var app = fb.New("643600385736129", "73cbc95ae6de7a6c26b16318330f796a")
	app.RedirectUri = "http://dmassassins.com"

	session := app.Session(token)
	return session
}

// Creates a user from a facebook_auth token

func CreateUserFromFacebookToken(facebook_token string) (*User, *ApplicationError) {

	session := getFacebookSession(facebook_token)
	res, err := session.Get("/me/", fb.Params{})
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}
	var first_name, last_name, email string
	var facebook, facebook_id string

	// Decodes all the fields from the facebook session and puts them in variables to be used with the new user
	err = res.DecodeField("first_name", &first_name)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField("last_name", &last_name)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField("email", &email)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField("link", &facebook)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField("id", &facebook_id)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	// Username's are a concatination of first/last names
	username := first_name + last_name

	// Set up user properties map, this will be inserted with the user
	properties := make(map[string]string)
	properties["Facebook"] = facebook

	picture := "https://graph.facebook.com/" + facebook_id
	properties["photo"] = picture + "?width=1000"
	properties["photo_thumb"] = picture + "?width=300&height=300"

	properties["first_name"] = first_name
	properties["last_name"] = last_name

	user, appErr := NewUser(username, email, "muggle", facebook_id, properties)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil

}

// Get a user from the db by it's facebook_id, confirms that the id matches the id in the token
// If there is no user in the DB with that facebook_id add them
func getUserFromFacebookId(facebook_id, facebook_token string) (*User, *ApplicationError) {
	var user_id string
	err := db.QueryRow(`SELECT user_id FROM dm_users WHERE facebook_id = $1`, facebook_id).Scan(&user_id)
	switch {
	case err == sql.ErrNoRows:
		return CreateUserFromFacebookToken(facebook_token)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)

	}
	test_id, appErr := GetFacebookIdFromToken(facebook_token)
	if appErr != nil {
		return nil, appErr
	}
	if test_id != facebook_id {
		return nil, NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)	
	}
	return GetUserById(user_id)
	

}

// Returns a user based on facebook_id and facebook_token, if no user exists in the db one will be created
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

// Returns a user's facebook_id from their token
func GetFacebookIdFromToken(token string) (interface{}, *ApplicationError) {

	session := getFacebookSession(token)
	res, err := session.Get("/debug_token", fb.Params{"input_token": token})
	if err != nil {
		return nil, NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}

	var user_id string
	err = res.DecodeField("data.user_id", &user_id)
	if err != nil {
		return nil, NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}	
	return user_id, nil

}
