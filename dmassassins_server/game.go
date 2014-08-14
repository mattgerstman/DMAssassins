package main

import (
	"code.google.com/p/go-uuid/uuid"
)

type Game struct {
	Game_id   string `json:"game_id"`
	Game_name string `json:"game_name"`
	Started   bool   `json:"game_started"`
}

func GetGameList() ([]*Game, *ApplicationError) {

	rows, err := db.Query(`SELECT game_id, game_name, game_started FROM dm_games ORDER BY game_name`)

	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var games []*Game
	for rows.Next() {
		var game_id, game_name string
		var started bool
		err = rows.Scan(&game_id, &game_name, &started)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		game := &Game{game_id, game_name, started}
		games = append(games, game)
	}
	return games, nil
}

func GetGameById(game_id string) (*Game, *ApplicationError) {
	var game_name string
	var started bool
	err := db.QueryRow(`SELECT game_name, started FROM dm_games WHERE game_id = $1`, game_id).Scan(&game_name, &started)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return &Game{game_id, game_name, started}, nil
}

func GetGameByName(game_name string) (*Game, *ApplicationError) {
	var game_id string
	var started bool
	err := db.QueryRow(`SELECT game_id, started FROM dm_games WHERE game_name = $1`, game_name).Scan(&game_id, &started)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return &Game{game_id, game_name, started}, nil
}

func (game *Game) End() *ApplicationError {

	res, err := db.Exec("UPDATE dm_games SET started = false WHERE game_id = $1", game.Game_id)
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
	res, err := db.Exec("UPDATE dm_games SET started = true WHERE game_id = $1", game.Game_id)
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

func NewGame(game_name, user_id string) (*Game, *ApplicationError) {
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	game_id := uuid.New()

	newGame, err := db.Prepare(`INSERT INTO dm_games (game_id, game_name) VALUES ($1, $2)`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	res, err := tx.Stmt(newGame).Exec(game_id, game_name)
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

	_, err = tx.Stmt(firstMapping).Exec(game_id, user_id)
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

	res, err = tx.Stmt(setAdmin).Exec(role, user_id)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	NoRowsAffectedAppErr = WereRowsAffected(res)
	if NoRowsAffectedAppErr != nil {
		tx.Rollback()
		return nil, NoRowsAffectedAppErr
	}
	tx.Commit()
	return &Game{game_id, game_name, false}, nil

}

// Assign all targets
func (game *Game) AssignTargets() (map[string]string, *ApplicationError) {

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
	_, err = tx.Stmt(deleteTargets).Exec(game.Game_id)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get new target list
	rows, err := db.Query(`SELECT user_id FROM dm_users WHERE game_id = $1 ORDER BY random()`, game.Game_id)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var user_id string
	var prev_user_id string
	var first_user_id string
	targets := make(map[string]string) // Map to return targets

	rows.Next()
	err = rows.Scan(&first_user_id)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	prev_user_id = first_user_id

	// Loop through rows
	for rows.Next() {

		// Get the user_id from the row
		err = rows.Scan(&user_id)
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
		_, err = tx.Stmt(insertTarget).Exec(prev_user_id, user_id, game.Game_id)
		if err != nil {
			tx.Rollback()
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Store the mapping to return
		targets[prev_user_id] = user_id
		// Increment to the next user
		prev_user_id = user_id
	}

	// Prepare the statement to have the last user target the first
	lastTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2, $3)`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to have the last user target the first
	_, err = tx.Stmt(lastTarget).Exec(user_id, first_user_id, game.Game_id)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	targets[user_id] = first_user_id

	tx.Commit()
	return targets, nil
}
