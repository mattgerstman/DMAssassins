package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Issue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

func postGithubIssue(title, message, email, name string) (appErr *ApplicationError) {

	var body string
	body = "Name: " + name + "\n"
	body += "Email: " + email + "\n"
	body += "Message: " + message + "\n"

	issue := &Issue{title, body, []string{"Support"}}
	url := Config.GithubRepo + "issues?access_token=" + Config.GithubApiKey

	data, err := json.Marshal(issue)
	if err != nil {
		return NewApplicationError(`Error creating issue`, err, ErrCodeInternalServerWTF)
	}

	buf := bytes.NewBuffer(data)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		return NewApplicationError(`Error creating issue`, err, ErrCodeExternalService)
	}

	_ = resp

	return nil
}
