package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	APIDomain         string `json:"api_domain"`
	DefaultTimeZone   string `json:"default_timezone"`
	FBAppId           string `json:"fb_app_id"`
	FBAppSecret       string `json:"fb_app_secret"`
	FBAccessToken     string `json:"fb_access_token"`
	GithubRepo        string `json:"github_repo"`
	GithubApiKey      string `json:"github_api_key"`
	DatabaseURL       string `json:"database_url"`
	SentryDSN         string `json:"sentry_dsn"`
	SupportEmail      string `json:"support_email"`
	MailGunPublicKey  string `json:"mailgun_public_key"`
	MailGunPrivateKey string `json:"mailgun_private_key"`
	MailGunDomain     string `json:"mailgun_domain"`
	MailGunSender     string `json:"mailgun_sender"`
	MailGunReplyTo    string `json:"mailgun_replyto"`
}

var Config *Configuration

// Loads config variables from file into global Config struct
func LoadConfig() (appErr *ApplicationError) {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to load config with message:", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("Failed to load config with message:", err)
	}
	return LoadPlotTwists()
}

type PlotTwistMap map[string]*PlotTwist

var PlotTwistConfig PlotTwistMap

// Loads plot twists from config file
func LoadPlotTwists() (appErr *ApplicationError) {
	file, err := os.Open("plot_twists.json")
	if err != nil {
		log.Fatal("Failed to load plot twist config with message:", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&PlotTwistConfig)
	if err != nil {
		log.Fatal("Failed to load plot twist config with message:", err)
	}
	return nil
}
