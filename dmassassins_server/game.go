package main

import (
	"code.google.com/p/go-uuid/uuid"
)

type Game struct {
	GameId   uuid.UUID `json:"game_id"`
	GameName string    `json:"game_name"`
	Started  bool      `json:"game_started"`
}

func GetGameList() ([]*Game, *ApplicationError) {

	rows, err := db.Query(`SELECT game_id, game_name, game_started FROM dm_games ORDER BY game_name`)

	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var games []*Game
	for rows.Next() {
		var gameId uuid.UUID
		var game_name string
		var started bool
		err = rows.Scan(&gameId, &game_name, &started)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		game := &Game{gameId, game_name, started}
		games = append(games, game)
	}
	return games, nil
}

func GetGameById(gameId uuid.UUID) (*Game, *ApplicationError) {
	var game_name string
	var started bool
	err := db.QueryRow(`SELECT game_name, started FROM dm_games WHERE game_id = $1`, gameId).Scan(&game_name, &started)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return &Game{gameId, game_name, started}, nil
}

func GetGameByName(game_name string) (*Game, *ApplicationError) {
	var gameId uuid.UUID
	var started bool
	err := db.QueryRow(`SELECT game_id, started FROM dm_games WHERE game_name = $1`, game_name).Scan(&gameId, &started)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return &Game{gameId, game_name, started}, nil
}

func (game *Game) End() *ApplicationError {

	res, err := db.Exec("UPDATE dm_games SET started = false WHERE game_id = $1", game.GameId)
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
	res, err := db.Exec("UPDATE dm_games SET started = true WHERE game_id = $1", game.GameId)
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

func (game *Game) GetLeaderBoard(alive bool) {
	_ := db.Query(`SELECT map.user_id, map.kills, first_name.value as first_name, last_name.value as last_name FROM dm_user_game_mapping as map, dm_user_properties as dm_first_name, user_properties as last_name WHERE map.user_id = first_name.user_id AND map.user_id = last_name.user_id AND first_name.key = 'first_name' AND last_name.key = 'last_name' game_id = $1 AND alive = $2 ORDER BY kills`, game.GameId, alive)
	return
}

func NewGame(game_name string, userId uuid.UUID) (*Game, *ApplicationError) {
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	gameId := uuid.NewUUID()

	newGame, err := db.Prepare(`INSERT INTO dm_games (game_id, game_name) VALUES ($1, $2)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	res, err := tx.Stmt(newGame).Exec(gameId, game_name)
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

	_, err = tx.Stmt(firstMapping).Exec(gameId, userId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	setAdmin, err := db.Prepare(`UPDATE dm_users SET user_role = $1 WHERE user_id = $2`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	role := "dm_admin"

	res, err = tx.Stmt(setAdmin).Exec(role, userId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr = WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}
	tx.Commit()
	return &Game{gameId, game_name, false}, nil

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
	_, err = tx.Stmt(deleteTargets).Exec(game.GameId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get new target list
	rows, err := db.Query(`SELECT user_id FROM dm_users WHERE game_id = $1 ORDER BY random()`, game.GameId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var userId uuid.UUID
	var prevUserId uuid.UUID
	var firstUserId uuid.UUID
	targets := make(map[string]uuid.UUID) // Map to return targets

	rows.Next()
	err = rows.Scan(&firstUserId)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	prevUserId = firstUserId

	// Loop through rows
	for rows.Next() {

		// Get the user_id from the row
		err = rows.Scan(&userId)
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
		_, err = tx.Stmt(insertTarget).Exec(prevUserId, userId, game.GameId)
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
	lastTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to have the last user target the first
	_, err = tx.Stmt(lastTarget).Exec(userId, firstUserId, game.GameId)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	targets[userId.String()] = firstUserId

	tx.Commit()
	return targets, nil
}
