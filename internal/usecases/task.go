package usecases

import (
	"context"

	"github.com/v1adhope/time-tracker/internal/entities"
)

type TaskUsecase struct {
	TaskRepo TaskRepo
}

func NewTask(tr TaskRepo) *TaskUsecase {
	return &TaskUsecase{tr}
}

func (u *TaskUsecase) Start(ctx context.Context, userID string) (entities.Task, error) {
	task, err := u.TaskRepo.Create(ctx, userID)
	if err != nil {
		return entities.Task{}, err
	}

	return task, nil
}
