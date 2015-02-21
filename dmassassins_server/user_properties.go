package main

import (
	"database/sql"
	"errors"
	"strings"
)

// Sets a single user property for a user, takes a transaction to go into another function
func (user *User) SetUserPropertyTransactional(tx *sql.Tx, key string, value string) (appErr *ApplicationError) {
	// First attempt to update it if the property currently exists
	tryUpdate, err := db.Prepare(`UPDATE dm_user_properties SET value = $1 WHERE user_id = $2 AND key ILIKE $3`)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	res, err := tx.Stmt(tryUpdate).Exec(value, user.UserId.String(), key)
	// Check how many rows were affected by the update
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// If no rows were affected insert the property
	if rowsAffected == 0 {
		tryInsert, err := db.Prepare(`INSERT INTO dm_user_properties (user_id, key, value) VALUES ($1,$2,$3)`)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		res, err = tx.Stmt(tryInsert).Exec(user.UserId.String(), key, value)
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

	user.Properties[key] = value
	return nil
}

// Sets a single user property for a user, wraps it in a transaction to make sure it doesn't blow up
func (user *User) SetUserProperty(key string, value string) (appErr *ApplicationError) {

	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	appErr = user.SetUserPropertyTransactional(tx, key, value)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	// Check transaction for errors and commit
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Sets a map of new user properties
func (user *User) SetUserProperties(newProperties map[string]string) (appErr *ApplicationError) {
	// Start a transaction so we can rollback if something blows up
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Loop through new properites and set them all
	for key, value := range newProperties {
		appErr = user.SetUserPropertyTransactional(tx, key, value)
		if appErr != nil {
			tx.Rollback()
			return appErr
		}
	}

	// Check transaction for errors and commit
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return nil
}

// Get a user property and cast it as a boolean
func (user *User) GetUserPropertyBool(key string) (property bool, appErr *ApplicationError) {
	stringProperty, appErr := user.GetUserProperty(key)
	return stringProperty == `true`, appErr
}

// Get a single User Property from the db
func (user *User) GetUserProperty(key string) (property string, appErr *ApplicationError) {

	// If we have the property in the user struct just return it
	if property, ok := user.Properties[key]; ok {
		return property, nil
	}

	err := db.QueryRow(`SELECT value FROM dm_user_properties WHERE user_id = $1 AND key ILIKE $2`, user.UserId.String(), key).Scan(&property)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return property, nil
}

// Get all properties for a user
func (user *User) GetUserProperties() (properties map[string]string, appErr *ApplicationError) {

	properties = make(map[string]string)

	// Query the db
	rows, err := db.Query(`SELECT key, value FROM dm_user_properties WHERE user_id = $1`, user.UserId.String())
	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Loop through rows and add properties
	for rows.Next() {
		var key, value string

		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		key = strings.ToLower(key)
		properties[key] = value

	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	properties["name"] = properties["first_name"] + " " + properties["last_name"]
	user.Properties = properties
	return properties, nil
}
