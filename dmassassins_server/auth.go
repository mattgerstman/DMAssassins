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
	RoleSuperAdmin = 4
	RoleAdmin      = 3
	RoleCaptain    = 2
	RoleUser       = 1
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

	if dbToken != token {
		apiFacebookId, appErr := GetFacebookIdFromToken(token)
		if appErr != nil {
			return appErr
		}
		if apiFacebookId != user.FacebookId {
			msg := "Permission Denied"
			err := errors.New("Permission Denied")
			return NewApplicationError(msg, err, ErrCodePermissionDenied)
		}

	}

	return nil
}

func compareRole(role string, roleId int) (greaterThanOrEqualTo bool) {
	var roles = map[string]int{
		"dm_super_admin": RoleSuperAdmin,
		"dm_admin":       RoleAdmin,
		"dm_captain":     RoleCaptain,
		"dm_user":        RoleUser,
	}
	return roles[role] >= roleId
}

func RequiresUser(r *http.Request) (appErr *ApplicationError) {
	role, teamId, userId, appErr := getRoleFromRequest(r)
	if appErr != nil {
		return appErr
	}


	vars := mux.Vars(r)
	reqUserId := uuid.Parse(vars["user_id"])
	if (uuid.Equal(userId, reqUserId)) || (reqUserId == nil) {
		return nil
	}

	theyAre, appErr := isTeamCaptain(role, teamId, r)
	if appErr != nil {
		return appErr
	}
	if theyAre {
		return nil
	}

	msg := "Permission Denied"
	err := errors.New("Permission Denied")
	return NewApplicationError(msg, err, ErrCodePermissionDenied)

}

func isTeamCaptain(role string, teamId uuid.UUID, r *http.Request) (isRightCaptain bool, appErr *ApplicationError) {

	if compareRole(role, RoleAdmin) {
		return true, nil
	}

	if !compareRole(role, RoleCaptain) {
		return false, nil
	}

	vars := mux.Vars(r)
	userId := uuid.Parse(vars["user_id"])
	gameId := uuid.Parse(vars["game_id"])

	GameMapping, appErr := GetGameMapping(userId, gameId)
	if appErr != nil {
		return false, appErr
	}
	reqTeamId := GameMapping.TeamId

	if (uuid.Equal(teamId, reqTeamId)) || (reqTeamId == nil) {
		return true, nil
	}
	return false, nil

}

func RequiresCaptain(r *http.Request) (appErr *ApplicationError) {
	role, teamId, _, appErr := getRoleFromRequest(r)
	if appErr != nil {
		return appErr
	}

	theyAre, appErr := isTeamCaptain(role, teamId, r)
	if appErr != nil {
		return appErr
	}
	if theyAre {
		return nil
	}

	msg := "Permission Denied"
	err := errors.New("Permission Denied")
	return NewApplicationError(msg, err, ErrCodePermissionDenied)

}

func RequiresAdmin(r *http.Request) (appErr *ApplicationError) {
	role, _, _, appErr := getRoleFromRequest(r)
	if appErr != nil {
		return appErr
	}
	if compareRole(role, RoleAdmin) {
		return nil
	}

	msg := "Permission Denied"
	err := errors.New("Permission Denied")
	return NewApplicationError(msg, err, ErrCodePermissionDenied)

}

func RequiresSuperAdmin(r *http.Request) (appErr *ApplicationError) {
	role, _, _, appErr := getRoleFromRequest(r)
	if appErr != nil {
		return appErr
	}
	if compareRole(role, RoleSuperAdmin) {
		return nil
	}

	msg := "Permission Denied"
	err := errors.New("Permission Denied")
	return NewApplicationError(msg, err, ErrCodePermissionDenied)
}

func getRoleFromRequest(r *http.Request) (userRole string, teamId uuid.UUID, userId uuid.UUID, appErr *ApplicationError) {
	userId, token, appErr := GetBasicAuth(r)
	if appErr != nil {
		return "", nil, nil, appErr
	}
	vars := mux.Vars(r)
	gameId := uuid.Parse(vars["game_id"])

	var facebookId, facebookToken, teamIdBuffer string
	err := db.QueryRow(`SELECT user.facebook_id, user.facebook_token, game.user_role, game.team_id FROM dm_users AS user, dm_user_game_mapping AS game WHERE user.user_id = game.user_id AND game.user_id = $1 AND (game.game_id = $2 OR user.user_role == 'dm_super_admin')`, userId.String(), gameId.String()).Scan(&facebookId, &facebookToken, &userRole, &teamIdBuffer)
	if err != nil {
		return "", nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId = uuid.Parse(teamIdBuffer)
	if facebookToken != token {
		apiFacebookId, appErr := GetFacebookIdFromToken(token)
		if appErr != nil {
			return "", nil, nil, appErr
		}
		if apiFacebookId != facebookId {
			return "", nil, nil, NewApplicationError("Invalid Token", err, ErrCodeInvalidFBToken)
		}
	}
	return userRole, teamId, userId, nil
}
