package main

import (
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/getsentry/raven-go"
	"strconv"
)

func startGameInner() (appErr *ApplicationError) {
	gameId := uuid.Parse("9202fcd2-ccbd-42d4-8c54-99968a38e5e6")
	game, appErr := GetGameById(gameId)
	if appErr != nil {
		return appErr
	}
	appErr = game.AssignTargetsBy(`normal`)
	if appErr != nil {
		return appErr
	}
	if appErr != nil {
		return appErr
	}
	return nil
}

func startGame() {
	appErr := startGameInner()
	if appErr != nil {
		LogWithSentry(appErr, nil, raven.ERROR, nil)
		fmt.Println(appErr.Msg)
		fmt.Println(appErr.Err)
		fmt.Println(appErr.Code)
	}
}

func generateTestUsers() {
	gameId := uuid.Parse("9202fcd2-ccbd-42d4-8c54-99968a38e5e6")
	game, appErr := GetGameById(gameId)
	teams, appErr := game.GetTeams()

	var teamsList []uuid.UUID
	for _, team := range teams {
		teamsList = append(teamsList, team.TeamId)
	}

	numTeams := len(teamsList)

	for i := 0; i < 300; i++ {
		firstName, _ := NewSecret(1)
		lastName, _ := NewSecret(1)
		username := firstName + lastName

		fmt.Println(username)

		email := firstName + `.` + lastName + `@gmail.com`
		facebookId := strconv.Itoa(i + 1000000)

		properties := make(map[string]string)
		properties["facebook"] = "https://facebook.com/" + facebookId

		picture := "https://graph.facebook.com/747160532/picture"
		properties["photo"] = picture + "?width=1000"
		properties["photo_thumb"] = picture + "?width=300&height=300"

		properties["first_name"] = firstName
		properties["last_name"] = lastName
		properties["allow_email"] = "true"

		user, appErr := NewUser(username, email, facebookId, properties)

		gameMapping, appErr := user.JoinGame(gameId, ``)

		teamIndex := i % numTeams
		teamId := teamsList[teamIndex]

		fmt.Println(teamId)
		gameMapping, appErr = user.JoinTeam(teamId)

		_ = appErr
		_ = gameMapping
	}
	_ = appErr
}
