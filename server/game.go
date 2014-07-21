package main

func assignTargets() (bool, *ApplicationError) {
	var appErr *ApplicationError
	msg := "Interal Error"
	code := ERROR_DATABASE

	transaction, err := db.Begin()
	appErr = CheckError(msg, err, code)
	if appErr != nil {
		return false, appErr
	}
	defer transaction.Commit()

	_, err = db.Exec(`DELETE FROM dm_user_targets`)
	appErr = CheckError(msg, err, code)
	if appErr != nil {
		return false, appErr
	}

	rows, err := db.Query(`SELECT user_id FROM dm_users ORDER BY random()`)
	appErr = CheckError(msg, err, code)
	if appErr != nil {
		return false, appErr
	}

	var user_id string
	var prev_user_id string
	var first_user_id string

	rows.Next()
	rows.Scan(&first_user_id)
	prev_user_id = first_user_id

	for rows.Next() {
		rows.Scan(&user_id)
		_, err = db.Exec(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, prev_user_id, user_id)

		appErr = CheckError(msg, err, code)
		if appErr != nil {
			return false, appErr
		}
		prev_user_id = user_id
	}

	_, err = db.Exec(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, user_id, first_user_id)
	appErr = CheckError(msg, err, code)
	if appErr != nil {
		return false, appErr
	}
	return true, nil
}
