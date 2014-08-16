package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/base64"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const (
	RoleAdmin   = 3
	RoleCaptain = 2
	RoleUser    = 1
)

func GetBasicAuth(r *http.Request) (userId uuid.UUID, token string, appErr *ApplicationError) {
	// Get the Authorization header.
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		msg := "Missing Header: Authorization"
		err := errors.New(msg)
		return nil, "", NewApplicationError(msg, err, ErrCodeMissingHeader)
	}

	// Check the Authorization type is "Basic"
	authComponents := strings.Split(authHeader, " ")
	if len(authComponents) != 2 {
		msg := "Invalid Header: Authorization"
		err := errors.New("Header was not formatted properly")
		return nil, "", NewApplicationError(msg, err, ErrCodeInvalidHeader)
	}
	if authComponents[0] != "Basic" {
		msg := "Invalid Header: Authorization"
		err := errors.New("Authorization type Basic was expected")
		return nil, "", NewApplicationError(msg, err, ErrCodeInvalidHeader)
	}

	// Base64 Decode the user-pass string.
	decoded, err := base64.StdEncoding.DecodeString(authComponents[1])
	if err != nil {
		return nil, "", NewApplicationError("Invalid Header: Authorization", err, ErrCodeInvalidHeader)
	}

	// Split apart the username and password.
	userTokenComponents := strings.Split(string(decoded), ":")
	if len(userTokenComponents) != 2 {
		msg := "Invalid Header: Authorization"
		err := errors.New("Header was not formatted properly")
		return nil, "", NewApplicationError(msg, err, ErrCodeInvalidHeader)
	}

	userId = uuid.Parse(userTokenComponents[0])
	if userId == nil {
		msg := "Invalid Header: Authorization"
		err := errors.New("UserId was not valid UUID")
		return nil, "", NewApplicationError(msg, err, ErrCodeInvalidHeader)
	}
	token = userTokenComponents[1]

	return userId, token, nil
}

func compareRole(role string, roleId int) bool {
	var roles = map[string]int{
		"dm_admin":   RoleAdmin,
		"dm_captain": RoleCaptain,
		"dm_user":    RoleUser,
	}
	return roles[role] >= roleId
}

func RequiresUser(r *http.Request) (bool, *ApplicationError) {
	role, appErr := getRoleFromHeaders(r)
	if appErr != nil {
		return false, appErr
	}
	return compareRole(role, RoleUser), nil

}

func RequiresCaptain(r *http.Request) (bool, *ApplicationError) {
	role, appErr := getRoleFromHeaders(r)
	if appErr != nil {
		return false, appErr
	}
	return compareRole(role, RoleCaptain), nil

}

func RequiresAdmin(r *http.Request) (bool, *ApplicationError) {
	role, appErr := getRoleFromHeaders(r)
	if appErr != nil {
		return false, appErr
	}
	return compareRole(role, RoleAdmin), nil

}

func getRoleFromHeaders(r *http.Request) (string, *ApplicationError) {
	userId, token, appErr := GetBasicAuth(r)
	if appErr != nil {
		return "", appErr
	}
	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])

	var facebookId, facebookToken, userRole string
	err := db.QueryRow(`SELECT user.facebook_id, user.facebook_token, game.user_role FROM dm_users AS user, dm_user_game_mapping AS game WHERE user.user_id = game.user_id AND game.user_id = $1 AND game.game_id = $2`, userId.String(), gameId.String()).Scan(&facebookId, &facebookToken, &userRole)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	if facebookToken != token {
		apiFacebookId, appErr := GetFacebookIdFromToken(token)
		if appErr != nil {
			return "", appErr
		}
		if apiFacebookId != facebookId {
			return "", NewApplicationError("Invalid Token", err, ErrCodeInvalidFBToken)
		}
	}
	return userRole, nil
}
