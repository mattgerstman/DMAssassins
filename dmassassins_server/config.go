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
func LoadConfig() *ApplicationError {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to load config with message:", err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("Failed to load config with message:", err)
	}
	return nil

}
