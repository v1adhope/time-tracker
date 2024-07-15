package repositories

import "github.com/v1adhope/time-tracker/pkg/postgresql"

type Repos struct {
	User *UserRepo
}

func New(driver *postgresql.Postgres) *Repos {
	return &Repos{
		User: NewUser(driver),
	}
}
