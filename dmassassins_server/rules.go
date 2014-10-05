package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"strings"
)

// Pulls the default rules from rules.mdown and
func loadDefaultRules(adminEmail string) (rules string, appErr *ApplicationError) {

	fileByte, _ := ioutil.ReadFile("templates/rules.mdown")
	outputString := string(fileByte)
	outputString = strings.Replace(outputString, `%ADMINEMAIL%`, adminEmail, -1)
	return outputString, nil

}

// Get rules for a game and load it into it's game properties
func (game *Game) GetRules() (rules string, appErr *ApplicationError) {

	rules, appErr = game.GetGameProperty("rules")
	if appErr != nil {
		return "", appErr
	}

	if rules == "" {
		admin, appErr := game.GetAdmin()
		if appErr != nil {
			return "", appErr
		}

		defaultRules, appErr := loadDefaultRules(admin.Email)
		if appErr != nil {
			return "", appErr
		}
		appErr = game.SetGameProperty("rules", defaultRules)
		if appErr != nil {
			return "", appErr
		}
		rules = defaultRules
	}

	return rules, nil

}

func (game *Game) GetHTMLRules() (rules string, appErr *ApplicationError) {

	var ok bool
	if rules, ok = game.Properties["rules"]; !ok {
		rules, appErr = game.GetRules()
		if appErr != nil {
			return "", appErr
		}
	}

	rulesByte := blackfriday.MarkdownBasic([]byte(rules))
	rules = string(rulesByte)
	game.Properties["rules"] = rules
	return rules, nil
}

// Set rules for a game and load it into it's game properties
func (game *Game) SetRules(rules string) (appErr *ApplicationError) {

	appErr = game.SetGameProperty("rules", rules)
	if appErr != nil {
		return appErr
	}
	return nil

}
