package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/v1adhope/time-tracker/internal/entities"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

type UserRepo struct {
	Driver *postgresql.Postgres
}

func NewUser(d *postgresql.Postgres) *UserRepo {
	return &UserRepo{d}
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
		return entities.User{}, fmt.Errorf("repositories: user: create: tosql: %w", err)
	}

	if err := r.Driver.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "users_passport_number_key" {
			return entities.User{}, entities.ErrorUserHasAlreadyExist
		}

		return entities.User{}, fmt.Errorf("repositories: user: create: queryrow: %w", err)
	}

	return user, nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	whereStatement := squirrel.Eq{
		"user_id": id,
	}

	sql, args, err := r.Driver.Builder.Delete("users").
		Where(whereStatement).
		ToSql()
	if err != nil {
		return fmt.Errorf("repositories: user: delete: tosql: %w", err)
	}

	tag, err := r.Driver.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repositories: user: delete: exec: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return entities.ErrorUserDoesNotExist
	}

	return nil
}

func (r *UserRepo) Update(ctx context.Context, user entities.User) error {
	whereStatement := squirrel.Eq{
		"user_id": user.ID,
	}

	valuesByColumns := squirrel.Eq{}

	if user.Surname != "" {
		valuesByColumns["surname"] = user.Surname
	}

	if user.Name != "" {
		valuesByColumns["name"] = user.Name
	}

	if user.Patronymic != "" {
		valuesByColumns["patronymic"] = user.Patronymic
	}

	if user.Address != "" {
		valuesByColumns["address"] = user.Address
	}

	if user.PassportNumber != "" {
		valuesByColumns["passport_number"] = user.PassportNumber
	}

	sql, args, err := r.Driver.Builder.Update("users").
		Where(whereStatement).
		SetMap(valuesByColumns).
		ToSql()
	if err != nil {
		return fmt.Errorf("repositories: user: update: tosql: %w", err)
	}

	tag, err := r.Driver.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("repositories: user: update: exec: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return entities.ErrorUserDoesNotExist
	}

	return nil
}

func (r *UserRepo) GetAll(ctx context.Context, representation entities.UserRepresentation) ([]entities.User, error) {
	sql, args, err := r.Driver.Builder.Select("user_id", "surname", "name", "patronymic", "address", "passport_number").
		From("users").
		Where(r.buildGetAllWhereFilterStatement(representation.Filter)).
		Limit(setLimitStatement(representation.Pagination.Limit)).
		Offset(setOffsetStatement(representation.Pagination.Offset)).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("repositories: user: getall: tosql: %w", err)
	}

	rows, err := r.Driver.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("repositories: user: getall: query: %w", err)
	}

	users := make([]entities.User, 0)
	user := entities.User{}

	_, err = pgx.ForEachRow(rows, []any{&user.ID, &user.Surname, &user.Name, &user.Patronymic, &user.Address, &user.PassportNumber}, func() error {
		users = append(users, user)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repositories: user: getall: forEachRow: %w", err)
	}

	return users, nil
}

func (r *UserRepo) buildGetAllWhereFilterStatement(filter entities.UserFilter) squirrel.And {
	statement := squirrel.And{}

	if filter.BySurname != "" {
		parts := strings.Split(filter.BySurname, ":")

		if isOperationEq(parts[0]) {
			statement = append(statement, squirrel.Eq{
				"surname": parts[1],
			})
		}

		if isOperationIlike(parts[0]) {
			statement = append(statement, squirrel.ILike{
				"surname": fmt.Sprintf("%%%s%%", parts[1]),
			})
		}
	}

	if filter.ByName != "" {
		parts := strings.Split(filter.ByName, ":")

		if isOperationEq(parts[0]) {
			statement = append(statement, squirrel.Eq{
				"name": parts[1],
			})
		}

		if isOperationIlike(parts[0]) {
			statement = append(statement, squirrel.ILike{
				"name": fmt.Sprintf("%%%s%%", parts[1]),
			})
		}
	}

	if filter.ByPatronymic != "" {
		parts := strings.Split(filter.ByPatronymic, ":")

		if isOperationEq(parts[0]) {
			statement = append(statement, squirrel.Eq{
				"patronymic": parts[1],
			})
		}

		if isOperationIlike(parts[0]) {
			statement = append(statement, squirrel.ILike{
				"patronymic": fmt.Sprintf("%%%s%%", parts[1]),
			})
		}
	}

	if filter.ByAddress != "" {
		parts := strings.Split(filter.ByAddress, ":")

		if isOperationEq(parts[0]) {
			statement = append(statement, squirrel.Eq{
				"address": parts[1],
			})
		}

		if isOperationIlike(parts[0]) {
			statement = append(statement, squirrel.ILike{
				"address": fmt.Sprintf("%%%s%%", parts[1]),
			})
		}
	}

	if filter.ByPassportNumber != "" {
		parts := strings.Split(filter.ByPassportNumber, ":")

		if isOperationEq(parts[0]) {
			statement = append(statement, squirrel.Eq{
				"passport_number": parts[1],
			})
		}

		if isOperationIlike(parts[0]) {
			statement = append(statement, squirrel.ILike{
				"passport_number": fmt.Sprintf("%%%s%%", parts[1]),
			})
		}
	}

	return statement
}
