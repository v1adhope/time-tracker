package entities

import "errors"

var (
	ErrorUserHasAlreadyExist = errors.New("user has already exist")
	ErrorUserDoesNotExist    = errors.New("user doesn't exist")

	ErrorTaskDoesNotExist      = errors.New("task doesn't exist")
	ErrorNoAnyTasksForThisUser = errors.New("no any tasks for this user")
)
