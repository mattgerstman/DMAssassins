package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	GameId      uuid.UUID         `json:"game_id"`
	GameName    string            `json:"game_name"`
	Started     bool              `json:"game_started"`
	HasPassword bool              `json:"game_has_password"`
	Properties  map[string]string `json:"game_properties"`
}

// Rename a game
func (game *Game) Rename(newName string) (appErr *ApplicationError) {
	res, err := db.Exec(`UPDATE dm_games SET game_name = $1 WHERE game_id = $2`, newName, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Make sure at least one game was affected
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	// Set the name in the struct if we use it later
	game.GameName = newName
	return nil

}

// Change a game's password
func (game *Game) ChangePassword(newPassword string) (appErr *ApplicationError) {
	res, err := db.Exec(`UPDATE dm_games SET game_password = $1 WHERE game_id = $2`, newPassword, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Make sure at least one game was affected
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	// Set the password in the struct if we use it later
	game.HasPassword = newPassword == ""
	game.Properties["game_password"] = newPassword
	return nil

}

// Get all games
func GetGameList() (games []*Game, appErr *ApplicationError) {
	// Query db for all games
	rows, err := db.Query(`SELECT game_id, game_name, game_started, game_password FROM dm_games ORDER BY game_name`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Convert the rows to an array of games
	return parseGameRows(rows)
}

// Get a game struct by it's ID
func GetGameById(gameId uuid.UUID) (game *Game, appErr *ApplicationError) {
	var gameName string
	var gameStarted bool
	var gamePasswordBuffer sql.NullString

	// Query DB for game
	err := db.QueryRow(`SELECT game_name, game_started, game_password FROM dm_games WHERE game_id = $1`, gameId.String()).Scan(&gameName, &gameStarted, &gamePasswordBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Check if it has a password
	hasPassword := false
	if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
		hasPassword = true
	}

	properties := make(map[string]string)

	// Return the game
	game = &Game{gameId, gameName, gameStarted, hasPassword, properties}
	_, appErr = game.GetGameProperties()
	if appErr != nil {
		return nil, appErr
	}
	return game, nil
}

// End a game
func (game *Game) End() (appErr *ApplicationError) {

	// Set game_started to false
	res, err := db.Exec("UPDATE dm_games SET game_started = false WHERE game_id = $1", game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Make sure at least one game was affected
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	game.Started = false
	return nil
}

// Determine if any players need a team before starting
func (game *Game) doAnyPlayersNeedTeams() (appErr *ApplicationError) {
	var count int
	err := db.QueryRow("SELECT count(user_id) FROM dm_user_game_mapping WHERE game_id = $1 AND team_id IS NULL", game.GameId.String()).Scan(&count)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	if count != 0 {
		return NewApplicationError("Every player must be assigned a team to start", err, ErrCodePlayerMissingTeam)
	}
	return nil
}

// Get the number of players for a game
func (game *Game) GetNumActivePlayers() (count int, appErr *ApplicationError) {

	err := db.QueryRow("SELECT count(user_id) FROM dm_user_game_mapping WHERE (user_role = 'dm_user' OR user_role = 'dm_captain') AND alive = true AND game_id = $1", game.GameId.String()).Scan(&count)
	if err != nil {
		return 0, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return count, nil
}

// Start a game
func (game *Game) Start() (appErr *ApplicationError) {

	// Make sure we have enough players to start the game
	count, appErr := game.GetNumActivePlayers()
	if count < 4 {
		err := errors.New("Not Enough Players")
		return NewApplicationError("You must have at least 4 players to start a game", err, ErrCodeNeedMorePlayers)
	}

	// If teams are enabled make sure all users have a team to start
	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return appErr
	}
	if teamsEnabled == `true` {
		anyPlayersNeedTeams := game.doAnyPlayersNeedTeams()
		if anyPlayersNeedTeams != nil {
			return anyPlayersNeedTeams
		}
	}

	// First assign targets for the game
	appErr = game.AssignTargets()
	if appErr != nil {
		return appErr
	}
	// Set started = true
	res, err := db.Exec("UPDATE dm_games SET game_started = true WHERE game_id = $1", game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Make sure we affected at least one row
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}

	// Update the given struct
	game.Started = true
	return nil
}

// Converts a set of game rows into an array of games
// Rows MUST appear in the order game_id, game_name, game_started, game_password
func parseGameRows(rows *sql.Rows) (games []*Game, appErr *ApplicationError) {

	// Loop through games
	for rows.Next() {
		var gameId uuid.UUID
		var gamePasswordBuffer sql.NullString
		var gameIdBuffer, gameName string
		var gameStarted bool

		// Scan in variables
		err := rows.Scan(&gameIdBuffer, &gameName, &gameStarted, &gamePasswordBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Check if a password exists
		hasPassword := false
		if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
			hasPassword = true
		}

		// Create the game struct and apparend it to the list
		gameId = uuid.Parse(gameIdBuffer)
		properties := make(map[string]string)
		game := &Game{gameId, gameName, gameStarted, hasPassword, properties}
		games = append(games, game)
	}
	return games, nil
}

// Gets a game's password or returns an empty string if there is none
func (game *Game) GetPassword() (gamePassword string, appErr *ApplicationError) {
	var storedPassword sql.NullString
	err := db.QueryRow(`SELECT game_password FROM dm_games WHERE game_id = $1`, game.GameId.String()).Scan(&storedPassword)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return storedPassword.String, nil
}

// Checks if a given plaintext password is right for a game
// Returns an error if the password doesn't match
func CheckPassword(gameId uuid.UUID, testPassword string) (appErr *ApplicationError) {
	var storedPassword sql.NullString
	err := db.QueryRow(`SELECT game_password FROM dm_games WHERE game_id = $1`, gameId.String()).Scan(&storedPassword)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	if strings.EqualFold(testPassword, storedPassword.String) {
		return nil
	}
	//CheckPasswordHash(hashedPassword, testPassword)
	msg := "Invalid Game Password: " + testPassword
	err = errors.New(msg)
	return NewApplicationError(msg, err, ErrCodeInvalidGamePassword)
}

// Compares a hashed password to a plaintext password
// Returns an error if they don't match
func CheckPasswordHash(hashedPassword []byte, plainPw string) (appErr *ApplicationError) {
	// If they're both nil return true
	if len(hashedPassword) == 0 && plainPw == "" {
		return nil
	}

	// Convert to a bytearray and skip whitespace
	bytePW := []byte(strings.TrimSpace(plainPw))
	err := bcrypt.CompareHashAndPassword(hashedPassword, bytePW)
	if err != nil {
		msg := "Invalid Game Password: " + plainPw
		err = errors.New(msg)
		return NewApplicationError(msg, err, ErrCodeInvalidGamePassword)
	}

	return nil
}

// Encrpyt a password for storage
func Crypt(plainPw string) (hashedPassword []byte, appErr *ApplicationError) {

	// If it's blank don't encrypt null
	if plainPw == "" {
		return nil, nil
	}

	// Convert to a bytearray and skip whitespace
	bytePw := []byte(strings.TrimSpace(plainPw))
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePw, bcrypt.DefaultCost)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeInternalServerWTF)
	}
	return hashedPassword, nil
}

// Creates a new game and saves it in the database
func NewGame(gameName string, userId uuid.UUID, gamePassword string) (game *Game, appErr *ApplicationError) {
	// Start a transaction, god knows we can't break anything
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare the statement to insert the game into the db
	newGame, err := db.Prepare(`INSERT INTO dm_games (game_id, game_name, game_password) VALUES ($1, $2, $3)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Generate a UUID and execute the statement to insert the game into the db
	gameId := uuid.NewRandom()
	res, err := tx.Stmt(newGame).Exec(gameId.String(), gameName, gamePassword)
	// Check to make sure the insert happened
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}

	// Create a user secret for the game
	secret, appErr := NewSecret(3)
	if appErr != nil {
		tx.Rollback()
		return nil, appErr
	}
	// Prepare the statement to insert the game creator(admin) into the game
	firstMapping, err := db.Prepare(`INSERT INTO dm_user_game_mapping (game_id, user_id, user_role, secret) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Executre the statement to insert the game creator(admin) into the game
	role := "dm_admin"
	res, err = tx.Stmt(firstMapping).Exec(gameId.String(), userId.String(), role, secret)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Check to make sure the insert happened
	NoRowsAffectedAppErr = WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}
	tx.Commit()
	hasPassword := gamePassword != ""

	properties := make(map[string]string)

	game = &Game{gameId, gameName, false, hasPassword, properties}
	_, appErr = game.GetGameProperties()
	if appErr != nil {
		return nil, appErr
	}

	return game, nil

}

// Assigns targets using a methodology
func (game *Game) AssignTargetsBy(assignmentType string) (appErr *ApplicationError) {

	// Reverse targets
	if assignmentType == `reverse` {
		return game.ReverseTargets()
	}

	if assignmentType == `strong_weak` {
		return game.StrongTargetWeak()
	}

	// If they don't have a special method check if teams are enabled
	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return appErr
	}

	// If teams are enabled assigned by team
	if teamsEnabled == `true` {
		canAssign, appErr := game.CanAssignByTeams()
		if appErr != nil {
			return appErr
		}

		if canAssign {
			return game.AssignTargetsByTeams()
		}
	}
	// Fallback to plain random assignment
	return game.AssignTargets()
}

type targetPair struct {
	AssassinId       uuid.UUID `json:assassin_id`
	AssassinTeamId   uuid.UUID `json:assassin_team_id`
	AssassinUserRole string    `json:assassin_user_role`
	TargetId         uuid.UUID `json:target_id`
	TargetTeamId     uuid.UUID `json:target_team_id`
	TargetUserRole   string    `json:target_user_role`
}

// Reverses targets for a game
func (game *Game) ReverseTargets() (appErr *ApplicationError) {
	var userIdBuffer, targetIdBuffer string
	var targets []*targetPair
	rows, err := db.Query(`SELECT user_id, target_id FROM dm_user_targets WHERE game_id = $1`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	for rows.Next() {
		rows.Scan(&userIdBuffer, &targetIdBuffer)
		userId := uuid.Parse(userIdBuffer)
		targetId := uuid.Parse(targetIdBuffer)
		pair := &targetPair{targetId, nil, "", userId, nil, ""}
		targets = append(targets, pair)
	}

	return game.insertTargets(targets)
}

// have the strongest players target the weakest ones, don't be concerned about teams/captains
func (game *Game) StrongTargetWeak() (appErr *ApplicationError) {
	var strong, weak []uuid.UUID

	// segregate strong users
	rows, err := db.Query(`SELECT DISTINCT ON (team_id) user_id FROM dm_user_game_mapping WHERE game_id = $1 ORDER BY team_id, kills desc`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// parse strong users into slice
	for rows.Next() {
		var strongUserIdBuffer string
		err = rows.Scan(&strongUserIdBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Add strong userId to strong slice
		strongUserId := uuid.Parse(strongUserIdBuffer)
		strong = append(strong, strongUserId)
	}

	// segregate weak users
	rows, err = db.Query(`SELECT DISTINCT ON (team_id) user_id FROM dm_user_game_mapping WHERE game_id = $1 ORDER BY team_id, kills asc`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// We need to offset strong and weak to ensure that strong players aren't targetting their own weakest palyer
	var firstWeakUserIdBuffer string
	rows.Next()
	rows.Scan(&firstWeakUserIdBuffer)
	firstWeakUserId := uuid.Parse(firstWeakUserIdBuffer)

	// parse weak users into slice
	for rows.Next() {
		var weakUserIdBuffer string
		err = rows.Scan(&weakUserIdBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Add weak userId to weak slice
		weakUserId := uuid.Parse(weakUserIdBuffer)
		weak = append(weak, weakUserId)
	}
	weak = append(weak, firstWeakUserId)

	// Construct params to get non strong/weak users
	params := `$2`
	numParams := (len(weak) * 2) + 1
	for i := 3; i < numParams; i++ {
		params += `, $` + strconv.Itoa(i)
	}

	// Select other users from the db
	var strongWeak []interface{}
	strongWeak = append(strongWeak, game.GameId.String())
	for _, strongId := range strong {
		strongWeak = append(strongWeak, strongId.String())
	}
	for _, weakId := range weak {
		strongWeak = append(strongWeak, weakId.String())
	}

	rows, err = db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') AND user_id NOT IN (`+params+`)`, strongWeak...)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get number of players
	numPlayers, appErr := game.GetNumActivePlayers()
	if appErr != nil {
		return appErr
	}

	// Get number of teams
	numTeams := len(weak)

	// Determine how often to insert a strong/weak pair
	insertNum := (numPlayers / numTeams) - 1

	// Create first strong/weak pair and make them the first target
	var targets []*targetPair
	firstStrong, strong := strong[0], strong[1:]
	firstWeak, weak := weak[0], weak[1:]
	firstPair := &targetPair{firstStrong, nil, "", firstWeak, nil, ""}

	targets = append(targets, firstPair)

	// we already have one pair
	i := 1
	lastUserId := firstWeak
	for rows.Next() {
		// Get userId
		var userIdBuffer string
		err = rows.Scan(&userIdBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Create target pair
		userId := uuid.Parse(userIdBuffer)
		pair := &targetPair{lastUserId, nil, "", userId, nil, ""}

		// Add pair to targets slice
		targets = append(targets, pair)

		// Set last userId and increment
		lastUserId = userId
		i++

		// Unless we're at the right point to insert just continue
		if (i % insertNum) != 0 {
			continue
		}

		// If we're out of strong/weak make it so we won't get here again
		if len(strong) == 0 {
			insertNum = numPlayers + 1
			continue
		}

		var nextStrong, nextWeak uuid.UUID
		// Get next strong/weak pair
		nextStrong, strong = strong[0], strong[1:]
		nextWeak, weak = weak[0], weak[1:]

		// Set up next strong/weak pair and insert it
		pair = &targetPair{lastUserId, nil, "", nextStrong, nil, ""}
		strongWeakPair := &targetPair{nextStrong, nil, "", nextWeak, nil, ""}
		targets = append(targets, pair, strongWeakPair)

		lastUserId = nextWeak
		i += 2
	}

	// have last user target first strong user
	lastTarget := &targetPair{lastUserId, nil, "", firstWeak, nil, ""}
	targets = append(targets, lastTarget)

	return nil
}

// Assign targets and space them out by team
func (game *Game) AssignTargetsByTeams() (appErr *ApplicationError) {
	// Get new target list
	rows, err := db.Query(`SELECT user_id, team_id, user_role FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY random()`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get the list of team ids
	teamsList, appErr := game.GetActiveTeamIds()
	if appErr != nil {
		return appErr
	}

	// Create userList, captainList, and buffer variables
	var userIdBuffer, teamIdBuffer, userRole string
	userList := make(map[string][]uuid.UUID)
	captainList := make(map[string]uuid.UUID)

	// Fill in userList and captainList with valid slices
	for _, team := range teamsList {
		userList[team.String()] = []uuid.UUID{}
	}

	// parse out users into teams and captain
	numUsers := 0
	for rows.Next() {
		// Get the user_id from the row
		err = rows.Scan(&userIdBuffer, &teamIdBuffer, &userRole)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// If user is a captain put them in the captains list
		userId := uuid.Parse(userIdBuffer)
		if CompareRole(userRole, RoleCaptain) {
			captainList[teamIdBuffer] = userId
			continue
		}
		// If they aren't a captain add them to the userList
		userList[teamIdBuffer] = append(userList[teamIdBuffer], userId)
		numUsers++
	}

	numTeams := len(teamsList)
	originalNumTeams := numTeams
	rand.Seed(time.Now().UTC().UnixNano())

	assigned := 1
	var targetList []*targetPair
	randTeam := teamsList[rand.Intn(numTeams)]

	var firstUser, firstUserTeam, currentUser, currentUserTeam, lastUser, lastUserTeam uuid.UUID

	firstUser, userList[randTeam.String()] = userList[randTeam.String()][0], userList[randTeam.String()][1:]

	firstUserTeam = randTeam
	lastUserTeam = firstUserTeam
	lastUser = firstUser

	// While we still have users to assign
	for assigned < numUsers {
		// Get a random team
		randTeamIndex := rand.Intn(numTeams)
		currentUserTeam = teamsList[randTeamIndex]

		// If the random team is the same as the last one just go to the next team
		if uuid.Equal(currentUserTeam, lastUserTeam) {
			randTeamIndex++
			if randTeamIndex >= numTeams {
				randTeamIndex = 0
			}
			currentUserTeam = teamsList[randTeamIndex]
		}

		// If our current team has no members delete it and go to the next one
		for len(userList[currentUserTeam.String()]) == 0 {
			// delete the current team from the userlist
			delete(userList, currentUserTeam.String())

			// delete the currentTeam from the teamsList
			teamsList = append(teamsList[:randTeamIndex], teamsList[(randTeamIndex+1):]...)

			// change numTeams to reflect the current number of teams
			numTeams = len(teamsList)
			// If our index is greater than the number of teams set it to 0
			if randTeamIndex >= numTeams {
				randTeamIndex = 0
			}

			// set the current user team
			currentUserTeam = teamsList[randTeamIndex]
		}

		// pop out currentUser from userList[currentUserTeam]
		currentUser, userList[currentUserTeam.String()] = userList[currentUserTeam.String()][0], userList[currentUserTeam.String()][1:]

		// append lastUser currentUser to the targetsList
		userTargetPair := &targetPair{lastUser, lastUserTeam, `dm_user`, currentUser, currentUserTeam, `dm_user`}
		targetList = append(targetList, userTargetPair)

		lastUser = currentUser
		lastUserTeam = currentUserTeam

		assigned++
	}

	userTargetPair := &targetPair{lastUser, lastUserTeam, `dm_user`, firstUser, firstUserTeam, `dm_user`}
	targetList = append(targetList, userTargetPair)

	// Space out captains by 5 unless we don't have enough users
	captainSpace := 5
	numTeams = originalNumTeams
	userTeamRatio := (numUsers / numTeams)
	if userTeamRatio < captainSpace {
		captainSpace = userTeamRatio
	}

	i := 0
	triesForRole := 5
	triesForTeam := 10
	// Loop through captains
	for captainTeam, captainId := range captainList {
		j := 0
		captainTeamId := uuid.Parse(captainTeam)
		// we have to use a marker to determine if we've found an appropriate pair to insert the captain in
		foundPair := false
		for !foundPair {
			if i >= numUsers {
				i = 0
				j++
			}

			// If we've lapped users enough times allow users form the same team to target each other
			if j > triesForTeam {
				foundPair = true
				continue
			}

			// Check if the assassin and target both have different teams than the captain
			assassinTeamId, targetTeamId := targetList[i].AssassinTeamId, targetList[i].TargetTeamId
			assassinUserRole, targetUserRole := targetList[i].AssassinUserRole, targetList[i].TargetUserRole
			if uuid.Equal(captainTeamId, assassinTeamId) {
				i++
				continue
			}
			if uuid.Equal(captainTeamId, targetTeamId) {
				i++
				continue
			}

			// If we've lapped users enough times allow captains to target each other
			if j > triesForRole {
				continue
			}

			if assassinUserRole == `dm_captain` {
				i++
				continue
			}

			if targetUserRole == `dm_captain` {
				i++
				continue
			}

			// If neither has the same team, we've found a match
			foundPair = true
		}
		// Get all the assassin/target information
		assassinId, assassinTeamId, assassinUserRole, targetId, targetTeamId, targetUserRole := targetList[i].AssassinId, targetList[i].AssassinTeamId, targetList[i].AssassinUserRole, targetList[i].TargetId, targetList[i].TargetTeamId, targetList[i].TargetUserRole
		// Insert the captain after the assassin
		targetList[i] = &targetPair{assassinId, assassinTeamId, assassinUserRole, captainId, captainTeamId, `dm_captain`}
		// Add a pair for the target
		captainPair := &targetPair{captainId, captainTeamId, `dm_captain`, targetId, targetTeamId, targetUserRole}
		targetList = append(targetList, captainPair)
	}

	return game.insertTargets(targetList)

}

// inserts a slice of targetPairs into the database
func (game *Game) insertTargets(targetList []*targetPair) (appErr *ApplicationError) {

	tx, err := db.Begin()

	// Prepare statement to delete previous targets
	deleteTargets, err := db.Prepare(`DELETE FROM dm_user_targets WHERE game_id = $1`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute statement to delete previous targets
	_, err = tx.Stmt(deleteTargets).Exec(game.GameId.String())
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// loop through all targets
	for _, pair := range targetList {
		// Prepare the statement to insert the target row
		insertTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id, game_id) VALUES ($1, $2, $3)`)
		if err != nil {
			tx.Rollback()
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Execute the statement to insert the target row
		_, err = tx.Stmt(insertTarget).Exec(pair.AssassinId.String(), pair.TargetId.String(), game.GameId.String())
		if err != nil {
			tx.Rollback()
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
	}

	tx.Commit()
	return nil
}

// Assign all targets plainly
func (game *Game) AssignTargets() (appErr *ApplicationError) {

	// Get new target list
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true ORDER BY random()`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var userIdBuffer, firstIdBuffer sql.NullString
	var userId, prevUserId, firstUserId uuid.UUID

	rows.Next()
	err = rows.Scan(&firstIdBuffer)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var targets []*targetPair
	firstUserId = uuid.Parse(firstIdBuffer.String)
	prevUserId = firstUserId

	// Loop through rows
	for rows.Next() {

		// Get the user_id from the row
		err = rows.Scan(&userIdBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Parse the user id
		userId = uuid.Parse(userIdBuffer.String)

		// Create a new target pair
		pair := &targetPair{prevUserId, nil, "", userId, nil, ""}

		// Append targetPair to targets
		targets = append(targets, pair)

		// Increment to the next user
		prevUserId = userId
	}

	// Set the last user to target the first
	pair := &targetPair{prevUserId, nil, "", firstUserId, nil, ""}
	targets = append(targets, pair)

	// Execute the actual insert code
	game.insertTargets(targets)
	return nil
}
