package main

import (
"code.google.com/p/go-uuid/uuid"
)

const (
	RoleAdmin   = 3
	RoleCaptain = 2
	RoleUser    = 1
)

func compareRole(role string, roleId int) bool {

	var roles = map[string]int{
		"dm_admin":   RoleAdmin,
		"dm_captain": RoleCaptain,
		"dm_user":    RoleUser,
	}
	return roles[role] >= roleId
}

func RequiresUser(userId, gameId uuid.UUID, token string) (bool, *ApplicationError) {
	role, appErr := getRoleFromToken(userId, gameId, token)
	if appErr != nil {
		return false, appErr
	}
	return compareRole(role, RoleUser), nil

}

func RequiresCaptain(userId, gameId uuid.UUID, token string) (bool, *ApplicationError) {
	role, appErr := getRoleFromToken(userId, gameId, token)
	if appErr != nil {
		return false, appErr
	}
	return compareRole(role, RoleCaptain), nil

}

func RequiresAdmin(userId, gameId uuid.UUID, token string) (bool, *ApplicationError) {
	role, appErr := getRoleFromToken(userId, gameId, token)
	if appErr != nil {
		return false, appErr
	}
	return compareRole(role, RoleAdmin), nil

}

func getRoleFromToken(userId, gameId uuid.UUID, token string) (string, *ApplicationError) {
	var facebookId, facebookToken, userRole string
	err := db.QueryRow(`SELECT user.facebook_id, user.facebook_token, game.user_role FROM dm_users AS user, dm_user_game_mapping AS game WHERE user.user_id = game.user_id AND game.user_id = $1 AND game.game_id = $2`, userId, gameId).Scan(&facebookId, &facebookToken, &userRole)
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
