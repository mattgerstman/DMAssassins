package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"fmt"
	"github.com/getsentry/raven-go"
	fb "github.com/huandu/facebook"
)

// Returns facebook app
func getFbApp() *fb.App {
	fb.Version = "v2.2"
	var app = fb.New(Config.FBAppId, Config.FBAppSecret)
	app.RedirectUri = "http://playassassins.com"
	return app
}

// Returns an authenticated facebook session with app id/secret
func getFacebookSession(token string) (fbSession *fb.Session) {

	app := getFbApp()
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
		email = `none-provided@playassassins.com`
		appErr := NewApplicationError("Internal Error", err, ErrCodeInvalidFBToken)
		extra := make(map[string]interface{})
		extra[`facebook`] = res
		extra[`facebook_id`] = facebookId
		LogWithSentry(appErr, map[string]string{}, raven.WARNING, extra, nil)
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
	username := firstName + ` ` + lastName

	// Set up user properties map, this will be inserted with the user
	properties := make(map[string]string)
	properties["facebook"] = "https://facebook.com/" + facebookId

	picture := "https://graph.facebook.com/" + facebookId + "/picture"
	properties["photo"] = picture + "?width=1000"
	properties["photo_thumb"] = picture + "?width=300&height=300"

	properties["first_name"] = firstName
	properties["last_name"] = lastName
	properties["allow_email"] = "true"
	properties["allow_posts"] = "true"

	// Create user
	user, appErr = NewUser(username, email, facebookId, properties)
	if appErr != nil {
		return nil, appErr
	}

	appErr = user.UpdateToken(facebookToken)
	if appErr != nil {
		sentryUser := NewSentryUser(user)
		LogWithSentry(appErr, map[string]string{"user_id": user.UserId.String()}, raven.WARNING, nil, sentryUser)
	}

	_, appErr = user.SendUserWelcomeEmail()
	if appErr != nil {
		sentryUser := NewSentryUser(user)
		LogWithSentry(appErr, map[string]string{"user_id": user.UserId.String()}, raven.WARNING, nil, sentryUser)
	}

	go user.StoreUserFriends()

	return user, nil

}

// wrapper for storeUserFreinds to log any errors
func (user *User) StoreUserFriends() {
	appErr := user.storeUserFriends()
	if appErr != nil {
		sentryUser := NewSentryUser(user)
		LogWithSentry(appErr, map[string]string{"user_id": user.UserId.String()}, raven.WARNING, nil, sentryUser)
	}
}

// store a user's friends
func (user *User) storeUserFriends() (appErr *ApplicationError) {
	friends, appErr := user.GetFacebookFriends()
	if appErr != nil {
		return appErr
	}

	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// prepare statement to delete friends already in the db
	deleteFriends, err := tx.Prepare(`DELETE FROM dm_friends where facebook_id = $1 OR friend_id = $1`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// execute statement to delete friends already in the db
	_, err = tx.Stmt(deleteFriends).Exec(user.FacebookId)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// insert all friend pairs
	for _, friendData := range friends {
		friendId := friendData[`id`]

		// Prepare statement to insert friend
		insertFriend, err := tx.Prepare(`INSERT INTO dm_friends (facebook_id, friend_id) VALUES ($1, $2), ($2, $1)`)
		if err != nil {
			tx.Rollback()
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// execute statement to insert friend
		_, err = tx.Stmt(insertFriend).Exec(user.FacebookId, friendId)
		if err != nil {
			tx.Rollback()
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
	}
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Get mutual friends for two users
func (user *User) GetMutualFriends(targetId string) (friends []map[string]string, count int, appErr *ApplicationError) {
	err := db.QueryRow(`SELECT count(*) FROM dm_users where facebook_id IN (SELECT user_friends.friend_id from (SELECT friend_id FROM dm_friends WHERE facebook_id = $1) user_friends, LATERAL ( SELECT friend_id FROM dm_friends WHERE facebook_id = $2) target_friends)`, user.FacebookId, targetId).Scan(&count)
	if err != nil {
		return nil, 0, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// If there are no friends skip the second sectin
	if count == 0 {
		return nil, 0, nil
	}

	rows, err := db.Query(`SELECT facebook_id, username FROM dm_users where facebook_id IN (SELECT user_friends.friend_id from (SELECT friend_id FROM dm_friends WHERE facebook_id = $1) user_friends, LATERAL ( SELECT friend_id FROM dm_friends WHERE facebook_id = $2) target_friends) LIMIT 5`, targetId, user.FacebookId)
	if err != nil {
		return nil, 0, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	for rows.Next() {
		var friendId string
		var friendName string
		err := rows.Scan(&friendId, &friendName)
		if err != nil {
			return nil, 0, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		friend := make(map[string]string)
		friend[`facebook_id`] = friendId
		friend[`user_name`] = friendName
		friends = append(friends, friend)
	}

	return friends, count, nil
}

func ExtendToken(facebookToken string) (longLivedToken string, appErr *ApplicationError) {
	// Query facebook session to make extend token
	app := getFbApp()
	longLivedToken, _, err := app.ExchangeToken(facebookToken)
	if err != nil {
		return "", NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}
	return longLivedToken, nil

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
	facebookId, err := session.User()
	if err != nil {
		return "", NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}

	return facebookId, nil
}

// get all of a user's friends playing assassins
func (user *User) GetFacebookFriends() (friends []map[string]string, appErr *ApplicationError) {
	res, err := fb.Get("/"+user.FacebookId+"/friends/", fb.Params{"access_token": Config.FBAccessToken})
	if err != nil {
		return nil, NewApplicationError("Error Contacting Facebook", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField(`data`, &friends)
	if err != nil {
		return nil, NewApplicationError("Error Contacting Facebook", err, ErrCodeInvalidFBToken)
	}

	return friends, nil
}

// get a user's facebook photos
func (user *User) GetFacebookPhotos() (photos []interface{}, appErr *ApplicationError) {
	token, appErr := user.GetToken()
	if appErr != nil {
		return nil, appErr
	}

	res, err := fb.Get("/"+user.FacebookId+"/photos/", fb.Params{"access_token": token})
	if err != nil {
		return nil, NewApplicationError("Error Contacting Facebook", err, ErrCodeInvalidFBToken)
	}

	err = res.DecodeField(`data`, &photos)
	if err != nil {
		return nil, NewApplicationError("Error Contacting Facebook", err, ErrCodeInvalidFBToken)
	}

	return photos, nil
}

func (game *Game) FacebookPost(message string) (appErr *ApplicationError) {

	pageId, appErr := game.GetGameProperty(`game_page_id`)
	if appErr != nil {
		return appErr
	}

	accessToken, appErr := game.GetGameProperty(`game_page_access_token`)
	if appErr != nil {
		return appErr
	}

	session := getFacebookSession(accessToken)
	res, err := session.Post(`/`+pageId+`/feed`, fb.Params{"message": message, "access_token": accessToken})
	fmt.Println(res)
	fmt.Println(err)
	if err != nil {
		return NewApplicationError("Invalid Facebook Token", err, ErrCodeInvalidFBToken)
	}

	return nil
}
