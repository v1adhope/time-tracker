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

type TaskRepo struct {
	Driver *postgresql.Postgres
}

func NewTask(d *postgresql.Postgres) *TaskRepo {
	return &TaskRepo{d}
}

func (r *TaskRepo) Create(ctx context.Context, userID string) (entities.Task, error) {
	task := entities.Task{UserID: userID}

	valuesByColumns := squirrel.Eq{
		"user_id": task.UserID,
	}
	sql, args, err := r.Driver.Builder.Insert("tasks").
		SetMap(valuesByColumns).
		Suffix("returning \"task_id\", \"created_at\"").
		ToSql()
	if err != nil {
		return entities.Task{}, fmt.Errorf("repositories: task: create: tosql: %w", err)
	}

	if err := r.Driver.Pool.QueryRow(ctx, sql, args...).Scan(&task.ID, &task.CreatedAt); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "fk_tasks_users_user_id" {
			return entities.Task{}, entities.ErrorUserDoesNotExist
		}

		return entities.Task{}, fmt.Errorf("repositories: task: create: queryRow: %w", err)
	}

	return task, nil
}
