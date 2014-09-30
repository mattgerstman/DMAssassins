package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"encoding/base64"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

const (
	RoleSuperAdmin = 4
	RoleAdmin      = 3
	RoleCaptain    = 2
	RoleUser       = 1
)

// Decode the Basic Auth header and return a userId and token to validate
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

// Requires that a user is logged in
func RequiresLogin(r *http.Request) (appErr *ApplicationError) {
	userId, token, appErr := GetBasicAuth(r)
	if appErr != nil {
		return appErr
	}

	user, appErr := GetUserById(userId)
	if appErr != nil {
		return appErr
	}

	dbToken, appErr := user.GetToken()
	if appErr != nil {
		return appErr
	}

	appErr = validateFacebookToken(dbToken, token, user.FacebookId)

	if appErr != nil {
		return appErr
	}

	return nil
}

// Compare two user roles by their int values
func CompareRole(role string, roleId int) (greaterThanOrEqualTo bool) {
	var roles = map[string]int{
		"dm_super_admin": RoleSuperAdmin,
		"dm_admin":       RoleAdmin,
		"dm_captain":     RoleCaptain,
		"dm_user":        RoleUser,
	}
	return roles[role] >= roleId
}

// Requires the same user, captain for that team, or admin for that game
func RequiresUser(r *http.Request) (role string, appErr *ApplicationError) {
	role, teamId, userId, appErr := getRoleFromRequest(r)
	if appErr != nil {
		return role, appErr
	}

	vars := mux.Vars(r)
	reqUserId := uuid.Parse(vars["user_id"])

	// If the userId's are equal or we're not requesting a user return no error
	if (uuid.Equal(userId, reqUserId)) || (reqUserId == nil) {
		return role, nil
	}

	// Check if the auth token is for a team captain for theuse given
	appErr = isTeamCaptain(role, teamId, r)
	if appErr != nil {
		return role, appErr
	}
	return role, nil
}

// Standard permission denied application error
func GetPermissionDeniedAppErr() (appErr *ApplicationError) {
	msg := "Permission Denied"
	err := errors.New("Permission Denied")
	return NewApplicationError(msg, err, ErrCodePermissionDenied)
}

// Check if a role/teamId match to be the team captain for the user_id/game_id in the request
func isTeamCaptain(role string, teamId uuid.UUID, r *http.Request) (appErr *ApplicationError) {

	if CompareRole(role, RoleAdmin) {
		return nil
	}

	if !CompareRole(role, RoleCaptain) {
		return GetPermissionDeniedAppErr()
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	if userId == nil {
		return nil
	}

	gameId := uuid.Parse(vars["game_id"])

	// Get the game mapping for the necessary user
	GameMapping, appErr := GetGameMapping(userId, gameId)
	if appErr != nil {
		return appErr
	}
	// Team id we need to validate against
	reqTeamId := GameMapping.TeamId

	if (uuid.Equal(teamId, reqTeamId)) || (reqTeamId == nil) {
		return nil
	}
	return GetPermissionDeniedAppErr()

}

// Requires the user is a team captain
func RequiresCaptain(r *http.Request) (role string, appErr *ApplicationError) {
	role, teamId, _, appErr := getRoleFromRequest(r)
	if appErr != nil {
		return role, appErr
	}

	appErr = isTeamCaptain(role, teamId, r)
	if appErr != nil {
		return role, appErr
	}

	return role, nil
}

// Requires the user is a game admin
func RequiresAdmin(r *http.Request) (role string, appErr *ApplicationError) {
	role, _, _, appErr = getRoleFromRequest(r)
	if appErr != nil {
		return role, appErr
	}
	if !CompareRole(role, RoleAdmin) {
		return role, GetPermissionDeniedAppErr()
	}

	return role, nil

}

// Requires the user is Matt Gerstman
func RequiresSuperAdmin(r *http.Request) (role string, appErr *ApplicationError) {
	role, _, _, appErr = getRoleFromRequest(r)
	if appErr != nil {
		return role, appErr
	}
	if !CompareRole(role, RoleSuperAdmin) {
		return role, GetPermissionDeniedAppErr()
	}
	return role, nil
}

// Validates a db token/facebook id against the given token
func validateFacebookToken(facebookToken, token, facebookId string) *ApplicationError {
	if facebookToken == token {
		return nil
	}
	apiFacebookId, appErr := GetFacebookIdFromToken(token)
	if appErr != nil {
		return appErr
	}
	if apiFacebookId != facebookId {
		msg := "Invalid Token"
		err := errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidFBToken)
	}
	return nil
}

func getRoleFromRequest(r *http.Request) (userRole string, teamId uuid.UUID, userId uuid.UUID, appErr *ApplicationError) {
	userId, token, appErr := GetBasicAuth(r)
	if appErr != nil {
		return "", nil, nil, appErr
	}
	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])

	var facebookId, facebookToken string
	var teamIdBuffer sql.NullString
	err := db.QueryRow(`SELECT dm_users.facebook_id, dm_users.facebook_token, game.user_role, game.team_id FROM dm_users, dm_user_game_mapping AS game WHERE dm_users.user_id = game.user_id AND game.user_id = $1 AND (game.game_id = $2 OR game.user_role = 'dm_super_admin')`, userId.String(), gameId.String()).Scan(&facebookId, &facebookToken, &userRole, &teamIdBuffer)
	if err != nil {
		return "", nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId = uuid.Parse(teamIdBuffer.String)
	appErr = validateFacebookToken(facebookToken, token, facebookId)
	if appErr != nil {
		return "", nil, nil, appErr
	}

	return userRole, teamId, userId, nil
}
