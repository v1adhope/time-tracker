package usecases

import (
	"context"

	"github.com/v1adhope/time-tracker/internal/entities"
)

type UserUsecase struct {
	userRepo UserRepo
}

func NewUser(ur UserRepo) *UserUsecase {
	return &UserUsecase{ur}
}

func (u *UserUsecase) Create(ctx context.Context, user entities.User) (entities.User, error) {
	user, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}
