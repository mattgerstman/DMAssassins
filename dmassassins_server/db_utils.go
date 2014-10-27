package main

import (
	"code.google.com/p/go-uuid/uuid"
	"database/sql"
	//	"fmt"
	"strconv"
)

// Gets a string of $1, $2, $3, ... etc
func GetParamsForSlice(startingNum int, slice []interface{}) (params string) {
	// Construct params to get non strong/weak users
	params = `$` + strconv.Itoa(startingNum)
	numParams := len(slice) + startingNum
	for i := (startingNum + 1); i < numParams; i++ {
		params += `, $` + strconv.Itoa(i)
	}
	return params
}

// Converts a set of rows with just a userId to a slice of uuids
func ConvertUserIdRowsToSlice(rows *sql.Rows) (users []uuid.UUID, appErr *ApplicationError) {
	var userIdBuffer string
	for rows.Next() {
		// load the uuid into a buffer string
		err := rows.Scan(&userIdBuffer)
		if err != nil {
			return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
		}
		// Parse the uuid
		userId := uuid.Parse(userIdBuffer)
		// append the uuid
		users = append(users, userId)
	}
	err := rows.Close()
	if err != nil {
		return nil, NewApplicationError("Internal Error", err, ErrCodeDatabase)
	}

	// Return the users list
	return users, nil
}

// Converst a slice of uuids to a slice of interfaces
func ConvertUUIDSliceToInterface(uuidSlice []uuid.UUID) (output []interface{}) {
	for _, uuid := range uuidSlice {
		output = append(output, uuid.String())
	}
	return output
}
