package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
	"strings"
)

func loadDefaultRules(adminEmail string) (rules string, appErr *ApplicationError) {

	fileByte, _ := ioutil.ReadFile("rules.mdown")
	output := blackfriday.MarkdownBasic(fileByte)
	//%ADMINEMAIL%
	outputString := string(output)
	outputString = strings.Replace(outputString, `%ADMINEMAIL%`, adminEmail, -1)
	return outputString, nil

}

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

func (game *Game) SetRules(rules string) (appErr *ApplicationError) {

	appErr = game.SetGameProperty("rules", rules)
	if appErr != nil {
		return appErr
	}
	return nil

}
