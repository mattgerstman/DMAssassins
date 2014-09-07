package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

// Creates a new secret from the json file
func NewSecret() (secret string, appErr *ApplicationError) {
	file, err := os.Open("secrets.json")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeFile)
	}

	var words []string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&words)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeFile)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	for i, j := range rand.Perm(len(words)) {
		if i >= 3 {
			break
		}
		secret += words[j]
	}
	return secret, nil
}
