package usecases

import "github.com/v1adhope/time-tracker/internal/usecases/repositories"

type Usecases struct {
	User *UserUsecase
}

func New(repos *repositories.Repos) *Usecases {
	return &Usecases{
		User: NewUser(repos.User),
	}
}
