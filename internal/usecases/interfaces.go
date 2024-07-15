package usecases

import (
	"context"

	"github.com/v1adhope/time-tracker/internal/entities"
)

type User interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user entities.User) error
}

type UserRepo interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, user entities.User) error
}
