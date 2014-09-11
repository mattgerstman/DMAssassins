package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"fmt"
)

type GameMapping struct {
	UserId   uuid.UUID `json:"user_id"`
	GameId   uuid.UUID `json:"game_id"`
	TeamId   uuid.UUID `json:"team_id"`
	UserRole string    `json:"user_role"`
	Secret   string    `json:"secret"`
	Kills    int       `json:"kills"`
	Alive    bool      `json:"alive"`
}

// Gets a game mapping from the database
func GetGameMapping(userId, gameId uuid.UUID) (gameMapping *GameMapping, appErr *ApplicationError) {
	var teamId uuid.UUID
	var userRole, secret string
	var teamIdBuffer sql.NullString
	var kills int
	var alive bool

	// Query the database
	err := db.QueryRow(`SELECT team_id, user_role, secret, kills, alive FROM dm_user_game_mapping WHERE user_id = $1 AND game_id = $2`, userId.String(), gameId.String()).Scan(&teamIdBuffer, &userRole, &secret, &kills, &alive)
	if err == sql.ErrNoRows {
		return nil, NewApplicationError("No Game Mappings", err, ErrCodeNotFoundGameMapping)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId = uuid.Parse(teamIdBuffer.String)
	return &GameMapping{userId, gameId, teamId, userRole, secret, kills, alive}, nil
}

// Adds a user to a game (creates a new game mapping)
func (user *User) JoinGame(gameId uuid.UUID, gamePassword string) (gameMapping *GameMapping, appErr *ApplicationError) {
	var teamId uuid.UUID
	var userRole string
	var teamIdBuffer sql.NullString
	var kills int
	var alive bool

	appErr = CheckPassword(gameId, gamePassword)
	if appErr != nil {
		return nil, appErr
	}

	secret, appErr := NewSecret()
	if appErr != nil {
		return nil, appErr
	}

	// Insert the GameMapping and get its default variables out of the database
	err := db.QueryRow(`INSERT INTO dm_user_game_mapping (user_id, game_id, secret) VALUES ($1, $2, $3) RETURNING team_id, user_role, kills, alive`, user.UserId.String(), gameId.String(), secret).Scan(&teamIdBuffer, &userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId = uuid.Parse(teamIdBuffer.String)
	return &GameMapping{user.UserId, gameId, teamId, userRole, secret, kills, alive}, nil
}

// Changes a user's role within a game
func (gameMapping *GameMapping) ChangeRole(role string) (appErr *ApplicationError) {

	res, err := db.Exec(`UPDATE dm_user_game_mapping SET user_role = $1 WHERE user_id = $2 AND game_id = $3`, role, gameMapping.UserId.String(), gameMapping.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}
	return nil
}

// Gets a user's assassins
func GetAssassin(targetId, gameId uuid.UUID) (assassin *User, appErr *ApplicationError) {

	var assassinIdBuffer string
	// Query the targets table
	err := db.QueryRow(`SELECT user_id from dm_user_targets WHERE target_id = $1 AND game_id = $2`, targetId.String(), gameId.String()).Scan(&assassinIdBuffer)
	if err == sql.ErrNoRows {
		msg := "Invalid target_id: " + targetId.String()
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundUserId)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	fmt.Println(assassinIdBuffer)
	// Grab assassin Id
	assassinId := uuid.Parse(assassinIdBuffer)

	// Get the user through the normal routes
	return GetUserById(assassinId)

}

// delete sthe actual game mapping from the db
func (gameMapping *GameMapping) delete() (appErr *ApplicationError) {
	res, err := db.Exec(`DELETE from dm_user_game_mapping WHERE user_id = $1 and game_id = $2`, gameMapping.UserId.String(), gameMapping.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	// the game mapping no longer exists so  set it to nil
	gameMapping = nil
	return nil
}

// Lets a user quit a game
func (gameMapping *GameMapping) LeaveGame(secret string) (appErr *ApplicationError) {
	// Get the user's assassin so we can have them "kill" their target
	fmt.Println("GetAssassin")
	assassin, appErr := GetAssassin(gameMapping.UserId, gameMapping.GameId)
	if appErr != nil && appErr.Code != ErrCodeNotFoundUserId {
		return appErr
	}

	if appErr == nil {
		fmt.Println("assassin Id" + assassin.UserId.String())
		_, appErr = assassin.KillTarget(gameMapping.GameId, secret, false)
		if appErr != nil {
			return appErr
		}
	}

	fmt.Println("delete game mapping")

	appErr = gameMapping.delete()
	if appErr != nil {
		return appErr
	}

	return nil
}

// Gets an arbitrary game for a user to start off with
func (user *User) GetArbitraryGameMapping() (gameMapping *GameMapping, appErr *ApplicationError) {

	var userRole, secret, gameIdBuffer string
	var teamIdBuffer sql.NullString
	var kills int
	var alive bool

	// Query the database
	err := db.QueryRow(`SELECT game_id, team_id, user_role, secret, kills, alive FROM dm_user_game_mapping WHERE user_id = $1 ORDER BY user_role, alive LIMIT 1`, user.UserId.String()).Scan(&gameIdBuffer, &teamIdBuffer, &userRole, &secret, &kills, &alive)
	if err == sql.ErrNoRows {
		return nil, NewApplicationError("No Games", err, ErrCodeNoGameMappings)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId := uuid.Parse(teamIdBuffer.String)
	gameId := uuid.Parse(gameIdBuffer)
	return &GameMapping{user.UserId, gameId, teamId, userRole, secret, kills, alive}, nil

}
