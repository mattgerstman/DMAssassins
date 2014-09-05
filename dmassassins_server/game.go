package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"errors"
	"strings"
)

type Game struct {
	GameId      uuid.UUID `json:"game_id"`
	GameName    string    `json:"game_name"`
	Started     bool      `json:"game_started"`
	HasPassword bool      `json:"game_has_password"`
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

	// Return the game
	return &Game{gameId, gameName, gameStarted, hasPassword}, nil
}

// DROIDS Likely removing this function, i'll let it sit for a while before I do that
// func GetGameByName(gameName string) (*Game, *ApplicationError) {
// 	var gameId uuid.UUID
// 	var gameIdBuffer sql.NullString
// 	var gameStarted bool
// 	var gamePasswordBuffer sql.NullString
// 	err := db.QueryRow(`SELECT game_id, game_started, game_password FROM dm_games WHERE game_name = $1`, gameName).Scan(&gameIdBuffer, &gameStarted, &gamePasswordBuffer)
// 	if err != nil {
// 		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
// 	}
// 	gameId = uuid.Parse(gameIdBuffer.String)

// 	hasPassword := false
// 	if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
// 		hasPassword = true
// 	}

// 	return &Game{gameId, gameName, gameStarted, hasPassword}, nil
// }

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

// Start a game
func (game *Game) Start() (appErr *ApplicationError) {
	// First assign targets for the game
	_, appErr = game.AssignTargets()
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

// Get all games for a user
func (user *User) GetGamesForUser() (games []*Game, appErr *ApplicationError) {

	// Select game_ids from the dm_user_game_mapping table and use those to get the games
	rows, err := db.Query(`SELECT game.game_id, game.game_name, game.game_started, game_password FROM dm_games AS game WHERE game_id IN (SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1) ORDER BY game_name`, user.UserId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// convert the rows to an array of gamess
	return parseGameRows(rows)
}

// Get all games a user is not present in so they can join one
func (user *User) GetNewGamesForUser() (games []*Game, appErr *ApplicationError) {

	// Select game_ids from the dm_user_game_mapping table and skip those in the dm_games datable
	rows, err := db.Query(`SELECT game.game_id, game.game_name, game.game_started, game_password FROM dm_games AS game WHERE game_id NOT IN (SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1) ORDER BY game_name`, user.UserId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return parseGameRows(rows)
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
		game := &Game{gameId, gameName, gameStarted, hasPassword}
		games = append(games, game)
	}
	return games, nil
}

// Checks if a given plaintext password is right for a game
// Returns an error if the password doesn't match
func CheckPassword(gameId uuid.UUID, testPassword string) (appErr *ApplicationError) {
	var hashedPassword []byte
	err := db.QueryRow(`SELECT game_password FROM dm_games WHERE game_id = $1`, gameId.String()).Scan(&hashedPassword)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return CheckPasswordHash(hashedPassword, testPassword)
}

// Compares a hashed password to a plaintext password
// Returns an error if they don't match
func CheckPasswordHash(hashedPassword []byte, plainPw string) (appErr *ApplicationError) {
	// If they're both nil return true
	if hashedPassword == nil && plainPw == "" {
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
		return nil, NewApplicationError("Internal Error", err, ErrCodeWtf)
	}
	return hashedPassword, nil
}

// Creates a new game and saves it in the database
func NewGame(gameName string, userId uuid.UUID, gamePassword string) (game *Game, appErr *ApplicationError) {
	// Encrypt the game's password
	encryptedPassword, appErr := Crypt(gamePassword)
	if appErr != nil {
		return nil, appErr
	}

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
	gameId := uuid.NewUUID()
	res, err := tx.Stmt(newGame).Exec(gameId.String(), gameName, encryptedPassword)
	// Check to make sure the insert happened
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}

	// Prepare the statement to insert the game creator(admin) into the game
	firstMapping, err := db.Prepare(`INSERT INTO dm_user_game_mapping (game_id, user_id, user_role) VALUES ($1, $2, $3)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Executre the statement to insert the game creator(admin) into the game
	role := "dm_admin"
	res, err = tx.Stmt(firstMapping).Exec(gameId.String(), userId.String(), role)
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
	hasPassword := encryptedPassword != nil

	return &Game{gameId, gameName, false, hasPassword}, nil

}

// Assign all targets
func (game *Game) AssignTargets() (targets map[string]uuid.UUID, appErr *ApplicationError) {

	// Begin Transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare statement to delete previous targets
	deleteTargets, err := db.Prepare(`DELETE FROM dm_user_targets WHERE game_id = $1`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute statement to delete previous targets
	_, err = tx.Stmt(deleteTargets).Exec(game.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get new target list
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = 1 ORDER BY random()`, game.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var userIdBuffer, firstIdBuffer sql.NullString
	var userId, prevUserId, firstUserId uuid.UUID

	targets = make(map[string]uuid.UUID) // Map to return targets

	rows.Next()

	err = rows.Scan(&firstIdBuffer)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	firstUserId = uuid.Parse(firstIdBuffer.String)
	prevUserId = firstUserId

	// Loop through rows
	for rows.Next() {

		// Get the user_id from the row
		err = rows.Scan(&userIdBuffer)
		userId = uuid.Parse(userIdBuffer.String)
		if err != nil {
			tx.Rollback()
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Prepare the statement to insert the target row
		insertTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id, game_id) VALUES ($1, $2, $3)`)
		if err != nil {
			tx.Rollback()
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Execute the statement to insert the target row
		_, err = tx.Stmt(insertTarget).Exec(prevUserId.String(), userId.String(), game.GameId.String())
		if err != nil {
			tx.Rollback()
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Store the mapping to return
		targets[prevUserId.String()] = userId
		// Increment to the next user
		prevUserId = userId
	}

	// Prepare the statement to have the last user target the first
	lastTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id, game_id) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to have the last user target the first
	_, err = tx.Stmt(lastTarget).Exec(userId.String(), firstUserId.String(), game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	targets[userId.String()] = firstUserId

	tx.Commit()
	return targets, nil
}
