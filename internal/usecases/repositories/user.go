package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/v1adhope/time-tracker/internal/entities"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

type UserRepo struct {
	Driver *postgresql.Postgres
}

func NewUser(driver *postgresql.Postgres) *UserRepo {
	return &UserRepo{driver}
}

func (r *UserRepo) Create(ctx context.Context, user entities.User) (entities.User, error) {
	valuesByColumns := squirrel.Eq{
		"surname":         user.Surname,
		"name":            user.Name,
		"patronymic":      user.Patronymic,
		"address":         user.Address,
		"passport_number": user.PassportNumber,
	}

	sql, args, err := r.Driver.Builder.Insert("users").
		SetMap(valuesByColumns).
		Suffix("returning \"user_id\"").
		ToSql()
	if err != nil {
		return entities.User{}, fmt.Errorf("repositories: create: tosql: %w", err)
	}

	if err := r.Driver.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "users_passport_number_key" {
			return entities.User{}, entities.ErrorUserHasAlreadyExist
		}

		return entities.User{}, fmt.Errorf("repositories: create: queryrow: %w", err)
	}

	return user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	valuesByColumns := squirrel.Eq{
		"user_id": id,
	}

	sql, args, err := r.Driver.Builder.Delete("users").Where(valuesByColumns).ToSql()
	if err != nil {
		return fmt.Errorf("repositories: delete: tosql: %w", err)
	}

	tag, err := r.Driver.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repositories: delete: exec: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return entities.ErrorUserDoesNotExist
	}

	return nil
}
