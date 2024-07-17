package entities

import "errors"

var (
	ErrorUserHasAlreadyExistWithThatPassport          = errors.New("user has already exist that passport data")
	ErrorUsersDoesNotExist                            = errors.New("user(s) doesn't exist")
	ErrorUserDoesNotExistWithThatPassportInfoExeption = errors.New("user doesn't exist by that passport data")

	ErrorTaskDoesNotExist      = errors.New("task doesn't exist")
	ErrorNoAnyTasksForThisUser = errors.New("no any tasks for this user")
)
