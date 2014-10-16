package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"math/rand"
	"time"
)

// Assigns targets using a methodology
func (game *Game) AssignTargetsBy(assignmentType string) (appErr *ApplicationError) {

	// Reverse targets
	if assignmentType == `reverse` {
		return game.ReverseTargets()
	}

	if assignmentType == `strong_weak` {
		return game.AssignStrongTargetWeak()
	}

	// If they don't have a special method check if teams are enabled
	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return appErr
	}

	// If teams are enabled assigned by team
	if teamsEnabled == `true` {
		canAssign, appErr := game.CanAssignByTeams()
		if appErr != nil {
			return appErr
		}

		if canAssign {
			// Get all players in the game
			rows, appErr := game.GetAllActivePlayersAsRows()
			if appErr != nil {
				return appErr
			}
			return game.AssignTargetsByTeams(rows)
		}
	}

	// Get players to assign targets with
	users, appErr := game.GetAllActivePlayersAsUUIDSlice()
	if appErr != nil {
		return appErr
	}

	// Fallback to plain random assignment
	return game.AssignTargets(users)
}

type targetPair struct {
	AssassinId       uuid.UUID `json:assassin_id`
	AssassinTeamId   uuid.UUID `json:assassin_team_id`
	AssassinUserRole string    `json:assassin_user_role`
	TargetId         uuid.UUID `json:target_id`
	TargetTeamId     uuid.UUID `json:target_team_id`
	TargetUserRole   string    `json:target_user_role`
}

// Reverses targets for a game
func (game *Game) ReverseTargets() (appErr *ApplicationError) {
	var userIdBuffer, targetIdBuffer string
	var targets []*targetPair
	rows, err := db.Query(`SELECT user_id, target_id FROM dm_user_targets WHERE game_id = $1`, game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	for rows.Next() {
		rows.Scan(&userIdBuffer, &targetIdBuffer)
		userId := uuid.Parse(userIdBuffer)
		targetId := uuid.Parse(targetIdBuffer)
		pair := &targetPair{targetId, nil, "", userId, nil, ""}
		targets = append(targets, pair)
	}

	return game.insertTargets(targets)
}

// Put strong users in a closed loop and other users in a regular loop
func (game *Game) AssignClosedStrongLoop() (appErr *ApplicationError) {
	strong, appErr := game.getStrongPlayers()
	if appErr != nil {
		return appErr
	}

	// DROIDS finish this loop
	_ = strong

	return nil
}

// Get a slice of the strongest player for each team, ties are broken arbitrarily
func (game *Game) getStrongPlayers() (strong []uuid.UUID, appErr *ApplicationError) {
	// segregate strong users
	rows, err := db.Query(`SELECT DISTINCT ON (team_id) user_id FROM dm_user_game_mapping WHERE game_id = $1 ORDER BY team_id, kills desc`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// parse strong users into slice
	for rows.Next() {
		var strongUserIdBuffer string
		err = rows.Scan(&strongUserIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Add strong userId to strong slice
		strongUserId := uuid.Parse(strongUserIdBuffer)
		strong = append(strong, strongUserId)
	}
	return strong, nil
}

// Get a slice of the weakest player for each team, ties are broken arbitrarily
func (game *Game) getWeakPlayers() (weak []uuid.UUID, appErr *ApplicationError) {
	// segregate weak users
	rows, err := db.Query(`SELECT DISTINCT ON (team_id) user_id FROM dm_user_game_mapping WHERE game_id = $1 ORDER BY team_id, kills asc`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// We need to offset strong and weak to ensure that strong players aren't targetting their own weakest palyer
	var firstWeakUserIdBuffer string
	rows.Next()
	rows.Scan(&firstWeakUserIdBuffer)
	firstWeakUserId := uuid.Parse(firstWeakUserIdBuffer)

	// parse weak users into slice
	for rows.Next() {
		var weakUserIdBuffer string
		err = rows.Scan(&weakUserIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Add weak userId to weak slice
		weakUserId := uuid.Parse(weakUserIdBuffer)
		weak = append(weak, weakUserId)
	}
	weak = append(weak, firstWeakUserId)
	return weak, nil
}

// Gets a list of rows for players not in the given uuid slice
func (game *Game) getPlayersNotInSlice(userSlice []interface{}) (rows *sql.Rows, appErr *ApplicationError) {

	params := GetParamsForSlice(2, userSlice)
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') AND user_id NOT IN (`+params+`)`, userSlice...)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return rows, nil
}

// have the strongest players target the weakest ones, don't be concerned about teams/captains
func (game *Game) AssignStrongTargetWeak() (appErr *ApplicationError) {
	// Get strong players
	strong, appErr := game.getStrongPlayers()
	if appErr != nil {
		return appErr
	}
	// Get weak players
	weak, appErr := game.getWeakPlayers()
	if appErr != nil {
		return appErr
	}

	// Select other users from the db
	var strongWeak []interface{}
	strongWeak = append(strongWeak, game.GameId.String())
	for _, strongId := range strong {
		strongWeak = append(strongWeak, strongId.String())
	}
	for _, weakId := range weak {
		strongWeak = append(strongWeak, weakId.String())
	}

	// Get players that aren't in the strong or weak category
	rows, appErr := game.getPlayersNotInSlice(strongWeak)
	if appErr != nil {
		return appErr
	}

	// Get number of players
	numPlayers, appErr := game.GetNumActivePlayers()
	if appErr != nil {
		return appErr
	}

	// Get number of teams
	numTeams := len(weak)

	// Determine how often to insert a strong/weak pair
	insertNum := (numPlayers / numTeams) - 1

	// Create first strong/weak pair and make them the first target
	var targets []*targetPair
	firstStrong, strong := strong[0], strong[1:]
	firstWeak, weak := weak[0], weak[1:]
	firstPair := &targetPair{firstStrong, nil, "", firstWeak, nil, ""}

	targets = append(targets, firstPair)

	// we already have one pair
	i := 1
	lastUserId := firstWeak
	for rows.Next() {
		// Get userId
		var userIdBuffer string
		err := rows.Scan(&userIdBuffer)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Create target pair
		userId := uuid.Parse(userIdBuffer)
		pair := &targetPair{lastUserId, nil, "", userId, nil, ""}

		// Add pair to targets slice
		targets = append(targets, pair)

		// Set last userId and increment
		lastUserId = userId
		i++

		// Unless we're at the right point to insert just continue
		if (i % insertNum) != 0 {
			continue
		}

		// If we're out of strong/weak make it so we won't get here again
		if len(strong) == 0 {
			insertNum = numPlayers + 1
			continue
		}

		var nextStrong, nextWeak uuid.UUID
		// Get next strong/weak pair
		nextStrong, strong = strong[0], strong[1:]
		nextWeak, weak = weak[0], weak[1:]

		// Set up next strong/weak pair and insert it
		pair = &targetPair{lastUserId, nil, "", nextStrong, nil, ""}
		strongWeakPair := &targetPair{nextStrong, nil, "", nextWeak, nil, ""}
		targets = append(targets, pair, strongWeakPair)

		lastUserId = nextWeak
		i += 2
	}

	// have last user target first strong user
	lastTarget := &targetPair{lastUserId, nil, "", firstWeak, nil, ""}
	targets = append(targets, lastTarget)

	return nil
}

func (game *Game) GetAllActivePlayersAsRows() (rows *sql.Rows, appErr *ApplicationError) {
	// Get new target list
	rows, err := db.Query(`SELECT user_id, team_id, user_role FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY random()`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return rows, nil
}

// Assign targets and space them out by team
func (game *Game) AssignTargetsByTeams(rows *sql.Rows) (appErr *ApplicationError) {
	// Get the list of team ids
	teamsList, appErr := game.GetActiveTeamIds()
	if appErr != nil {
		return appErr
	}

	// Create userList, captainList, and buffer variables
	var userIdBuffer, teamIdBuffer, userRole string
	userList := make(map[string][]uuid.UUID)
	captainList := make(map[string]uuid.UUID)

	// Fill in userList and captainList with valid slices
	for _, team := range teamsList {
		userList[team.String()] = []uuid.UUID{}
	}

	// parse out users into teams and captain
	numUsers := 0
	for rows.Next() {
		// Get the user_id from the row
		err := rows.Scan(&userIdBuffer, &teamIdBuffer, &userRole)
		if err != nil {
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// If user is a captain put them in the captains list
		userId := uuid.Parse(userIdBuffer)
		if CompareRole(userRole, RoleCaptain) {
			captainList[teamIdBuffer] = userId
			continue
		}
		// If they aren't a captain add them to the userList
		userList[teamIdBuffer] = append(userList[teamIdBuffer], userId)
		numUsers++
	}

	numTeams := len(teamsList)
	originalNumTeams := numTeams
	rand.Seed(time.Now().UTC().UnixNano())

	assigned := 1
	var targetList []*targetPair
	randTeam := teamsList[rand.Intn(numTeams)]

	var firstUser, firstUserTeam, currentUser, currentUserTeam, lastUser, lastUserTeam uuid.UUID

	firstUser, userList[randTeam.String()] = userList[randTeam.String()][0], userList[randTeam.String()][1:]

	firstUserTeam = randTeam
	lastUserTeam = firstUserTeam
	lastUser = firstUser

	// While we still have users to assign
	for assigned < numUsers {
		// Get a random team
		randTeamIndex := rand.Intn(numTeams)
		currentUserTeam = teamsList[randTeamIndex]

		// If the random team is the same as the last one just go to the next team
		if uuid.Equal(currentUserTeam, lastUserTeam) {
			randTeamIndex++
			if randTeamIndex >= numTeams {
				randTeamIndex = 0
			}
			currentUserTeam = teamsList[randTeamIndex]
		}

		// If our current team has no members delete it and go to the next one
		for len(userList[currentUserTeam.String()]) == 0 {
			// delete the current team from the userlist
			delete(userList, currentUserTeam.String())

			// delete the currentTeam from the teamsList
			teamsList = append(teamsList[:randTeamIndex], teamsList[(randTeamIndex+1):]...)

			// change numTeams to reflect the current number of teams
			numTeams = len(teamsList)
			// If our index is greater than the number of teams set it to 0
			if randTeamIndex >= numTeams {
				randTeamIndex = 0
			}

			// set the current user team
			currentUserTeam = teamsList[randTeamIndex]
		}

		// pop out currentUser from userList[currentUserTeam]
		currentUser, userList[currentUserTeam.String()] = userList[currentUserTeam.String()][0], userList[currentUserTeam.String()][1:]

		// append lastUser currentUser to the targetsList
		userTargetPair := &targetPair{lastUser, lastUserTeam, `dm_user`, currentUser, currentUserTeam, `dm_user`}
		targetList = append(targetList, userTargetPair)

		lastUser = currentUser
		lastUserTeam = currentUserTeam

		assigned++
	}

	userTargetPair := &targetPair{lastUser, lastUserTeam, `dm_user`, firstUser, firstUserTeam, `dm_user`}
	targetList = append(targetList, userTargetPair)

	// Space out captains by 5 unless we don't have enough users
	captainSpace := 5
	numTeams = originalNumTeams
	userTeamRatio := (numUsers / numTeams)
	if userTeamRatio < captainSpace {
		captainSpace = userTeamRatio
	}

	i := 0
	triesForRole := 5
	triesForTeam := 10
	// Loop through captains
	for captainTeam, captainId := range captainList {
		j := 0
		captainTeamId := uuid.Parse(captainTeam)
		// we have to use a marker to determine if we've found an appropriate pair to insert the captain in
		foundPair := false
		for !foundPair {
			if i >= numUsers {
				i = 0
				j++
			}

			// If we've lapped users enough times allow users from the same team to target each other
			if j > triesForTeam {
				foundPair = true
				continue
			}

			// Check if the assassin and target both have different teams than the captain
			assassinTeamId, targetTeamId := targetList[i].AssassinTeamId, targetList[i].TargetTeamId
			assassinUserRole, targetUserRole := targetList[i].AssassinUserRole, targetList[i].TargetUserRole
			if uuid.Equal(captainTeamId, assassinTeamId) {
				i++
				continue
			}
			if uuid.Equal(captainTeamId, targetTeamId) {
				i++
				continue
			}

			// If we've lapped users enough times allow captains to target each other
			if j > triesForRole {
				continue
			}

			if assassinUserRole == `dm_captain` {
				i++
				continue
			}

			if targetUserRole == `dm_captain` {
				i++
				continue
			}

			// If neither has the same team, we've found a match
			foundPair = true
		}
		// Get all the assassin/target information
		assassinId, assassinTeamId, assassinUserRole, targetId, targetTeamId, targetUserRole := targetList[i].AssassinId, targetList[i].AssassinTeamId, targetList[i].AssassinUserRole, targetList[i].TargetId, targetList[i].TargetTeamId, targetList[i].TargetUserRole
		// Insert the captain after the assassin
		targetList[i] = &targetPair{assassinId, assassinTeamId, assassinUserRole, captainId, captainTeamId, `dm_captain`}
		// Add a pair for the target
		captainPair := &targetPair{captainId, captainTeamId, `dm_captain`, targetId, targetTeamId, targetUserRole}
		targetList = append(targetList, captainPair)
	}

	return game.insertTargets(targetList)

}

// inserts a slice of targetPairs into the database
func (game *Game) insertTargets(targetList []*targetPair) (appErr *ApplicationError) {

	tx, err := db.Begin()

	// Prepare statement to delete previous targets
	deleteTargets, err := db.Prepare(`DELETE FROM dm_user_targets WHERE game_id = $1`)
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute statement to delete previous targets
	_, err = tx.Stmt(deleteTargets).Exec(game.GameId.String())
	if err != nil {
		tx.Rollback()
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// loop through all targets
	for _, pair := range targetList {
		// Prepare the statement to insert the target row
		insertTarget, err := db.Prepare(`INSERT INTO dm_user_targets (user_id, target_id, game_id) VALUES ($1, $2, $3)`)
		if err != nil {
			tx.Rollback()
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}

		// Execute the statement to insert the target row
		_, err = tx.Stmt(insertTarget).Exec(pair.AssassinId.String(), pair.TargetId.String(), game.GameId.String())
		if err != nil {
			tx.Rollback()
			return NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
	}

	tx.Commit()
	return nil
}

func (game *Game) GetAllActivePlayersAsUUIDSlice() (users []uuid.UUID, appErr *ApplicationError) {
	// Get new target list
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true ORDER BY random()`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	var userIdBuffer string
	// Loop through rows
	for rows.Next() {

		// Get the user_id from the row
		err = rows.Scan(&userIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Parse the user id
		userId := uuid.Parse(userIdBuffer)

		// Add to user array
		users = append(users, userId)
	}
	return users, nil
}

// Assign all targets plainly
func (game *Game) AssignTargets(users []uuid.UUID) (appErr *ApplicationError) {

	var targets []*targetPair
	var prevUserId uuid.UUID
	// Loop through rows
	for _, userId := range users {
		// Create a new target pair
		pair := &targetPair{prevUserId, nil, "", userId, nil, ""}

		// Append targetPair to targets
		targets = append(targets, pair)

		// Increment to the next user
		prevUserId = userId
	}

	// Set the last user to target the first
	pair := &targetPair{prevUserId, nil, "", users[0], nil, ""}
	targets = append(targets, pair)

	// Execute the actual insert code
	game.insertTargets(targets)
	return nil
}
