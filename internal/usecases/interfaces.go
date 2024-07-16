package usecases

import (
	"context"

	"github.com/v1adhope/time-tracker/internal/entities"
)

type User interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user entities.User) error
	GetAll(ctx context.Context, representation entities.UserRepresentation) ([]entities.User, error)
}

type UserRepo interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user entities.User) error
	GetAll(ctx context.Context, representation entities.UserRepresentation) ([]entities.User, error)
}

type Task interface {
	Start(ctx context.Context, userID string) (entities.Task, error)
	End(ctx context.Context, id string) (string, error)
	GetReportSummaryTime(ctx context.Context, userID string, sort entities.TaskSort) ([]entities.TaskSummary, error)
}

type TaskRepo interface {
	Create(ctx context.Context, userID string) (entities.Task, error)
	SetFinishedAt(ctx context.Context, id string) (string, error)
	GetReportSummaryTime(ctx context.Context, userID string, sort entities.TaskSort) ([]entities.TaskSummary, error)
}
