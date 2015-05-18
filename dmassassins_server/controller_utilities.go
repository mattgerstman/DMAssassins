package main

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"net/http"
)

type Params struct {
	data map[string]interface{}
}

func NewParams(r *http.Request) (params *Params, appErr *ApplicationError) {

	data := make(map[string]interface{})

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, NewApplicationError("Invalid JSON", err, ErrCodeInvalidJSON)
	}
	return &Params{data}, nil
}

func (params *Params) GetParam(key string) (param interface{}, appErr *ApplicationError) {
	data := params.data
	if _, ok := data[key]; !ok {
		msg := "Missing Parameter: " + key
		err := errors.New(msg)
		return nil, NewApplicationError(msg, err, ErrCodeMissingParameter)
	}

	return data[key], nil
}

func getInvalidParameterAppErr(key string) (appErr *ApplicationError) {
	msg := "Invalid Parameter: " + key
	err := errors.New(msg)
	return NewApplicationError(msg, err, ErrCodeInvalidParameter)
}

func (params *Params) GetIntParam(key string) (intParam int, appErr *ApplicationError) {
	param, appErr := params.GetParam(key)
	if appErr != nil {
		return 0, appErr
	}

	if param == nil {
		return 0, nil
	}

	if intParam, ok := param.(int); ok {
		return intParam, nil
	}

	return 0, getInvalidParameterAppErr(key)
}

func (params *Params) GetStringParam(key string) (stringParam string, appErr *ApplicationError) {
	param, appErr := params.GetParam(key)
	if appErr != nil {
		return "", appErr
	}

	if param == nil {
		return "", nil
	}

	if stringParam, ok := param.(string); ok {
		return stringParam, nil
	}

	return "", getInvalidParameterAppErr(key)
}

func (params *Params) GetUUIDParam(key string) (uuidParam uuid.UUID, appErr *ApplicationError) {
	stringParam, appErr := params.GetStringParam(key)
	if appErr != nil {
		return nil, appErr
	}
	return uuid.Parse(stringParam), nil
}
