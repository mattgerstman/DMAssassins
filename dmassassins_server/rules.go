package main

import (
	"github.com/russross/blackfriday"
	"io/ioutil"
)

func loadDefaultRules() (rules string, appErr *ApplicationError) {

	fileByte, _ := ioutil.ReadFile("rules.mdown")
	output := blackfriday.MarkdownBasic(fileByte)
	return string(output), nil

}

func (game *Game) GetRules() (rules string, appErr *ApplicationError) {

	rules, appErr = game.GetGameProperty("rules")
	if appErr != nil {
		return "", appErr
	}

	if rules == "" {
		defaultRules, appErr := loadDefaultRules()
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
