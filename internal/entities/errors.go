package entities

import "errors"

var (
	ErrorUserHasAlreadyExist = errors.New("User has already exist")
	ErrorUserDoesNotExist    = errors.New("User doesn't exist")

	ErrorTaskDoesNotExist = errors.New("Task doesn't exist")
)
