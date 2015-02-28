package main

import (
	"database/sql"
	"errors"
)

// Wrapper for SetGamePropertyTransactional that opens up and commits a transaction
func (game *Game) SetGameProperty(key string, value string) (appErr *ApplicationError) {

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Set game property
	appErr = game.SetGamePropertyTransactional(tx, key, value)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	// Commit transaction, check for errors
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Sets a single game property for a game
func (game *Game) SetGamePropertyTransactional(tx *sql.Tx, key string, value string) (appErr *ApplicationError) {
	// First attempt to update it if the property currently exists
	tryUpdate, err := db.Prepare(`UPDATE dm_game_properties SET value = $1 WHERE game_id = $2 AND key ILIKE $3`)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	res, err := tx.Stmt(tryUpdate).Exec(value, game.GameId.String(), key)
	// Check how many rows were affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// If no rows were affected insert the property
	if rowsAffected == 0 {
		tryInsert, err := db.Prepare(`INSERT INTO dm_game_properties (game_id, key, value) VALUES ($1, $2, $3)`)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		res, err = tx.Stmt(tryInsert).Exec(game.GameId.String(), key, value)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		if rowsAffected == 0 {
			err = errors.New("Failed insert for " + key + " : " + value)
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
	}

	game.Properties[key] = value
	return nil
}

// Get a single Game Property from the db
func (game *Game) GetGameProperty(key string) (property string, appErr *ApplicationError) {

	// If we have the property readily available just return it
	if property, ok := game.Properties[key]; ok {
		return property, nil
	}

	// Otherwise query the DB
	err := db.QueryRow(`SELECT value FROM dm_game_properties WHERE game_id = $1 AND key ILIKE $2`, game.GameId.String(), key).Scan(&property)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return property, nil
}

// Get all properties for a game
func (game *Game) GetGameProperties() (properties map[string]string, appErr *ApplicationError) {

	properties = make(map[string]string)

	// Query the db
	rows, err := db.Query(`SELECT key, value FROM dm_game_properties WHERE game_id = $1`, game.GameId.String())
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Loop through rows and add properties
	for rows.Next() {
		var key string
		var value string

		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		properties[key] = value

	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Add the properties to the struct
	game.Properties = properties
	return properties, nil
}
