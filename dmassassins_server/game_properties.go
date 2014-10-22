package main

import (
	"database/sql"
	"errors"
	"github.com/getsentry/raven-go"
	"strings"
)

// Sets a single game property for a game
func (game *Game) SetGameProperty(key string, value string) (appErr *ApplicationError) {

	// First attempt to update it if the property currently exists
	res, err := db.Exec(`UPDATE dm_game_properties SET value = $1 WHERE game_id = $2 AND key ILIKE $3`, value, game.GameId.String(), key)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Check how many rows were affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// If no rows were affected insert the property
	if rowsAffected == 0 {
		res, err := db.Exec(`INSERT INTO dm_game_properties (game_id, key, value) VALUES ($1,$2,$3)`, game.GameId.String(), key, value)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		if rowsAffected == 0 {
			err = errors.New("Failed insert for " + key + " : " + string(value))
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
	}

	// Set the property in the struct
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
		if err == nil {
			key = strings.ToLower(key)
			properties[key] = value
		} else {
			// Fail silently if a single property spazzes out (should never happen but who knows)
			appErr := NewApplicationError("Error getting game properties", err, ErrCodeDatabase)
			LogWithSentry(appErr, map[string]string{"game_id": game.GameId.String()}, raven.WARNING)
		}
	}

	// Add the properties to the struct
	game.Properties = properties
	return properties, nil
}
