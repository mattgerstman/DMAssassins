package main

import (
	"database/sql"
	"errors"
	"github.com/getsentry/raven-go"
	"strings"
)

// Sets a single user property for a user
func (user *User) SetUserProperty(key string, value string) (*User, *ApplicationError) {
	// First attempt to update it if the property currently exists
	res, err := db.Exec(`UPDATE dm_user_properties SET value = $1 WHERE user_id = $2 AND key ILIKE $3`, value, user.User_id, key)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Check how many rows were affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// If no rows were affected insert the property
	if rowsAffected == 0 {
		res, err := db.Exec(`INSERT INTO dm_user_properties (user_id, key, value) VALUES ($1,$2,$3)`, user.User_id, key, value)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		if rowsAffected == 0 {
			err = errors.New("Failed insert for " + key + " : " + value)
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
	}

	user.Properties[key] = value
	return user, nil
}

// Get a single User Property from the db
func (user *User) GetUserProperty(key string) (string, *ApplicationError) {
	var property string
	err := db.QueryRow(`SELECT value FROM dm_user_properties WHERE user_id = $1 AND key ILIKE $2`, user.User_id, key).Scan(&property)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return property, nil
}

// Get all properties for a user
func (user *User) GetUserProperties() (map[string]string, *ApplicationError) {

	properties := make(map[string]string)

	rows, err := db.Query(`SELECT key, value FROM dm_user_properties WHERE user_id = $1`, user.User_id)
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	for rows.Next() {
		var key, value string

		err := rows.Scan(&key, &value)
		if err == nil {
			key = strings.ToLower(key)
			properties[key] = value
		} else {
			appErr := NewApplicationError("Error getting user properties", err, ErrCodeDatabase)
			LogWithSentry(appErr, map[string]string{"user_id": user.User_id}, raven.WARNING)
		}

	}
	properties["name"] = properties["first_name"] + " " + properties["last_name"]
	user.Properties = properties
	return properties, nil
}
