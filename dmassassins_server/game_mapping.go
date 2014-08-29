package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
)

type GameMapping struct {
	UserId   uuid.UUID `json:"user_id"`
	GameId   uuid.UUID `json:"game_id"`
	TeamId   uuid.UUID `json:"team_id"`
	UserRole string    `json:"user_role"`
	kills    int       `json:"kills"`
	alive    bool      `json:"alive"`
}

func GetGameMapping(userId, gameId uuid.UUID) (*GameMapping, *ApplicationError) {
	var teamId uuid.UUID
	var userRole string
	var teamIdBuffer sql.NullString
	var kills int
	var alive bool

	err := db.QueryRow(`SELECT team_id, user_role, kills, alive FROM dm_user_game_mapping WHERE user_id = $1 AND game_id = $2`, userId.String(), gameId.String()).Scan(&teamIdBuffer, &userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId = uuid.Parse(teamIdBuffer.String)
	return &GameMapping{userId, gameId, teamId, userRole, kills, alive}, nil
}

func CheckPassword(gameId uuid.UUID, testPassword string) *ApplicationError {
	var gamePasswordBuffer sql.NullString
	err := db.QueryRow(`SELECT game_password FROM dm_games WHERE game_id = $1`, gameId.String()).Scan(&gamePasswordBuffer)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	if !gamePasswordBuffer.Valid || (gamePasswordBuffer.String == "") {
		return nil
	}
	gamePassword := gamePasswordBuffer.String
	if gamePassword == testPassword {
		return nil
	}
	msg := "Invalid Game Password: " + testPassword
	err = errors.New(msg)
	return NewApplicationError(msg, err, ErrCodeInvalidGamePassword)

}

func (user *User) JoinGame(gameId uuid.UUID, gamePassword string) (*GameMapping, *ApplicationError) {
	var teamId uuid.UUID
	var userRole string
	var teamIdBuffer sql.NullString
	var kills int
	var alive bool

	appErr := CheckPassword(gameId, gamePassword)
	if appErr != nil {
		return nil, appErr
	}

	err := db.QueryRow(`INSERT INTO dm_user_game_mapping (user_id, game_id) VALUES ($1, $2) RETURNING team_id, user_role, kills, alive`, user.UserId.String(), gameId.String()).Scan(&teamIdBuffer, &userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId = uuid.Parse(teamIdBuffer.String)
	return &GameMapping{user.UserId, gameId, teamId, userRole, kills, alive}, nil
}

func (gameMapping *GameMapping) JoinTeam(teamId uuid.UUID) *ApplicationError {

	res, err := db.Exec(`UPDATE dm_user_game_mapping SET team_id = $1 WHERE user_id = $2 AND game_id = $3`, teamId.String(), gameMapping.UserId.String(), gameMapping.GameId)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}
	return nil
}

func (gameMapping *GameMapping) ChangeRole(role string) *ApplicationError {

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
