package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
)

type Game struct {
	GameId      uuid.UUID `json:"game_id"`
	GameName    string    `json:"game_name"`
	Started     bool      `json:"game_started"`
	HasPassword bool      `json:"game_has_password"`
}

func GetGameList() ([]*Game, *ApplicationError) {

	rows, err := db.Query(`SELECT game_id, game_name, game_started, game_password FROM dm_games ORDER BY game_name`)

	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var games []*Game
	for rows.Next() {
		var gameId uuid.UUID
		var gameIdBuffer, gamePasswordBuffer sql.NullString
		var gameName string
		var gameStarted bool
		err = rows.Scan(&gameIdBuffer, &gameName, &gameStarted, &gamePasswordBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		hasPassword := false
		if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
			hasPassword = true
		}

		gameId = uuid.Parse(gameIdBuffer.String)
		game := &Game{gameId, gameName, gameStarted, hasPassword}
		games = append(games, game)
	}
	return games, nil
}

func GetGameById(gameId uuid.UUID) (*Game, *ApplicationError) {
	var gameName string
	var gameStarted bool
	var gamePasswordBuffer sql.NullString
	err := db.QueryRow(`SELECT game_name, game_started, game_password FROM dm_games WHERE game_id = $1`, gameId.String()).Scan(&gameName, &gameStarted, &gamePasswordBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	hasPassword := false
	if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
		hasPassword = true
	}

	return &Game{gameId, gameName, gameStarted, hasPassword}, nil
}

func GetGameByName(gameName string) (*Game, *ApplicationError) {
	var gameId uuid.UUID
	var gameIdBuffer sql.NullString
	var gameStarted bool
	var gamePasswordBuffer sql.NullString
	err := db.QueryRow(`SELECT game_id, game_started, game_password FROM dm_games WHERE game_name = $1`, gameName).Scan(&gameIdBuffer, &gameStarted, &gamePasswordBuffer)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	gameId = uuid.Parse(gameIdBuffer.String)

	hasPassword := false
	if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
		hasPassword = true
	}

	return &Game{gameId, gameName, gameStarted, hasPassword}, nil
}

func (game *Game) End() *ApplicationError {

	res, err := db.Exec("UPDATE dm_games SET game_started = false WHERE game_id = $1", game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}
	game.Started = false
	return nil
}

func (game *Game) Start() *ApplicationError {
	_, appErr := game.AssignTargets()
	if appErr != nil {
		return appErr
	}
	res, err := db.Exec("UPDATE dm_games SET game_started = true WHERE game_id = $1", game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		return NoRowsAffectedAppErr
	}
	game.Started = true
	return nil
}

func (user *User) GetGamesForUser() ([]*Game, *ApplicationError) {

	rows, err := db.Query(`SELECT game.game_id, game.game_name, game.game_started, game_password FROM dm_games AS game WHERE game_id IN (SELECT game_id FROM dm_user_game_mapping WHERE user_id = $1) ORDER BY game_name`, user.UserId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	var games []*Game
	for rows.Next() {
		var gameId uuid.UUID
		var gameIdBuffer, gamePasswordBuffer sql.NullString
		var gameName string
		var gameStarted bool

		err = rows.Scan(&gameIdBuffer, &gameName, &gameStarted, &gamePasswordBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		gameId = uuid.Parse(gameIdBuffer.String)

		hasPassword := false
		if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
			hasPassword = true
		}

		game := &Game{gameId, gameName, gameStarted, hasPassword}
		games = append(games, game)
	}
	return games, nil
}

func NewGame(gameName string, userId uuid.UUID, gamePassword string) (*Game, *ApplicationError) {
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.NewUUID()
	var gamePasswordBuffer sql.NullString
	if gamePassword != "" {
		gamePasswordBuffer.String = gamePassword
		gamePasswordBuffer.Valid = true
	}

	newGame, err := db.Prepare(`INSERT INTO dm_games (game_id, game_name, game_password) VALUES ($1, $2, $3)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	res, err := tx.Stmt(newGame).Exec(gameId.String(), gameName, gamePasswordBuffer)
	NoRowsAffectedAppErr := WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}

	firstMapping, err := db.Prepare(`INSERT INTO dm_user_game_mapping (game_id, user_id) VALUES ($1, $2)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	_, err = tx.Stmt(firstMapping).Exec(gameId.String(), userId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	setAdmin, err := db.Prepare(`UPDATE dm_user_game_mapping SET user_role = $1 WHERE user_id = $2 AND game_id = $3`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	role := "dm_admin"

	res, err = tx.Stmt(setAdmin).Exec(role, userId.String(), gameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr = WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}
	tx.Commit()
	hasPassword := false
	if (gamePasswordBuffer.Valid != false) && (gamePasswordBuffer.String != "") {
		hasPassword = true
	}
	return &Game{gameId, gameName, false, hasPassword}, nil

}

// Assign all targets
func (game *Game) AssignTargets() (map[string]uuid.UUID, *ApplicationError) {

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
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 ORDER BY random()`, game.GameId.String())
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var userIdBuffer, firstIdBuffer sql.NullString
	var userId, prevUserId, firstUserId uuid.UUID

	targets := make(map[string]uuid.UUID) // Map to return targets

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
