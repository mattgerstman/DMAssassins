package main

import (
	"bytes"
	"github.com/mailgun/mailgun-go"
	"text/template"
)

// Remove users from the list that shouldn't receive emails
// Users who opt out will always be removed
// Sometimes we need to determine if they're alive or not
func parseEmailableUsers(users map[string]*User, onlyAlive bool) (userList []*User) {
	for _, user := range users {
		// If we only want to email living users, skip the dead oens
		if onlyAlive {
			if alive, ok := user.Properties["alive"]; ok {
				if alive != "true" {
					continue
				}
			}
		}

		// Add users who have opted in to the userList
		if allowEmail, ok := user.Properties["allow_email"]; ok {
			if allowEmail == "true" {
				userList = append(userList, user)
			}
		}
	}
	return userList
}

// Gets all users for a game we can email
func (game *Game) getEmailableUsersForGame(onlyAlive bool) (userList []*User, appErr *ApplicationError) {
	// Get all of the game's users
	userMap, appErr := game.GetAllUsersForGame()
	if appErr != nil {
		return nil, appErr
	}

	// Parse out users that opted out of emails
	userList = parseEmailableUsers(userMap, onlyAlive)
	return userList, nil
}

// Gets a list of emailable users for a game
func (game *Game) GetEmailsForGame(onlyAlive bool) (emails []string, appErr *ApplicationError) {

	userList, appErr := game.getEmailableUsersForGame(onlyAlive)
	if appErr != nil {
		return nil, appErr
	}

	for _, user := range userList {
		emails = append(emails, user.Email)
	}

	return emails, nil
}

// Inform a user they've died
func (user *User) SendDeadEmail(GameName string) (id string, appErr *ApplicationError) {
	// Make sure we're allowed to email the user
	allowEmail, appErr := user.GetUserProperty("allow_email")
	if appErr != nil {
		return "", nil
	}
	if allowEmail != "true" {
		return "", nil
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/you-died.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/you-died.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `You Have Been Killed In ` + GameName + ` DMAssassins`
	tag := `Killed`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Inform a user we've revived them
func (user *User) SendReviveEmail(GameName string) (id string, appErr *ApplicationError) {
	// Make sure we're allowed to email the user
	allowEmail, appErr := user.GetUserProperty("allow_email")
	if appErr != nil {
		return "", nil
	}
	if allowEmail != "true" {
		return "", nil
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/revive.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/revive.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `You Have Been Revived In ` + GameName + ` DMAssassins`
	tag := `Revived`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Inform a user they have a new target
func (user *User) SendNewTargetEmail(GameName string) (id string, appErr *ApplicationError) {
	// Make sure we're allowed to email the user
	allowEmail, appErr := user.GetUserProperty("allow_email")
	if appErr != nil {
		return "", nil
	}
	if allowEmail != "true" {
		return "", nil
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/new-target.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Set up the subject and contents of the email
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/new-target.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `You Have A New Target In ` + GameName + ` DMAssassins`
	tag := `NewTarget`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Inform a user they've been banned
func (user *User) SendBanhammerEmail(GameName string) (id string, appErr *ApplicationError) {
	// Make sure we're allowed to email the user
	allowEmail, appErr := user.GetUserProperty("allow_email")
	if appErr != nil {
		return "", nil
	}
	if allowEmail != "true" {
		return "", nil
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/banhammer.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/banhammer.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `You Have Been Banned From ` + GameName + ` DMAssassins`
	tag := `WelcomeUser`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Welcome an Admin to the Game
func (user *User) SendAdminWelcomeEmail() (id string, appErr *ApplicationError) {
	// Make sure we're allowed to email the user
	allowEmail, appErr := user.GetUserProperty("allow_email")
	if appErr != nil {
		return "", nil
	}
	if allowEmail != "true" {
		return "", nil
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/welcome-admin.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/welcome-admin.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `Welcome to DMAssassins!`
	tag := `WelcomeAdmin`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Welcome a user to the game
func (user *User) SendUserWelcomeEmail() (id string, appErr *ApplicationError) {
	// Make sure we're allowed to email the user
	allowEmail, appErr := user.GetUserProperty("allow_email")
	if appErr != nil {
		return "", nil
	}
	if allowEmail != "true" {
		return "", nil
	}
	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/welcome-user.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/welcome-user.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `Welcome to DMAssassins!`
	tag := `WelcomeUser`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Inform users the game has started
func (game *Game) sendStartGameEmail() (id string, appErr *ApplicationError) {

	// Get all users for the game that we're allowed to email
	users, appErr := game.getEmailableUsersForGame(false)
	if appErr != nil {
		return "", appErr
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  game.GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/game-started.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/game-started.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := game.GameName + ` DMAssassins Has Begun!`
	tag := `StartGame`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Inform users the game has ended
func (game *Game) sendGameOverEmail() (id string, appErr *ApplicationError) {

	// Get all users for the game that we're allowed to email
	users, appErr := game.getEmailableUsersForGame(false)
	if appErr != nil {
		return "", appErr
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  game.GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/game-ended.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/game-ended.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := game.GameName + ` DMAssassins Is Over!`
	tag := `StartGame`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Send an email through mailgun
func sendEmail(subject, body, htmlBody, tag string, users []*User) (id string, appErr *ApplicationError) {
	// Instantiate the mailgun object
	mg := mailgun.NewMailgun(Config.MailGunDomain, Config.MailGunPrivateKey, Config.MailGunPublicKey)

	// Instantiate the message
	m := mg.NewMessage(
		Config.MailGunSender,
		subject,
		body,
	)

	// Add data about the message and set the HTML body
	m.AddTag(tag)
	m.SetTracking(true)
	m.SetHtml(htmlBody)

	// Add all of the users as receipients with their name/user_id as variables
	for _, user := range users {
		err := m.AddRecipientAndVariables(user.Email, map[string]interface{}{
			"first_name": user.Properties[`first_name`],
			"user_id":    user.UserId.String(),
		})
		if err != nil {
			return "", NewApplicationError("Internal Error", err, ErrCodeEmail)
		}
	}

	// Send the email to mailgun
	_, id, err := mg.Send(m)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeEmail)
	}

	return id, nil

}
