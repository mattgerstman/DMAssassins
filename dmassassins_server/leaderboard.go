package main

type Leaderboard struct {
	TeamsEnabled bool                             `json:"teams_enabled"`
	Users        []*UserLeaderboardEntry          `json:"users"`
	Teams        map[string]*TeamLeaderboardEntry `json:"teams"`
}

type UserLeaderboardEntry struct {
	Name     string `json:"name"`
	Kills    int    `json:"kills"`
	Alive    bool   `json:"alive"`
	TeamName string `json:"team_name"`
}

type TeamLeaderboardEntry struct {
	Kills   int `json:"kills"`
	Alive   int `json:"alive"`
	Players int `json:"players"`
}

// Get Query to generate leaderboard based on whether or not teams are enabled
func getQuery(teamsEnabled bool) (query string) {
	query = `SELECT map.kills, map.alive, first_name.value as first_name, last_name.value as last_name`
	if teamsEnabled {
		query += `, team.team_name`
	} else {
		query += `, ''`
	}
	query += ` FROM dm_user_game_mapping as map, dm_user_properties as first_name, dm_user_properties as last_name`
	if teamsEnabled {
		query += `, dm_teams as team`
	}
	query += ` WHERE map.user_id = first_name.user_id AND map.user_id = last_name.user_id AND first_name.key = 'first_name' AND last_name.key = 'last_name' AND map.game_id = $1 AND (map.user_role = 'dm_user' OR map.user_role = 'dm_captain')`
	if teamsEnabled {
		query += ` AND (team.team_id = map.team_id)`
	}
	query += ` ORDER BY kills DESC`
	return query
}

// Returns the game leaderboard for user rankings
func (game *Game) GetLeaderboard() (leaderboard *Leaderboard, appErr *ApplicationError) {

	//game.SetGameProperty("teams_enabled", "false")
	teamsEnabledString, _ := game.GetGameProperty("teams_enabled")
	teamsEnabled := teamsEnabledString == "true"

	query := getQuery(teamsEnabled)

	// Query the db
	rows, err := db.Query(query, game.GameId.String())

	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	teamKills := make(map[string]*TeamLeaderboardEntry)

	var userLeaderboard []*UserLeaderboardEntry

	// Loop through leaderboard
	for rows.Next() {
		var firstName, lastName, name, teamName string
		var kills int
		var alive bool
		err = rows.Scan(&kills, &alive, &firstName, &lastName, &teamName)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		if teamsEnabled {
			if teamKills[teamName] == nil {
				teamKills[teamName] = &TeamLeaderboardEntry{0, 0, 0}
			}
			teamKills[teamName].Kills += kills
			teamKills[teamName].Players++
			if alive {
				teamKills[teamName].Alive++
			}

		}

		// Concatenate first + last name
		name = firstName + " " + lastName

		// Create the entry and append it
		entry := &UserLeaderboardEntry{name, kills, alive, teamName}
		userLeaderboard = append(userLeaderboard, entry)
	}
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Set up the overall leaderboard object
	leaderboard = &Leaderboard{teamsEnabled, userLeaderboard, teamKills}

	return leaderboard, nil
}
