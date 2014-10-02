package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/mailgun/mailgun-go"
	"fmt"
)

func parseEmailableUsers(users map[string]*User) (userList []*User) {
	for _, user := range users {
		if allowEmail, ok := user.Properties["allow_email"]; ok {
			if allowEmail == "true" {
				userList = append(userList, user)
			}
		}
	}
	return userList
}

func (game *Game) getEmailableUsersForGame() (userList []*User, appErr *ApplicationError) {
	userMap, appErr := game.GetAllUsersForGame()
	if appErr != nil {
		return nil, appErr
	}

	userList = parseEmailableUsers(userMap)
	return userList, nil
}

func (game *Game) sendStartGameEmail() (id string, appErr *ApplicationError) {
	// users, appErr := game.getEmailableUsersForGame()
	// if appErr != nil {
	// 	return appErr
	// }

	// Temporrary code to only email me. I don't want to waste emails with MailGun on testing
	user, appErr := GetUserById(uuid.Parse("5759a74a-2f1b-11e4-9241-685b35b45205"))
	if appErr != nil {
		return "", appErr
	}
	users := []*User{user}

	mg := mailgun.NewMailgun(Config.MailGunDomain, Config.MailGunPrivateKey, Config.MailGunPublicKey)

	subject := game.GameName + `: Assassins Has Started!`

	body := "Hey %recipient.first_name%!"
	body += "The DMAssassins game \"" + game.GameName + "\" has started. "
	body += "Login to http://dmassassins.com to see your first target.\n\n"
	body += "Sincerely,\n"
	body += "The DMAssassins Team\n\n\n"
	body += "To Unsubscribe go to " + Config.APIDomain + "unsubscribe/%recipient.user_id%\""

	htmlBody := "Hey %recipient.first_name%!<br />"
	htmlBody += "<p>The DMAssassins game \"" + game.GameName + "\" has started. "
	htmlBody += "Login to <a href=\"http://dmassassins.com\">http://dmassassins.com</a> to see your first target.</p>"
	htmlBody += "Sincerely,<br />"
	htmlBody += "The DMAssassins Team<br /><br /><br />"
	htmlBody += "To Unsubscribe click <a href=\"" + Config.APIDomain + "unsubscribe/%recipient.user_id%\">here</a>"

	m := mg.NewMessage(
		Config.MailGunSender,
		subject,
		body,
	)
	m.SetHtml(htmlBody)
	m.AddTag("StartGame")
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
	fmt.Println(id)
	fmt.Println(err)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeEmail)
	}


	return id, nil

}
