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

func (u *TaskUsecase) End(ctx context.Context, id string) (string, error) {
	finishedAt, err := u.TaskRepo.SetFinishedAt(ctx, id)
	if err != nil {
		return "", err
	}

	return finishedAt, nil
}

func (u *TaskUsecase) GetReportSummaryTime(ctx context.Context, userID string, sort entities.TaskSort) ([]entities.TaskSummary, error) {
	tasks, err := u.TaskRepo.GetReportSummaryTime(ctx, userID, sort)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
