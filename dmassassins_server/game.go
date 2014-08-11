package main

import (
	"code.google.com/p/go-uuid/uuid"
)

type Game struct {
	Game_id   string `json:"game_id"`
	Game_name string `json:"game_name"`
	Started   bool   `json:"game_name"`
}

func (game *Game) StartGame() {
		targets, appErr := game.AssignTargets();
		if (appErr != nil) {
			return appErr;
		}
		db.Query("UPDATE dm_games SET started = true WHERE game_id = $1", game.Game_id);
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

	_, err = tx.Stmt(newGame).Exec(game_id, game_name)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
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

	_, err = tx.Stmt(setAdmin).Exec(user_id, role)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	tx.Commit()
	return &Game{game_id, game_name, false}, nil

}

// Assign all targets
func AssignTargets() (map[string]string, *ApplicationError) {

	// Begin Transaction
	tx, err := db.Begin()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Prepare statement to delete previous targets
	deleteTargets, err := db.Prepare(`DELETE FROM dm_user_targets`)
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute statement to delete previous targets
	_, err = tx.Stmt(deleteTargets).Exec()
	if err != nil {
		tx.Rollback()
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Get new target list
	rows, err := db.Query(`SELECT user_id FROM dm_users ORDER BY random()`)
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
		insertTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`)
		if err != nil {
			tx.Rollback()
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Execute the statement to insert the target row
		_, err = tx.Stmt(insertTarget).Exec(prev_user_id, user_id)
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
	lastTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute the statement to have the last user target the first
	_, err = tx.Stmt(lastTarget).Exec(user_id, first_user_id)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	targets[user_id] = first_user_id

	tx.Commit()
	return targets, nil
}
