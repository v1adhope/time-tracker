package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	createdAt := time.Now().UTC()

	task := entities.Task{
		UserID:    userID,
		CreatedAt: createdAt,
	}

	valuesByColumns := squirrel.Eq{
		"user_id":    task.UserID,
		"created_at": task.CreatedAt,
	}

	sql, args, err := r.Driver.Builder.Insert("tasks").
		SetMap(valuesByColumns).
		Suffix("returning \"task_id\"").
		ToSql()
	if err != nil {
		return entities.Task{}, fmt.Errorf("repositories: task: create: tosql: %w", err)
	}

	if err := r.Driver.Pool.QueryRow(ctx, sql, args...).Scan(&task.ID); err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.ConstraintName == "fk_tasks_users_user_id" {
			return entities.Task{}, entities.ErrorUserDoesNotExist
		}

		return entities.Task{}, fmt.Errorf("repositories: task: create: queryRow: %w", err)
	}

	return task, nil
}

func (r *TaskRepo) SetFinishedAt(ctx context.Context, id string) (time.Time, error) {
	finishedAt := time.Now().UTC()

	valuesByColumns := squirrel.Eq{
		"finished_at": finishedAt,
	}

	whereStatement := squirrel.Eq{
		"task_id": id,
	}

	sql, args, err := r.Driver.Builder.Update("tasks").
		SetMap(valuesByColumns).
		Where(whereStatement).
		ToSql()
	if err != nil {
		return time.Time{}, fmt.Errorf("repositories: task: setFinishedAt: tosql: %w", err)
	}

	tag, err := r.Driver.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return time.Time{}, fmt.Errorf("repositories: task: setFinishedAt: exec: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return time.Time{}, entities.ErrorTaskDoesNotExist
	}

	return finishedAt, nil
}
