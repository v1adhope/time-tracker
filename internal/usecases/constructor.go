package usecases

import "github.com/v1adhope/time-tracker/internal/usecases/repositories"

type Usecases struct {
	User *UserUsecase
	Task *TaskUsecase
}

func New(repos *repositories.Repos) *Usecases {
	return &Usecases{
		User: NewUser(repos.User),
		Task: NewTask(repos.Task),
	}
}
