package repositories

import "strconv"

const (
	defaultLimit  = 10
	defaultOffset = 0
)

func isOperationEq(target string) bool {
	if target != "eq" {
		return false
	}

	return true
}

func isOperationIlike(target string) bool {
	if target != "ilike" {
		return false
	}

	return true
}

func setLimitStatement(target string) uint64 {
	if target == "" {
		return defaultLimit
	}

	value, err := strconv.ParseUint(target, 10, 64)
	if err != nil {
		return defaultLimit
	}

	return value
}

func setOffsetStatement(target string) uint64 {
	if target == "" {
		return defaultOffset
	}

	value, err := strconv.ParseUint(target, 10, 64)
	if err != nil {
		return defaultOffset
	}

	return value
}
