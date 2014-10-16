package main

import (
	"strconv"
)

// Gets a string of $1, $2, $3, ... etc
func GetParamsForSlice(startingNum int, slice []interface{}) (params string) {
	// Construct params to get non strong/weak users
	params = `$` + strconv.Itoa(startingNum)
	numParams := len(slice) + startingNum - 1
	for i := (startingNum + 1); i < numParams; i++ {
		params += `, $` + strconv.Itoa(i)
	}
	return params
}
