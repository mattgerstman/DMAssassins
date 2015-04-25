package main

import (
	"bytes"
	"code.google.com/p/go-uuid/uuid"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/mailgun/mailgun-go"
	"text/template"
	"time"
)

// Remove users from the list that shouldn't receive emails
// Users who opt out will always be removed
// Sometimes we need to determine if they're alive or not
func parseEmailableUsers(users map[string]*User, onlyAlive bool) (userList []*User) {
	for _, user := range users {
		// If we only want to email living users, skip the dead ones
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
	subject := `You Have Been Killed In ` + GameName + ` Assassins`
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
	subject := `You Have Been Revived In ` + GameName + ` Assassins`
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
	subject := `You Have A New Target In ` + GameName + ` Assassins`
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
	subject := `You Have Been Banned From ` + GameName + ` Assassins`
	tag := `Banhammer`
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
	subject := `Welcome to Assassins!`
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
	subject := `Welcome to Assassins!`
	tag := `WelcomeUser`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Inform a user their email has changed
func (user *User) SendChangeEmailEmail() (id string, appErr *ApplicationError) {
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
	t, err := template.ParseFiles("templates/change-email.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/change-email.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := `Assassins - Email Notification`
	tag := `ChangeEmail`
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
	subject := game.GameName + ` Assassins Has Begun!`
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
	subject := game.GameName + ` Assassins Is Over!`
	tag := `StartGame`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Gets email data for a plot twist
func getPlotTwistEmailData(twistName string) (emailData map[string]interface{}) {
	emailData = map[string]interface{}{
		"APIDomain": Config.APIDomain,
	}
	twist, appErr := GetPlotTwist(twistName)
	if appErr != nil {
		LogWithSentry(appErr, map[string]string{"plot_twist": twistName}, raven.WARNING, nil)
		return
	}
	if twist.KillTimer == 0 {
		return
	}
	// calculate killTimer deadline
	now := time.Now()
	fmt.Println(twist.KillTimer)
	executeTime := now.Add(time.Duration(twist.KillTimer) * time.Hour)
	deadline := executeTime.Format(`Monday at 3:04 PM MST`)
	emailData[`Deadline`] = deadline
	return

}

// Inform users of a plot twist
func (game *Game) SendPlotTwistEmail(twistName string) (id string, appErr *ApplicationError) {

	// Get all users for the game that we're allowed to email
	users, appErr := game.getEmailableUsersForGame(true)
	if appErr != nil {
		return "", appErr
	}

	var bodyBuffer bytes.Buffer
	emailData := getPlotTwistEmailData(twistName)

	// Compile the plain email template
	t, err := template.ParseFiles("templates/plot_twists/" + twistName + ".txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/plot_twists/" + twistName + ".html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := game.GameName + ` Assassins - Plot Twist!`
	tag := `PlotTwist`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Send an email to a user who has a new target due to the defend the weak kill mode
func (user *User) SendDefendWeakNewTargetEmail(gameName string) (id string, appErr *ApplicationError) {
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
	t, err := template.ParseFiles("templates/plot_twists/defend_weak_new_target.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/plot_twists/defend_weak_new_target.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := gameName + ` Assassins - You have a new target!`
	tag := `DefendWeakNewTarget`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Send an email to a user who has died due to the defend the weak kill mode
func (user *User) SendDefendWeakKilledEmail(gameName string) (id string, appErr *ApplicationError) {
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
		"GameName":  gameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/plot_twists/defend_weak_died.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/plot_twists/defend_weak_died.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := gameName + ` Assassins - Plot Twist: You're Dead`
	tag := `DefendWeakKilled`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Wrap the user in a slice for the sending function
	users := []*User{user}

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Seperate users who were killed in a plot twist from those who are alive/already dead
func splitTheDead(users []*User, killedUsers []uuid.UUID) (aliveUsers, deadUsers []*User) {
	// Convert killedUsers to a map for constant lookup
	killedUsersMap := make(map[string]bool)
	for _, userId := range killedUsers {
		killedUsersMap[userId.String()] = true
	}

	// split alive and killed users throw out already dead users

	for _, user := range users {
		userIdKey := user.UserId.String()
		// If the user was killed append them to the dead list
		if _, ok := killedUsersMap[userIdKey]; ok {
			deadUsers = append(deadUsers, user)
			continue
		}
		// if the user wasn't alive going in ignore them
		if alive, ok := user.Properties["alive"]; ok {
			if alive != "true" {
				continue
			}
		}

		// Append all the living users to the ok list
		aliveUsers = append(aliveUsers, user)

	}
	return aliveUsers, deadUsers
}

// Email users who survived the timer
func (game *Game) sendSurvivedTimerEmail(users []*User) (id string, appErr *ApplicationError) {
	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  game.GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/plot_twists/timer_survived.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/plot_twists/timer_survived.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := game.GameName + ` Assassins - You Survived The Countdown!`
	tag := `SurvivedTimer`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)
}

// Email users who were killed by the timer
func (game *Game) sendKilledByTimerEmail(users []*User) (id string, appErr *ApplicationError) {
	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  game.GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/plot_twists/timer_killed.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/plot_twists/timer_killed.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := game.GameName + ` Assassins - You Missed The Countdown!`
	tag := `KilledByTimer`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)
}

func (game *Game) SendTimerDisabledEmail() (id string, appErr *ApplicationError) {
	// Get all users for the game that we're allowed to email
	users, appErr := game.getEmailableUsersForGame(true)
	if appErr != nil {
		return "", appErr
	}

	var bodyBuffer bytes.Buffer
	emailData := map[string]interface{}{
		"GameName":  game.GameName,
		"APIDomain": Config.APIDomain,
	}

	// Compile the plain email template
	t, err := template.ParseFiles("templates/plot_twists/timer_disabled.txt")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	t.Execute(&bodyBuffer, emailData)

	// Compile the HTML email template
	var htmlBodyBuffer bytes.Buffer
	htmlT, err := template.ParseFiles("templates/plot_twists/timer_disabled.html")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeBadTemplate)
	}
	htmlT.Execute(&htmlBodyBuffer, emailData)

	// Set up the subject and contents of the email
	subject := game.GameName + ` Assassins Cancelled Countdown`
	tag := `TimerDisabled`
	body := bodyBuffer.String()
	htmlBody := htmlBodyBuffer.String()

	// Send the email
	return sendEmail(subject, body, htmlBody, tag, users)

}

// Send an email for an expired timer
func (game *Game) SendTimerExpiredEmail(killedUsers []uuid.UUID) (appErr *ApplicationError) {
	// Get all users for the game
	users, appErr := game.getEmailableUsersForGame(false)
	if appErr != nil {
		return appErr
	}

	// Seperate the killed from the living
	aliveUsers, deadUsers := splitTheDead(users, killedUsers)

	if len(aliveUsers) > 0 {

		// Inform users they survived the countdown
		_, appErr = game.sendSurvivedTimerEmail(aliveUsers)
		if appErr != nil {
			return appErr
		}
	}

	if len(deadUsers) > 0 {
		// Inform users they didn't survive the countdown
		_, appErr = game.sendKilledByTimerEmail(deadUsers)
		if appErr != nil {
			return appErr
		}
	}
	return nil
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
	if htmlBody != "" {
		m.SetHtml(htmlBody)
	}

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
