package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"strconv"
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
	err := db.QueryRow(`SELECT user_id from dm_user_targets WHERE target_id = $1 AND game_id = $2 LIMIT 1`, targetId.String(), gameId.String()).Scan(&assassinIdBuffer)
	if err == sql.ErrNoRows {
		msg := "Invalid target_id: " + targetId.String()
		return nil, NewApplicationError(msg, err, ErrCodeNotFoundUserId)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Grab assassin Id
	assassinId := uuid.Parse(assassinIdBuffer)

	// Get the user through the normal routes
	return GetUserById(assassinId)

}

// Gets a game's admin
func (game *Game) GetAdmin() (admin *User, appErr *ApplicationError) {
	var userIdBuffer string

	// Query the database
	err := db.QueryRow(`SELECT user_id FROM dm_user_game_mapping WHERE user_role = 'dm_admin' AND game_id = $1`, game.GameId.String()).Scan(&userIdBuffer)
	if err == sql.ErrNoRows {
		return nil, NewApplicationError("No Admin", err, ErrCodeNotFoundGameMapping)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	userId := uuid.Parse(userIdBuffer)
	if userId == nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return GetUserById(userId)
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
	assassin, appErr := GetAssassin(gameMapping.UserId, gameMapping.GameId)
	if appErr != nil && appErr.Code != ErrCodeNotFoundUserId {
		return appErr
	}

	if appErr == nil {
		_, appErr = assassin.KillTarget(gameMapping.GameId, secret, false)
		if appErr != nil {
			return appErr
		}
	}

	appErr = gameMapping.delete()
	if appErr != nil {
		return appErr
	}

	return nil
}

// Get all games for a user
func (user *User) GetGamesForUser() (games []*Game, appErr *ApplicationError) {

	// Select game_ids from the dm_user_game_mapping table and use those to get the games
	rows, err := db.Query(`SELECT game.game_id, game.game_name, game.game_started, game_password FROM dm_games AS game WHERE game_id IN (SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1) ORDER BY game_name`, user.UserId.String())
	if err == sql.ErrNoRows {
		return nil, NewApplicationError("No Game Mappings", err, ErrCodeNoGameMappings)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// convert the rows to an array of gamess
	return parseGameRows(rows)
}

// Get all games a user is not present in so they can join one
func (user *User) GetNewGamesForUser() (games []*Game, appErr *ApplicationError) {

	// Select game_ids from the dm_user_game_mapping table and skip those in the dm_games datable
	rows, err := db.Query(`SELECT game.game_id, game.game_name, game.game_started, game_password FROM dm_games AS game WHERE game_started = false AND game_id NOT IN (SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1) ORDER BY game_name`, user.UserId.String())
	if err == sql.ErrNoRows {
		return nil, NewApplicationError("No Game Mappings", err, ErrCodeNoGameMappings)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return parseGameRows(rows)
}

// Get all users for a game
func (game *Game) GetAllUsersForGame() (users map[string]*User, appErr *ApplicationError) {
	teams, appErr := game.GetTeamsMap()
	if appErr != nil && appErr.Code != ErrCodeNoTeams{
		return nil, appErr
	}

	rows, err := db.Query(`SELECT users.user_id, users.email, users.facebook_id, map.team_id, map.user_role FROM dm_users AS users, dm_user_game_mapping AS map WHERE users.user_id = map.user_id AND map.game_id = $1`, game.GameId.String())
	if err == sql.ErrNoRows {
		return nil, NewApplicationError("No Users For Game", err, ErrCodeNoGameMappings)
	}
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var userIds []interface{}
	users = make(map[string]*User)
	// Loop through users
	for rows.Next() {
		var userIdBuffer, email, facebookId, userRole string
		var teamIdBuffer sql.NullString

		// Scan in variables
		err := rows.Scan(&userIdBuffer, &email, &facebookId, &teamIdBuffer, &userRole)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Organize user info
		userId := uuid.Parse(userIdBuffer)
		teamId := uuid.Parse(teamIdBuffer.String)
		properties := make(map[string]string)
		
		if teams != nil {
			if team, ok := teams[teamId.String()]; ok {
				properties["team"] = team.TeamName
			} else {
				properties["team"] = "No Team"
			}
			
		}

		// Create the user struct and add it
		user := &User{userId, "", email, facebookId, properties}
		users[userId.String()] = user
		userIds = append(userIds, userId.String())
	}

	// Build Sql query for properties
	max := len(userIds)
	params := `$1`
	for i := 0; i < max; i++ {
		params += `, $` + strconv.Itoa(i+1)
	}
	sql := `SELECT user_id, key, value FROM dm_user_properties WHERE user_id IN ( ` + params + ` )`

	// Prepare Sql Query for properties
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute Sql Query for properties
	rows, err = stmt.Query(userIds...)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Add properties to users
	for rows.Next() {
		var userIdBuffer, key, value string
		rows.Scan(&userIdBuffer, &key, &value)
		users[userIdBuffer].Properties[key] = value
	}

	return users, nil
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
