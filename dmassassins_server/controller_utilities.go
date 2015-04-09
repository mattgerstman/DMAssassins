package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func GetParams(r *http.Request) (params map[string]interface{}, appErr *ApplicationError) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}
	return params, nil
}

func GetParam(r *http.Request, required bool, key string) (param interface{}, appErr *ApplicationError) {
	params, appErr := GetParams(r)
	if appErr != nil {
		return "", appErr
	}

	if _, ok := params[key]; !ok && required {
		msg := "Missing Parameter: " + key
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	if _, ok := params[key]; !ok {
		return nil, nil
	}
	return params[key], nil
}

func GetStringParam(r *http.Request, required bool, key string) (stringParam string, appErr *ApplicationError) {

	param, appErr := GetParam(r, required, key)
	if appErr != nil {
		return "", appErr
	}

	if param == nil {
		return "", nil
	}

	return param.(string), nil
}
