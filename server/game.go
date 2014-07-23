package main

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
