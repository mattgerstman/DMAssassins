package main

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"github.com/mailgun/mailgun-go"
	"text/template"
	"fmt"
)

type EmailData struct {
	GameName  string
	APIDomain string
}

func parseEmailableUsers(users map[string]*User, onlyAlive bool) (userList []*User) {
	for _, user := range users {
		if onlyAlive {
			if alive, ok := user.Properties["alive"]; ok {
				if alive != "true" {
					continue
				}
			}
		}

		if allowEmail, ok := user.Properties["allow_email"]; ok {
			if allowEmail == "true" {
				userList = append(userList, user)
			}
		}
	}
	return userList
}

func (game *Game) getEmailableUsersForGame(onlyAlive bool) (userList []*User, appErr *ApplicationError) {
	userMap, appErr := game.GetAllUsersForGame()
	if appErr != nil {
		return nil, appErr
	}

	userList = parseEmailableUsers(userMap, onlyAlive)
	return userList, nil
}

func (game *Game) sendStartGameEmail() (id string, appErr *ApplicationError) {
	// users, appErr := game.getEmailableUsersForGame(false)
	// if appErr != nil {
	// 	return appErr
	// }

	// Temporary code to only email me. I don't want to waste emails with MailGun on testing
	user, appErr := GetUserById(uuid.Parse("5759a74a-2f1b-11e4-9241-685b35b45205"))
	if appErr != nil {
		return "", appErr
	}
	users := []*User{user}

	var bodyBuffer bytes.Buffer
	emailData := &EmailData{game.GameName, Config.APIDomain}
	t, err := template.ParseFiles("templates/game-started.txt")
	if err != nil {
		// TODO be less lazy
		fmt.Println(err)
	}
	t.Execute(&bodyBuffer, emailData)

	subject := game.GameName + ` DMAssassins Has Begun!`
	tag := `StartGame`
	body := bodyBuffer.String()

	return sendEmail(subject, body, tag, users)

}

func sendEmail(subject, body, tag string, users []*User) (id string, appErr *ApplicationError) {

	mg := mailgun.NewMailgun(Config.MailGunDomain, Config.MailGunPrivateKey, Config.MailGunPublicKey)

	m := mg.NewMessage(
		Config.MailGunSender,
		subject,
		body,
	)
	m.AddTag(tag)
	m.SetTracking(true)

	for _, user := range users {
		err := m.AddRecipientAndVariables(user.Email, map[string]interface{}{
			"first_name": user.Properties[`first_name`],
			"user_id":    user.UserId.String(),
		})
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeEmail)
		}
	}
	_, id, err := mg.Send(m)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeEmail)
	}

	return id, nil

}
