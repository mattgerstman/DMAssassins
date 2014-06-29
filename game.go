package main

func assignTargets() bool {
	transaction, _ := db.Begin()
	defer transaction.Commit()

	rows, _ := db.Query(`DELETE FROM dm_user_targets`)
	rows, _ = db.Query(`SELECT user_id FROM dm_users ORDER BY random()`)

	var user_id string
	var prev_user_id string
	var first_user_id string

	rows.Next()
	rows.Scan(&first_user_id)
	prev_user_id = first_user_id
	for rows.Next() {
		rows.Scan(&user_id)
		db.Query(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, prev_user_id, user_id)
		prev_user_id = user_id
	}
	db.Query(`INSERT INTO dm_user_targets (user_id, target_id) VALUES ($1, $2)`, user_id, first_user_id)
	return true
}
