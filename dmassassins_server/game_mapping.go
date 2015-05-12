package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

	// Confirm the user has the right password
	appErr = CheckPassword(gameId, gamePassword)
	if appErr != nil {
		return nil, appErr
	}

	// Assign the user a secret
	secret, appErr := NewSecret(3)
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

// Can a user become a captain
func (gameMapping *GameMapping) CanBecomeCaptain() (appErr *ApplicationError) {
	// Get the team
	team, appErr := GetTeamById(gameMapping.TeamId)
	if appErr != nil {
		return appErr
	}

	// Get the team captain
	captainId, appErr := team.GetTeamCaptainId()
	if appErr != nil {
		return appErr
	}

	// If we don't have a team captain demote them
	if captainId == nil {
		return nil
	}

	demoteString := ``
	user, appErr := GetUserById(captainId)
	if appErr == nil {
		demoteString = "\nDemote " + user.Properties[`first_name`] + ` ` + user.Properties[`last_name`] + ` to promote another user.`
	}

	// Return an error if we have a captain
	msg := `A team cannot have two captains!` + demoteString
	err := errors.New(msg)
	return NewApplicationError(msg, err, ErrCodeCaptainExists)
}

// Changes a user's role within a game
func (gameMapping *GameMapping) ChangeRole(role string) (appErr *ApplicationError) {

	if role == `dm_captain` {
		appErr = gameMapping.CanBecomeCaptain()
		if appErr != nil {
			return appErr
		}
	}

	// Update the user role
	_, err := db.Exec(`UPDATE dm_user_game_mapping SET user_role = $1 WHERE user_id = $2 AND game_id = $3`, role, gameMapping.UserId.String(), gameMapping.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Change the role on the user struct
	gameMapping.UserRole = role
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
	err := db.QueryRow(`SELECT user_id FROM dm_user_game_mapping WHERE user_role = 'dm_admin' AND game_id = $1 LIMIT 1`, game.GameId.String()).Scan(&userIdBuffer)
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

// deletes the actual game mapping from the db
func (gameMapping *GameMapping) delete() (appErr *ApplicationError) {
	_, err := db.Exec(`DELETE from dm_user_game_mapping WHERE user_id = $1 and game_id = $2`, gameMapping.UserId.String(), gameMapping.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// the game mapping no longer exists so set it to nil
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

	// Kill the assassins target (silently)
	if appErr == nil {
		_, _, appErr = assassin.KillTarget(gameMapping.GameId, secret, false)
		if appErr != nil {
			return appErr
		}
		// If we don't have an assassin confirm the secret
	} else if !strings.EqualFold(secret, gameMapping.Secret) {
		msg := fmt.Sprintf("Invalid secret: %s", secret)
		err := errors.New("Invalid Secret")
		return NewApplicationError(msg, err, ErrCodeInvalidSecret)
	}

	// Delete the game mapping from the db
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
	if appErr != nil && appErr.Code != ErrCodeNoTeams {
		return nil, appErr
	}

	rows, err := db.Query(`SELECT users.user_id, users.email, users.facebook_id, map.team_id, map.user_role, map.alive FROM dm_users AS users, dm_user_game_mapping AS map WHERE users.user_id = map.user_id AND map.game_id = $1`, game.GameId.String())
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
		var alive bool

		// Scan in variables
		err := rows.Scan(&userIdBuffer, &email, &facebookId, &teamIdBuffer, &userRole, &alive)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Organize user info
		userId := uuid.Parse(userIdBuffer)
		teamId := uuid.Parse(teamIdBuffer.String)
		properties := make(map[string]string)

		properties["user_role"] = userRole
		properties["alive"] = strconv.FormatBool(alive)

		if teams != nil {
			properties["team_id"] = teamId.String()
			if team, ok := teams[teamId.String()]; ok {
				properties["team"] = team.TeamName
			} else {
				properties["team"] = "null"
			}

		}

		// Create the user struct and add it
		user := &User{userId, "", email, facebookId, properties}
		users[userId.String()] = user
		userIds = append(userIds, userId.String())
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
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
	err = stmt.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Add properties to users
	for rows.Next() {
		var userIdBuffer, key, value string
		rows.Scan(&userIdBuffer, &key, &value)
		users[userIdBuffer].Properties[key] = value
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return users, nil
}

// Revive a previously killed user
func (gameMapping *GameMapping) Revive() (assassinId, targetId uuid.UUID, appErr *ApplicationError) {

	// Get a random assassin/target pair
	var assassinIdBuffer, targetIdBuffer string
	err := db.QueryRow(`SELECT user_id, target_id FROM dm_user_targets where game_id = $1 ORDER BY RANDOM() LIMIT 1`, gameMapping.GameId.String()).Scan(&assassinIdBuffer, &targetIdBuffer)
	if err == sql.ErrNoRows {
		return nil, nil, NewApplicationError("No Users left", err, ErrCodeNoUsers)
	}
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare the statement to kill the old target
	setAssassin, err := db.Prepare(`UPDATE dm_user_targets SET target_id = $1 WHERE user_id = $2 AND game_id = $3`)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setAssassin).Exec(gameMapping.UserId.String(), assassinIdBuffer, gameMapping.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare the statement to kill the old target
	setTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id, game_id) VALUES ($1, $2, $3)`)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setTarget).Exec(gameMapping.UserId.String(), targetIdBuffer, gameMapping.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare the statement to kill the old target
	setAlive, err := db.Prepare(`UPDATE dm_user_game_mapping SET alive = true WHERE user_id = $1 AND game_id = $2`)
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to kill the old target
	_, err = tx.Stmt(setAlive).Exec(gameMapping.UserId.String(), gameMapping.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// check transaction for errors
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameMapping.Alive = true
	assassinId = uuid.Parse(assassinIdBuffer)
	targetId = uuid.Parse(targetIdBuffer)
	return assassinId, targetId, nil

}

// Get a slice of the strongest player for each team, ties are broken arbitrarily
func (game *Game) GetStrongPlayers() (strong []uuid.UUID, appErr *ApplicationError) {
	// segregate strong users
	rows, err := db.Query(`SELECT DISTINCT ON (team_id) user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY team_id, kills desc`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// parse strong users into slice
	for rows.Next() {
		var strongUserIdBuffer string
		err = rows.Scan(&strongUserIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Add strong userId to strong slice
		strongUserId := uuid.Parse(strongUserIdBuffer)
		strong = append(strong, strongUserId)
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return strong, nil
}

// Get a slice of the strongest player for each team, ties are broken arbitrarily
func (game *Game) getStrongPlayersWithState(alive bool) (strong []uuid.UUID, appErr *ApplicationError) {
	// segregate strong users
	rows, err := db.Query(`SELECT DISTINCT ON (team_id) user_id, alive FROM dm_user_game_mapping WHERE game_id = $1 AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY team_id, kills desc`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// parse strong users into slice
	for rows.Next() {
		var strongUserIdBuffer string
		var playerState bool
		err = rows.Scan(&strongUserIdBuffer, &playerState)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		if playerState != alive {
			continue
		}
		// Add strong userId to strong slice
		strongUserId := uuid.Parse(strongUserIdBuffer)
		strong = append(strong, strongUserId)
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return strong, nil
}

// Get a slice of the weakest player for each team, ties are broken arbitrarily
func (game *Game) GetWeakPlayers() (weak []uuid.UUID, appErr *ApplicationError) {
	// segregate weak users
	rows, err := db.Query(`SELECT DISTINCT ON (team_id) user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY team_id, kills asc`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// parse weak users into slice
	for rows.Next() {
		var weakUserIdBuffer string
		err = rows.Scan(&weakUserIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Add weak userId to weak slice
		weakUserId := uuid.Parse(weakUserIdBuffer)
		weak = append(weak, weakUserId)
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return weak, nil
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

// Sets a user as dead without going through the motions of using their assassin
func (gameMapping *GameMapping) MarkDead() (appErr *ApplicationError) {
	_, err := db.Exec(`UPDATE dm_user_game_mapping SET alive = false WHERE user_id = $1 AND game_id = $2`, gameMapping.UserId.String(), gameMapping.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	gameMapping.Alive = false
	return nil
}
