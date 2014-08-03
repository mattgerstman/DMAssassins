package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Configuration struct {
	FBAppId           string
	FBAppSecret       string
	FBUserAccessToken string
	DatabaseURL       string
	SentryDSN         string
}

var Config *Configuration

func LoadConfig() *ApplicationError {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("error:", err)
	}
	return nil

}
