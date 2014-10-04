package main

import (
	"bytes"
	"github.com/mailgun/mailgun-go"
	"text/template"
)

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

func (user *User) SendBanhammerEmail(GameName string) (id string, appErr *ApplicationError) {

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  GameName,
		"APIDomain": Config.APIDomain,
	}

	t, err := template.ParseFiles("templates/banhammer.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/banhammer.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	subject := `Welcome to DMAssassins!`
	tag := `WelcomeUser`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	users := []*User{user}

	return sendEmail(subject, body, htmlBody, tag, users)

}

func (user *User) SendAdminWelcomeEmail() (id string, appErr *ApplicationError) {

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"APIDomain": Config.APIDomain,
	}

	t, err := template.ParseFiles("templates/welcome-admin.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/welcome-admin.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	subject := `Welcome to DMAssassins!`
	tag := `WelcomeUser`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	users := []*User{user}

	return sendEmail(subject, body, htmlBody, tag, users)

}

func (user *User) SendUserWelcomeEmail() (id string, appErr *ApplicationError) {

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"APIDomain": Config.APIDomain,
	}

	t, err := template.ParseFiles("templates/welcome-user.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/welcome-user.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	subject := `Welcome to DMAssassins!`
	tag := `WelcomeUser`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	users := []*User{user}

	return sendEmail(subject, body, htmlBody, tag, users)

}

func (game *Game) sendStartGameEmail() (id string, appErr *ApplicationError) {
	users, appErr := game.getEmailableUsersForGame(false)
	if appErr != nil {
		return "", appErr
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  game.GameName,
		"APIDomain": Config.APIDomain,
	}
	t, err := template.ParseFiles("templates/game-started.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/game-started.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	subject := game.GameName + ` DMAssassins Has Begun!`
	tag := `StartGame`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	return sendEmail(subject, body, htmlBody, tag, users)

}

func sendEmail(subject, body, htmlBody, tag string, users []*User) (id string, appErr *ApplicationError) {

	mg := mailgun.NewMailgun(Config.MailGunDomain, Config.MailGunPrivateKey, Config.MailGunPublicKey)

	m := mg.NewMessage(
		Config.MailGunSender,
		subject,
		body,
	)
	m.AddTag(tag)
	m.SetTracking(true)
	m.SetHtml(htmlBody)

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
