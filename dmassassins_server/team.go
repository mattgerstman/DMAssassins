package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
)

const (
	MinUsersPerTeam = 5
)

type Team struct {
	TeamId   uuid.UUID `json:"team_id"`
	GameId   uuid.UUID `json:"game_id"`
	TeamName string    `json:"team_name"`
}

// Gets a team by it's team_id
func GetTeamById(teamId uuid.UUID) (team *Team, appErr *ApplicationError) {
	var gameIdBuffer, teamName string
	err := db.QueryRow(`SELECT team_id, team_name FROM dm_teams WHERE team_id = $1 ORDER BY team_name`, teamId.String()).Scan(&gameIdBuffer, &teamName)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	gameId := uuid.Parse(gameIdBuffer)
	return &Team{teamId, gameId, teamName}, nil
}

// Gets a team for a user by a game_id
func (user *User) GetTeamByGameId(gameId uuid.UUID) (team *Team, appErr *ApplicationError) {
	var teamIdBuffer, teamName string
	err := db.QueryRow(`SELECT team_id, team_name FROM dm_teams WHERE team_id = (SELECT team_id FROM dm_user_game_mapping WHERE user_id = $1 AND game_id = $2)`, user.UserId.String(), gameId.String()).Scan(&teamIdBuffer, &teamName)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	teamId := uuid.Parse(teamIdBuffer)
	user.Properties["team"] = teamName
	return &Team{teamId, gameId, teamName}, nil
}

// Get a mapping of TeamId to Team
func (game *Game) GetTeamsMap() (teams map[string]*Team, appErr *ApplicationError) {
	// Query Db
	rows, err := db.Query(`SELECT team_id, team_name FROM dm_teams WHERE game_id = $1 ORDER BY team_name`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	teams = make(map[string]*Team)

	// Loop through rows
	for rows.Next() {
		var teamIdBuffer sql.NullString
		var teamName string

		err = rows.Scan(&teamIdBuffer, &teamName)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Append team to teams array
		teamId := uuid.Parse(teamIdBuffer.String)
		team := &Team{teamId, game.GameId, teamName}
		teams[teamId.String()] = team
	}
	return teams, nil
}

// Get a list of team ids with players currently in the game
func (game *Game) GetActiveTeamIds() (teamsList []uuid.UUID, appErr *ApplicationError) {
	rows, err := db.Query(`SELECT distinct(team_id) FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	for rows.Next() {
		var teamIdBuffer string
		err = rows.Scan(&teamIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		teamId := uuid.Parse(teamIdBuffer)
		teamsList = append(teamsList, teamId)
	}
	return teamsList, nil
}

// Gets a list of teams for a game
func (game *Game) GetTeams() (teams []*Team, appErr *ApplicationError) {
	// Query Db
	rows, err := db.Query(`SELECT team_id, team_name FROM dm_teams WHERE game_id = $1 ORDER BY team_name`, game.GameId.String())
	if err == sql.ErrNoRows {
		fmt.Println("No teams")
		return nil, NewApplicationError("No Team", err, ErrCodeNoTeams)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Loop through rows
	for rows.Next() {
		var teamIdBuffer sql.NullString
		var teamName string

		err = rows.Scan(&teamIdBuffer, &teamName)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Append team to teams array
		teamId := uuid.Parse(teamIdBuffer.String)
		team := &Team{teamId, game.GameId, teamName}
		teams = append(teams, team)
	}
	return teams, nil
}

// Creates a new team and returns it
func (game *Game) NewTeam(teamName string) (team *Team, appErr *ApplicationError) {
	// Generate a uuid and insert the team
	teamId := uuid.NewRandom()
	res, err := db.Exec(`INSERT INTO dm_teams (team_id, game_id, team_name) VALUES ($1,$2,$3)`, teamId.String(), game.GameId.String(), teamName)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Make sure the insert worked
	appErr = WereRowsAffected(res)
	if appErr != nil {
		return nil, appErr
	}

	// Return the team
	return &Team{teamId, game.GameId, teamName}, nil
}

// Removes all players from a team and deletes it
func DeleteTeam(teamId uuid.UUID) (appErr *ApplicationError) {

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

	appErr = WereRowsAffected(res)
	if appErr != nil {
		return appErr
	}

	tx.Commit()
	return nil
}

// Gets a uuid slice of team captains who are alive
func (game *Game) GetAllTeamCaptains() (captains []uuid.UUID, appErr *ApplicationError) {
	// Get list of Dead Captain Ids
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE user_role = 'dm_captain' AND game_id = $1`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return ConvertUserIdRowsToSlice(rows)
}

// Gets a uuid slice of team captains who are alive
func (game *Game) GetAliveTeamCaptains() (captains []uuid.UUID, appErr *ApplicationError) {
	// Get list of Dead Captain Ids
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE user_role = 'dm_captain' AND alive = true AND game_id = $1`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return ConvertUserIdRowsToSlice(rows)
}

// Gets a uuid slice of team captains who are dead
func (game *Game) GetDeadTeamCaptains() (captains []uuid.UUID, appErr *ApplicationError) {
	// Get list of Dead Captain Ids
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE user_role = 'dm_captain' AND alive = false AND game_id = $1`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return ConvertUserIdRowsToSlice(rows)
}

// Gets the team for a user and gameId
func (user *User) GetTeam(gameId uuid.UUID) (team *Team, appErr *ApplicationError) {
	gameMapping, appErr := GetGameMapping(user.UserId, gameId)
	if appErr != nil {
		return nil, appErr
	}
	return GetTeamById(gameMapping.TeamId)
}

// Adds a user to a team
func (user *User) JoinTeam(teamId uuid.UUID) (gameMapping *GameMapping, appErr *ApplicationError) {
	// Get the game_id (it's easier to enforce this constraint here than the DB)
	var gameIdBuffer string
	err := db.QueryRow(`SELECT game_id FROM dm_teams WHERE team_id = $1`, teamId.String()).Scan(&gameIdBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.Parse(gameIdBuffer)
	if gameId == nil {
		msg := "Invalid Team Id: " + teamId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeDatabase)
	}

	var userRole, secret string
	var kills int
	var alive bool

	// Updates the user's game_mapping to include their new team id
	err = db.QueryRow(`UPDATE dm_user_game_mapping SET team_id = $1 WHERE game_id = $2 and user_id = $3 RETURNING user_role, secret, kills, alive`, teamId.String(), gameId.String(), user.UserId.String()).Scan(&userRole, &secret, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return &GameMapping{user.UserId, gameId, teamId, userRole, secret, kills, alive}, nil

}

// removes a user from a team
func (user *User) LeaveTeam(teamId uuid.UUID) (gameMapping *GameMapping, appErr *ApplicationError) {
	// Glean the game_id from the team to enforce the game_id constraint
	var gameIdBuffer string
	err := db.QueryRow(`SELECT game_id FROM dm_teams WHERE team_id = $1`, teamId.String()).Scan(&gameIdBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.Parse(gameIdBuffer)
	if gameId == nil {
		msg := "Invalid Team Id: " + teamId.String()
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeDatabase)
	}

	var userRole, secret string
	var kills int
	var alive bool

	// Sets the user's team_id to null
	err = db.QueryRow(`UPDATE dm_user_game_mapping SET team_id = null WHERE game_id = $1 and user_id = $2 RETURNING user_role, secret, kills, alive`, gameId.String(), user.UserId.String()).Scan(&userRole, &secret, &kills, &alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return &GameMapping{user.UserId, gameId, nil, userRole, secret, kills, alive}, nil
}

// Rename a team
func (team *Team) Rename(newName string) (appErr *ApplicationError) {
	// Run teh update
	res, err := db.Exec(`UPDATE dm_teams SET team_name = $1 WHERE team_id = $2`, newName, team.TeamId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr = WereRowsAffected(res)
	if appErr != nil {
		return appErr
	}

	team.TeamName = newName
	return nil
}

// Gets the user id for the team captain
func (team *Team) GetTeamCaptainId() (captainId uuid.UUID, appErr *ApplicationError) {
	var captainIdBuffer string

	// Get captain id from db
	err := db.QueryRow(`SELECT user_id FROM dm_user_game_mapping WHERE team_id = $1 and user_role = 'dm_captain'`, team.TeamId.String()).Scan(&captainIdBuffer)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return uuid.Parse(captainIdBuffer), nil
}

// is it safe to assign targets by teams
func (game *Game) CanAssignByTeams() (canAssign bool, appErr *ApplicationError) {
	var numUsers, numCaptains int
	var teamIdBuffer string

	rows, err := db.Query(`SELECT count(user_id), team_id from dm_user_game_mapping WHERE alive = true AND game_id = $1 GROUP BY team_id`, game.GameId.String())
	if err != nil {
		return false, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	for rows.Next() {
		err = rows.Scan(&numUsers, &teamIdBuffer)
		if err != nil {
			return false, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		if numUsers < MinUsersPerTeam {
			return false, nil
		}
	}

	rows, err = db.Query(`SELECT count(user_id), team_id from dm_user_game_mapping WHERE alive = true AND user_role = 'dm_captain' AND game_id = $1 GROUP BY team_id`, game.GameId.String())
	if err != nil {
		return false, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	for rows.Next() {
		err = rows.Scan(&numCaptains, &teamIdBuffer)
		if err != nil {
			return false, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		if numCaptains != 1 {
			return false, nil
		}
	}

	return true, nil

}
