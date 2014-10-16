package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

// Creates a new secret from the json file
func NewSecret(length int) (secret string, appErr *ApplicationError) {
	file, err := os.Open("secrets.json")
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeFile)
	}

	// Read in words list
	var words []string
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&words)
	if err != nil {
		return "", NewApplicationError("Internal Error", err, ErrCodeFile)
	}

	// Grab 3 words and string them together
	// I'm using a permutation of the array because it's the simplest way i found to get three unique words from it
	rand.Seed(time.Now().UTC().UnixNano())
	for i, j := range rand.Perm(len(words)) {
		if i >= length {
			break
		}
		secret += words[j]
	}
	return secret, nil
}
