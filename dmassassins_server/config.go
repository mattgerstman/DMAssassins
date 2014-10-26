package main

import (
	"encoding/json"
	"log"
	"os"
)

type Configuration struct {
	APIDomain         string
	FBAppId           string
	FBAppSecret       string
	FBAccessToken     string
	DatabaseURL       string
	SentryDSN         string
	SupportEmail      string
	MailGunPublicKey  string
	MailGunPrivateKey string
	MailGunDomain     string
	MailGunSender     string
	MailGunReplyTo    string
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
