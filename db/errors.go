package db

import "errors"

var UserExistsError = errors.New("Username already Exists")
var UserDoesNotExistsError = errors.New("User doesn't exist")
