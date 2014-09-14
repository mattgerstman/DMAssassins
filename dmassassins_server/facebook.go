package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	fb "github.com/huandu/facebook"
	"github.com/polds/imgbase64"
)

// Returns an authenticated facebook session with app id/secret
// Need to move app id/secret to config file
func getFacebookSession(token string) (fbSession *fb.Session) {

	fb.Version = "v2.0"
	var app = fb.New(Config.FBAppId, Config.FBAppSecret)
	app.RedirectUri = "http://dmassassins.com"

	session := app.Session(token)
	return session
}

// Creates a user from a facebook_auth token

func CreateUserFromFacebookToken(facebookToken string) (user *User, appErr *ApplicationError) {

	session := getFacebookSession(facebookToken)
	res, err := session.Get("/me/", fb.Params{})
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}
	var firstName, lastName, email string
	var facebook, facebookId string

	// Decodes all the fields from the facebook session and puts them in variables to be used with the new user
	err = res.DecodeField("first_name", &firstName)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField("last_name", &lastName)
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

	err = res.DecodeField("id", &facebookId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
	}

	// Username's are a concatination of first/last names
	username := firstName + lastName

	// Set up user properties map, this will be inserted with the user
	properties := make(map[string]string)
	properties["Facebook"] = "https://facebook.com/" + facebookId

	picture := "https://graph.facebook.com/" + facebookId + "/picture"
	photo := picture + "?width=1000"
	imgbase64.SetDefaultImage(photo)
	img := imgbase64.FromRemote(photo)
	properties["photo"] = img

	photo_thumb := picture + "?width=300&height=300"
	imgbase64.SetDefaultImage(photo_thumb)
	img_thumb := imgbase64.FromRemote(photo_thumb)
	properties["photo_thumb"] = img_thumb

	properties["first_name"] = firstName
	properties["last_name"] = lastName

	// Create user
	user, appErr = NewUser(username, email, facebookId, properties)
	if appErr != nil {
		return nil, appErr
	}

	return user, nil

}

// Get a user from the db by it's facebook_id, confirms that the id matches the id in the token
// If there is no user in the DB with that facebook_id add them
func getUserFromFacebookId(facebookId, facebookToken string) (user *User, appErr *ApplicationError) {
	var userId uuid.UUID
	var userIdBuffer sql.NullString

	// See if we already have the facebook id in the database
	err := db.QueryRow(`SELECT user_id FROM dm_users WHERE facebook_id = $1`, facebookId).Scan(&userIdBuffer)
	userId = uuid.Parse(userIdBuffer.String)

	switch {
	// If we don't have the user create it
	case err == sql.ErrNoRows:
		return CreateUserFromFacebookToken(facebookToken)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)

	}
	// Now lets confirm that the db and the facebook id in the token match
	testId, appErr := GetFacebookIdFromToken(facebookToken)
	if appErr != nil {
		return nil, appErr
	}
	if testId != facebookId {
		return nil, NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}

	// Get the user object
	user, appErr = GetUserById(userId)
	if appErr != nil {
		return nil, appErr
	}

	// Set the user's facebook token to the most recent one
	appErr = user.UpdateToken(facebookToken)
	if appErr != nil {
		return nil, appErr
	}

	// return the user
	return user, nil

}

// Returns a user based on facebook_id and facebook_token, if no user exists in the db one will be created
func GetUserFromFacebookData(facebookId, facebookToken string) (user *User, appErr *ApplicationError) {

	var userId uuid.UUID
	var userIdBuffer sql.NullString
	// See if we have a user with the given facebook_id/facebook_token in the db
	err := db.QueryRow(`SELECT user_id FROM dm_users WHERE facebook_id = $1 AND facebook_token = $2`, facebookId, facebookToken).Scan(&userIdBuffer)
	userId = uuid.Parse(userIdBuffer.String)

	switch {
	// If we have no user in the db check just the id and see if it's in the database
	case err == sql.ErrNoRows:
		return getUserFromFacebookId(facebookId, facebookToken)
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	//If we have a user_id from the db just return a user
	return GetUserById(userId)
}

// Returns a user's facebook_id from their token
func GetFacebookIdFromToken(token string) (facebookId string, appErr *ApplicationError) {

	// Query facebook session to make sure token is valid
	session := getFacebookSession(token)
	res, err := session.Get("/debug_token", fb.Params{"input_token": token, "access_token":Config.FBAccessToken})
	if err != nil {
		return "", NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}

	// Decode the facebook Id from the sessoin data
	err = res.DecodeField("data.user_id", &facebookId)
	if err != nil {
		return "", NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}
	return facebookId, nil

}
