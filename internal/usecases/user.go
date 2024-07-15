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

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	if err := u.userRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) Update(ctx context.Context, user entities.User) error {
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) GetAll(ctx context.Context, representation entities.UserRepresentation) ([]entities.User, error) {
	users, err := u.userRepo.GetAll(ctx, representation)
	if err != nil {
		return nil, err
	}

	return users, nil
}
