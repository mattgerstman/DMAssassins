package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Assigns targets using a methodology, wrapps the inner function in a transaction
func (game *Game) AssignTargetsBy(assignmentType string) (appErr *ApplicationError) {
	// begin transaction
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	// Run internal assign targets
	appErr = game.AssignTargetsByTransactional(tx, assignmentType)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}
	// error check the transaction commit
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return nil
}

// Assigns targets using a methodology using a transaction
func (game *Game) AssignTargetsByTransactional(tx *sql.Tx, assignmentType string) (appErr *ApplicationError) {

	anyLeft, appErr := game.checkAlive()
	if appErr != nil {
		return appErr
	}

	if !anyLeft || assignmentType == `delete` {
		return game.DeleteTargetsTransactional(tx)
	}

	// Reverse targets
	if assignmentType == `reverse` {
		fmt.Println(`reverse`)
		return game.reverseTargets(tx)
	}

	if assignmentType == `strong_weak` {
		fmt.Println(`strong_weak`)
		return game.assignStrongTargetWeak(tx)
	}

	if assignmentType == `strong_closed` {
		fmt.Println(`closed_strong`)
		return game.assignClosedStrongLoop(tx)
	}

	// If they don't have a special method check if teams are enabled
	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if appErr != nil {
		return appErr
	}

	// If teams are enabled assigned by team
	if teamsEnabled == `true` {
		fmt.Println(`teams`)
		return game.assignTargetsByTeams(tx)
	}
	fmt.Println(`regular`)
	// Get players to assign targets with
	users, appErr := game.GetAllActivePlayersAsUUIDSlice()
	if appErr != nil {
		return appErr
	}

	// Fallback to plain random assignment
	return game.assignTargets(tx, users, false)
}

func (game *Game) checkAlive() (anyLeft bool, appErr *ApplicationError) {
	var numLeft int
	err := db.QueryRow(`SELECT count(*) FROM dm_user_game_mapping WHERE user_role != 'dm_admin' AND user_role != 'dm_super_admin' AND alive = true AND game_id = $1`, game.GameId.String()).Scan(&numLeft)
	if err != nil {
		return false, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return (numLeft > 0), nil
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
func (game *Game) reverseTargets(tx *sql.Tx) (appErr *ApplicationError) {
	var userIdBuffer, targetIdBuffer string
	var targets []*targetPair
	rows, err := tx.Query(`SELECT user_id, target_id FROM dm_user_targets WHERE game_id = $1`, game.GameId.String())
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
	err = rows.Close()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	return game.insertTargetsWithDelete(tx, targets)
}

// Put strong users in a closed loop and other users in a regular loop
func (game *Game) assignClosedStrongLoop(tx *sql.Tx) (appErr *ApplicationError) {
	// Gets the strongest players
	strong, appErr := game.GetStrongPlayers()
	if appErr != nil {
		return appErr
	}

	// Assigns the strong players as targets
	appErr = game.assignTargets(tx, strong, false)
	if appErr != nil {
		return appErr
	}

	// Converts strong players to an interface
	strongInterface := ConvertUUIDSliceToInterface(strong)

	// Gets players not in the strong slice
	regularPlayers, appErr := game.getPlayersNotInSlice(strongInterface)
	if appErr != nil {
		return appErr
	}

	// Converts the rows from get players not in slice to a slice
	users, appErr := ConvertUserIdRowsToSlice(regularPlayers)
	if appErr != nil {
		return appErr
	}

	return game.assignTargets(tx, users, true)
}

// Gets a list of rows for players not in the given uuid slice
func (game *Game) getPlayersNotInSlice(userSlice []interface{}) (rows *sql.Rows, appErr *ApplicationError) {

	var gameSlice []interface{}
	gameSlice = append(gameSlice, game.GameId.String())
	gameSlice = append(gameSlice, userSlice...)

	params := GetParamsForSlice(2, userSlice)
	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') AND user_id NOT IN (`+params+`) ORDER BY random()`, gameSlice...)
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return rows, nil
}

// have the strongest players target the weakest ones, don't be concerned about teams/captains
func (game *Game) assignStrongTargetWeak(tx *sql.Tx) (appErr *ApplicationError) {
	// Get strong players
	strong, appErr := game.GetStrongPlayers()
	if appErr != nil {
		return appErr
	}
	// Get weak players
	weak, appErr := game.GetWeakPlayers()
	if appErr != nil {
		return appErr
	}

	// if we don't have even strong and weak players we're in trouble
	if len(strong) != len(weak) {
		return NewApplicationError("Internal Error", errors.New(`Strong List and Weak List have different lengths`), ErrCodeDatabase)
	}

	// Get number of teams
	numTeams := len(weak)

	for i := 0; i < numTeams; i++ {
		if !uuid.Equal(strong[i], weak[i]) {
			continue
		}
		numTeams--
		next := i + 1
		if next == numTeams {
			strong = strong[:i]
			weak = weak[:i]
			i--
			continue
		}
		strong = append(strong[:i], strong[next:]...)
		weak = append(weak[:i], weak[next:]...)
		i--

	}

	weak = append(weak[1:], weak[0])

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

	// INSERT REGULAR PAIRS

	// Scan first id
	var firstUserIdBuffer string

	rows.Next()
	err := rows.Scan(&firstUserIdBuffer)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// set first id
	firstUserId := uuid.Parse(firstUserIdBuffer)

	// set lastUserId to firstUser Id
	lastUserId := firstUserId

	var targets []*targetPair

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
	}

	err = rows.Close()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// have last user target the first
	pair := &targetPair{lastUserId, nil, "", firstUserId, nil, ""}
	// Add pair to targets slice
	targets = append(targets, pair)

	// INSERT STRONG WEAK PAIRS

	// Get number of players
	numPlayers, appErr := game.GetNumActivePlayers()
	if appErr != nil {
		return appErr
	}

	// Determine how often to insert a strong/weak pair
	insertNum := (numPlayers / numTeams) - 1
	modCounter := 0
	for len(strong) != 0 {
		for index, currentPair := range targets {
			if (modCounter % insertNum) == 0 {

				if (currentPair.AssassinUserRole != ``) || (currentPair.TargetUserRole != ``) {
					continue
				}
				if len(strong) == 0 {
					break
				}

				var nextStrong, nextWeak uuid.UUID
				// Get next strong/weak pair
				nextStrong, strong = strong[0], strong[1:]
				nextWeak, weak = weak[0], weak[1:]

				oldAssassin := currentPair.AssassinId
				oldTarget := currentPair.TargetId

				// Replace current pair with the old assassin targeting the strong player
				targets[index] = &targetPair{oldAssassin, nil, "", nextStrong, nil, "strong"}

				// Create strong/weak pair
				strongWeakPair := &targetPair{nextStrong, nil, "strong", nextWeak, nil, "weak"}
				weakTargetPair := &targetPair{nextWeak, nil, "weak", oldTarget, nil, ""}
				targets = append(targets, strongWeakPair, weakTargetPair)

			}
			modCounter++
		}
	}

	return game.insertTargetsWithDelete(tx, targets)
}

func (game *Game) getAllActivePlayersAsRows() (rows *sql.Rows, appErr *ApplicationError) {
	// Get new target list
	rows, err := db.Query(`SELECT user_id, team_id, user_role FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY random()`, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return rows, nil
}

// Assign targets and space them out by team
func (game *Game) assignTargetsByTeams(tx *sql.Tx) (appErr *ApplicationError) {
	// Get active players
	rows, appErr := game.getAllActivePlayersAsRows()
	if appErr != nil {
		return appErr
	}

	// Get the list of team ids
	teamsList, appErr := game.GetTeamsWithRegularPlayersLeft()
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
	err := rows.Close()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
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

	// If we have more than 5 teams we may need to save some for later
	saveSomeForLater := len(userList) > 5
	if saveSomeForLater {
		var largestTeamSize, smallestTeamSize int
		for team := range userList {
			currentTeamSize := len(team)
			if largestTeamSize < currentTeamSize {
				largestTeamSize = currentTeamSize
			}
			if smallestTeamSize > currentTeamSize {
				smallestTeamSize = currentTeamSize
			}
		}
		// If the largest team size is sufficiently larger than the smallest team size we need to save some targets for later
		saveSomeForLater = largestTeamSize >= (smallestTeamSize * 3 / 2)
	}

	// While we still have users to assign
	for assigned < numUsers {
		if numTeams <= 3 && saveSomeForLater {
			break
		}

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
		i += captainSpace
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

	// loop through all of our teams
	for teamIdString, team := range userList {
		// user the same counters for all sers to insert them
		i := 0
		j := 0
		teamId := uuid.Parse(teamIdString)

		// loop through each user on each team
		for _, userId := range team {
			foundPair := false

			// keep looping until we have a proper pair
			for !foundPair {
				//
				if i >= numUsers {
					i = 0
					j++
				}

				if j > triesForTeam {
					foundPair = true
					continue
				}

				// check the assassin/target teams to make sure they don't conflict
				assassinTeamId, targetTeamId := targetList[i].AssassinTeamId, targetList[i].TargetTeamId

				if uuid.Equal(teamId, assassinTeamId) {
					i++
					continue
				}
				if uuid.Equal(teamId, targetTeamId) {
					i++
					continue
				}
				foundPair = true
			}
			// Get all the assassin/target information
			assassinId, assassinTeamId, assassinUserRole, targetId, targetTeamId, targetUserRole := targetList[i].AssassinId, targetList[i].AssassinTeamId, targetList[i].AssassinUserRole, targetList[i].TargetId, targetList[i].TargetTeamId, targetList[i].TargetUserRole
			// Insert the captain after the assassin
			targetList[i] = &targetPair{assassinId, assassinTeamId, assassinUserRole, userId, teamId, `dm_user`}
			// Add a pair for the target
			captainPair := &targetPair{userId, teamId, `dm_user`, targetId, targetTeamId, targetUserRole}
			targetList = append(targetList, captainPair)
			i += 3
		}
	}

	return game.insertTargetsWithDelete(tx, targetList)

}

// inserts targets into the db, requires transaction to be wrapped in
func (game *Game) insertTargets(tx *sql.Tx, targetList []*targetPair) (appErr *ApplicationError) {

	// loop through all targets
	for _, pair := range targetList {
		// Prepare the statement to insert the target row
		insertTarget, err := tx.Prepare(`INSERT INTO dm_user_targets (user_id, target_id, game_id) VALUES ($1, $2, $3)`)
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
	return nil
}

// inserts a slice of targetPairs into the database and deletes the current pairs in the db
func (game *Game) insertTargetsWithDelete(tx *sql.Tx, targetList []*targetPair) (appErr *ApplicationError) {

	// Delete targets
	appErr = game.DeleteTargetsTransactional(tx)
	if appErr != nil {
		return appErr
	}

	// if we have no targets just clear the db
	if len(targetList) == 0 {
		return
	}

	appErr = game.insertTargets(tx, targetList)
	if appErr != nil {
		return appErr
	}
	return nil
}

// Gets all active players for a game as a slice of uuids
func (game *Game) GetAllActivePlayersAsUUIDSlice() (users []uuid.UUID, appErr *ApplicationError) {

	rows, err := db.Query(`SELECT user_id FROM dm_user_game_mapping WHERE game_id = $1 AND alive = true AND (user_role = 'dm_user' OR user_role = 'dm_captain') ORDER BY random()`, game.GameId.String())
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
	// Close the rows
	err = rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return users, nil
}

// Deletes all targets for a game
func (game *Game) DeleteTargetsTransactional(tx *sql.Tx) (appErr *ApplicationError) {

	// Prepare statement to delete targets
	deleteTargets, err := tx.Prepare(`DELETE FROM dm_user_targets WHERE game_id = $1`)
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Execute statement to delete targets
	_, err = tx.Stmt(deleteTargets).Exec(game.GameId.String())
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return nil
}

// Deletes all targets for a game
func (game *Game) DeleteTargets() (appErr *ApplicationError) {

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Run internal delete code
	appErr = game.DeleteTargetsTransactional(tx)
	if appErr != nil {
		tx.Rollback()
		return appErr
	}

	// Check commit for error
	err = tx.Commit()
	if err != nil {
		return NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}
	return nil
}

// Assign all targets plainly
func (game *Game) assignTargets(tx *sql.Tx, users []uuid.UUID, skipDelete bool) (appErr *ApplicationError) {

	if len(users) == 0 {
		if !skipDelete {
			return game.DeleteTargets()
		}
		return
	}

	var targets []*targetPair
	// Loop through rows
	firstUserId := users[0]
	prevUserId := firstUserId
	users = users[1:]
	for _, userId := range users {
		// Create a new target pair
		pair := &targetPair{prevUserId, nil, "", userId, nil, ""}

		// Append targetPair to targets
		targets = append(targets, pair)

		// Increment to the next user
		prevUserId = userId
	}

	// Set the last user to target the first
	pair := &targetPair{prevUserId, nil, "", firstUserId, nil, ""}
	targets = append(targets, pair)

	// Execute the actual insert code
	if skipDelete {
		return game.insertTargets(tx, targets)
	}
	return game.insertTargetsWithDelete(tx, targets)

}

type SuperTargetPair struct {
	AssassinId       uuid.UUID `json:"assassin_id"`
	AssassinName     string    `json:"assassin_name"`
	AssassinTeam     string    `json:"assassin_team_name"`
	AssassinUserRole string    `json:"assassin_user_role"`
	AssassinKills    int       `json:"assassin_kills"`
	TargetId         uuid.UUID `json:"target_id"`
	TargetName       string    `json:"target_name"`
	TargetTeam       string    `json:"target_team_name"`
	TargetUserRole   string    `json:"target_user_role"`
	TargetKills      int       `json:"target_kills"`
}

// Get all targets for a game for the super admin panel
func (game *Game) GetTargets() (targets map[string]*SuperTargetPair, appErr *ApplicationError) {

	sql := `select t.user_id, p1.value, p2.value, team1.team_name, m.user_role, m.kills, t.target_id, p3.value, p4.value, team2.team_name, g.user_role, g.kills FROM dm_user_targets AS t, dm_user_properties AS p1, dm_user_properties AS p2, dm_user_properties AS p3, dm_user_properties AS p4, dm_teams as team1, dm_teams as team2, dm_user_game_mapping AS m, dm_user_game_mapping AS g WHERE t.user_id = p1.user_id AND p1.key='first_name' AND t.user_id = p2.user_id AND p2.key='last_name' AND team1.team_id = m.team_id AND team2.team_id = g.team_id AND m.user_id = t.user_id AND g.user_id = t.target_id AND p3.user_id = t.target_id AND p3.key='first_name' AND p4.user_id = t.target_id AND p4.key='last_name' AND t.game_id = $1`

	// If teams aren't enabled use the no team version of the query
	teamsEnabled, appErr := game.GetGameProperty(`teams_enabled`)
	if (appErr != nil) || (teamsEnabled != `true`) {
		sql = `select t.user_id, p1.value, p2.value, 'No Team', m.user_role, m.kills, t.target_id, p3.value, p4.value, 'No Team', g.user_role, g.kills FROM dm_user_targets AS t, dm_user_properties AS p1, dm_user_properties AS p2, dm_user_properties AS p3, dm_user_properties AS p4, dm_user_game_mapping AS m, dm_user_game_mapping AS g WHERE t.user_id = p1.user_id AND p1.key='first_name' AND t.user_id = p2.user_id AND p2.key='last_name' AND m.user_id = t.user_id AND g.user_id = t.target_id AND p3.user_id = t.target_id AND p3.key='first_name' AND p4.user_id = t.target_id AND p4.key='last_name' AND t.game_id = $1`
	}

	rows, err := db.Query(sql, game.GameId.String())
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	targets = make(map[string]*SuperTargetPair)

	for rows.Next() {
		var assassinIdBuffer, assassinFirstName, assassinLastName, assassinTeam, assassinUserRole, targetIdBuffer, targetFirstName, targetLastName, targetTeam, targetUserRole string
		var assassinKills, targetKills int

		err := rows.Scan(&assassinIdBuffer, &assassinFirstName, &assassinLastName, &assassinTeam, &assassinUserRole, &assassinKills, &targetIdBuffer, &targetFirstName, &targetLastName, &targetTeam, &targetUserRole, &targetKills)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		assassinId := uuid.Parse(assassinIdBuffer)
		assassinName := assassinFirstName + ` ` + assassinLastName
		targetId := uuid.Parse(targetIdBuffer)
		targetName := targetFirstName + ` ` + targetLastName

		newPair := &SuperTargetPair{assassinId, assassinName, assassinTeam, assassinUserRole, assassinKills, targetId, targetName, targetTeam, targetUserRole, targetKills}
		targets[assassinIdBuffer] = newPair
	}
	return targets, nil
}
