package usecases

import (
	"context"

	"github.com/v1adhope/time-tracker/internal/entities"
)

type User interface {
	Create(ctx context.Context, user entities.User) (string, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user entities.User) error
	GetAll(ctx context.Context, representation entities.UserRepresentation) ([]entities.User, error)
	Get(ctx context.Context, passportNumber string) (entities.User, error)
}

type UserRepo interface {
	Create(ctx context.Context, user entities.User) (string, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user entities.User) error
	GetAll(ctx context.Context, representation entities.UserRepresentation) ([]entities.User, error)
	Get(ctx context.Context, passportNumber string) (entities.User, error)
}

type Task interface {
	Start(ctx context.Context, userID string) error
	End(ctx context.Context, id string) (string, error)
	GetReportSummaryTime(ctx context.Context, userID string, sort entities.TaskSort) ([]entities.TaskSummary, error)
}

type TaskRepo interface {
	Create(ctx context.Context, userID string) (string, error)
	SetFinishedAt(ctx context.Context, id string) (string, error)
	GetReportSummaryTime(ctx context.Context, userID string, sort entities.TaskSort) ([]entities.TaskSummary, error)
}
