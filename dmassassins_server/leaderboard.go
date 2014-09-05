package main

type UserLeaderboardEntry struct {
	Name  string `json:"name"`
	Kills int    `json:"kills"`
}

// Returns the game leaderboard for user rankings
func (game *Game) GetUserLeaderboard(alive bool) (leaderboard []*UserLeaderboardEntry, appErr *ApplicationError) {
	// Query the db
	rows, err := db.Query(`SELECT map.kills, first_name.value as first_name, last_name.value as last_name FROM dm_user_game_mapping as map, dm_user_properties as first_name, dm_user_properties as last_name WHERE map.user_id = first_name.user_id AND map.user_id = last_name.user_id AND first_name.key = 'first_name' AND last_name.key = 'last_name' AND game_id = $1 AND alive = $2 ORDER BY kills DESC`, game.GameId.String(), alive)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Loop through leaderboard
	for rows.Next() {
		var firstName, lastName, name string
		var kills int
		err = rows.Scan(&kills, &firstName, &lastName)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Concatenate first + last name
		name = firstName + " " + lastName

		// Create the entry and append it
		entry := &UserLeaderboardEntry{name, kills}
		leaderboard = append(leaderboard, entry)
	}
	return leaderboard, nil
}
