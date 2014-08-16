package main

import (
	"code.google.com/p/go-uuid/uuid"
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
	var kills int
	var alive bool

	err := db.QueryRow(`SELECT (team_id, user_role, kills, alive) FROM dm_user_game_mapping WHERE user_id = $1 AND game_id = $2`, userId, gameId).Scan(&teamId, &userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return &GameMapping{userId, gameId, teamId, userRole, kills, alive}, nil
}

func (user *User) JoinGame(gameId uuid.UUID) (*GameMapping, *ApplicationError) {
	var teamId uuid.UUID
	var userRole string
	var kills int
	var alive bool

	err := db.QueryRow(`INSERT INTO dm_user_game_mapping (user_id, game_id) VALUES ($1, $2) RETURNING team_id, user_role, kills, alive`, user.UserId, gameId).Scan(&teamId, &userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return &GameMapping{user.UserId, gameId, teamId, userRole, kills, alive}, nil
}

func (gameMapping *GameMapping) JoinTeam(teamId uuid.UUID) *ApplicationError {

	res, err := db.Exec(`UPDATE dm_user_game_mapping SET team_id = $1 WHERE user_id = $2 AND game_id = $3`, teamId, gameMapping.UserId, gameMapping.GameId)
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

	res, err := db.Exec(`UPDATE dm_user_game_mapping SET user_role = $1 WHERE user_id = $2 AND game_id = $3`, role, gameMapping.UserId, gameMapping.GameId)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}
	return nil
}
