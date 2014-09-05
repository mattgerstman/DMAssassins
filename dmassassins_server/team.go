package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
)

type Team struct {
	TeamId   uuid.UUID `json:"team_id"`
	GameId   uuid.UUID `json:"game_id"`
	TeamName string    `json:"team_name"`
}

func GetTeamById(teamId uuid.UUID) (*Team, *ApplicationError) {
	var gameIdBuffer, teamName string
	err := db.QueryRow(`SELECT team_id, team_name FROM dm_teams WHERE game_id = $1 ORDER BY team_name`, teamId.String()).Scan(&gameIdBuffer, &teamName)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.Parse(gameIdBuffer)
	return &Team{teamId, gameId, teamName}, nil
}

func (game *Game) GetTeams() ([]*Team, *ApplicationError) {
	rows, err := db.Query(`SELECT team_id, team_name FROM dm_teams WHERE game_id = $1 ORDER BY team_name`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var teams []*Team
	for rows.Next() {
		var teamIdBuffer sql.NullString
		var teamName string

		err = rows.Scan(&teamIdBuffer, &teamName)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		teamId := uuid.Parse(teamIdBuffer.String)
		team := &Team{teamId, game.GameId, teamName}
		teams = append(teams, team)
	}
	return teams, nil
}

func (game *Game) CreateTeam(teamName string) (*Team, *ApplicationError) {
	teamId := uuid.NewUUID()

	res, err := db.Exec(`INSERT INTO dm_teams (team_id, game_id, team_name) VALUES ($1,$2,$3)`, teamId.String(), game.GameId.String(), teamName)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr := WereRowsAffected(res)
	if appErr != nil {
		return nil, appErr
	}
	return &Team{teamId, game.GameId, teamName}, nil
}

func DeleteTeam(teamId uuid.UUID) *ApplicationError {

	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare the statement to remove members of the team
	removeMembers, err := db.Prepare(`UPDATE dm_user_game_mapping SET team_id = null WHERE team_id = $1`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to remove members of the team
	_, err = tx.Stmt(removeMembers).Exec(teamId.String())
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare the statement to delete the team
	deleteTeam, err := db.Prepare(`DELETE from dm_teams WHERE team_id = $1`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to delete the team
	res, err := tx.Stmt(deleteTeam).Exec(teamId.String())
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr := WereRowsAffected(res)
	if appErr != nil {
		return appErr
	}
	return nil
}

func (user *User) GetTeam(gameId uuid.UUID) (*Team, *ApplicationError) {
	gameMapping, appErr := GetGameMapping(user.UserId, gameId)
	if appErr != nil {
		return nil, appErr
	}
	return GetTeamById(gameMapping.TeamId)
}

func (user *User) JoinTeam(teamId uuid.UUID) (*GameMapping, *ApplicationError) {
	var gameIdBuffer sql.NullString
	err := db.QueryRow(`SELECT game_id FROM dm_teams WHERE team_id = $1`, teamId.String()).Scan(&gameIdBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.Parse(gameIdBuffer.String)
	if gameId == nil {
		msg := "Invalid Team Id: " + teamId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeDatabase)
	}

	var userRole string
	var kills int
	var alive bool

	err = db.QueryRow(`UPDATE dm_user_game_mapping SET team_id = $1 WHERE game_id = $2 and user_id = $3 RETURNING user_role, kills, alive`, teamId.String(), gameId.String(), user.UserId.String()).Scan(&userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return &GameMapping{user.UserId, gameId, teamId, userRole, kills, alive}, nil

}

func (user *User) LeaveTeam(teamId uuid.UUID) (*GameMapping, *ApplicationError) {
	var gameIdBuffer sql.NullString
	err := db.QueryRow(`SELECT game_id FROM dm_teams WHERE team_id = $1`, teamId.String()).Scan(&gameIdBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.Parse(gameIdBuffer.String)
	if gameId == nil {
		msg := "Invalid Team Id: " + teamId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeDatabase)
	}

	var userRole string
	var kills int
	var alive bool

	err = db.QueryRow(`UPDATE dm_user_game_mapping SET team_id = null WHERE game_id = $1 and user_id = $2 RETURNING user_role, kills, alive`, gameId.String(), user.UserId.String()).Scan(&userRole, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return &GameMapping{user.UserId, gameId, nil, userRole, kills, alive}, nil
}

func (team *Team) Rename(newName string) (*Team, *ApplicationError) {
	res, err := db.Exec(`UPDATE dm_teams SET team_name = $1 WHERE team_id = $2`, newName, team.TeamId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr := WereRowsAffected(res)
	if appErr != nil {
		return nil, appErr
	}

	team.TeamName = newName
	return team, nil
}