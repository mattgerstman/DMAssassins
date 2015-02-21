package main

import (
	"github.com/getsentry/raven-go"
	"github.com/gorilla/context"
	"net/http"
)

func SetUserForRequest(r *http.Request, user *User) {
	context.Set(r, "user", user)
}

func GetUserForRequest(r *http.Request) (user *User) {
	if rv := context.Get(r, "user"); rv != nil {
		return rv.(*User)
	}
	return nil
}

func GetSentryUserForRequest(r *http.Request) (sentryUser *raven.User) {
	user := GetUserForRequest(r)
	if user == nil {
		return nil
	}
	return NewSentryUser(user)
}
