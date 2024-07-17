package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
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
	createdAt := time.Now().UTC().Format(time.RFC3339)

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
			return entities.Task{}, entities.ErrorUsersDoesNotExist
		}

		return entities.Task{}, fmt.Errorf("repositories: task: create: queryRow: %w", err)
	}

	return task, nil
}

func (r *TaskRepo) SetFinishedAt(ctx context.Context, id string) (string, error) {
	finishedAt := time.Now().UTC().Format(time.RFC3339)

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
		return "", fmt.Errorf("repositories: task: setFinishedAt: tosql: %w", err)
	}

	tag, err := r.Driver.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return "", fmt.Errorf("repositories: task: setFinishedAt: exec: %w", err)
	}

	if tag.RowsAffected() != 1 {
		return "", entities.ErrorTaskDoesNotExist
	}

	return finishedAt, nil
}

type taskSummaryTimeDTO struct {
	ID          string
	CreatedAt   time.Time
	FinishedAt  *time.Time
	SummaryTime *time.Duration
}

func (r *TaskRepo) GetReportSummaryTime(ctx context.Context, userID string, sort entities.TaskSort) ([]entities.TaskSummary, error) {
	whereStatement := squirrel.Eq{
		"user_id": userID,
	}

	sql, args, err := r.Driver.Builder.Select("task_id", "created_at", "finished_at", "finished_at - created_at as summary_time").
		From("tasks").
		Where(whereStatement).
		Where(r.buildGetReportSummaryTimeWhereSortStatement(sort)).
		OrderBy("summary_time desc").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("repositories: task: getReportSummaryTime: tosql: %w", err)
	}

	rows, err := r.Driver.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("repositories: task: getReportSummaryTime: query: %w", err)
	}

	tasks := make([]entities.TaskSummary, 0)
	taskDTO := taskSummaryTimeDTO{}

	_, err = pgx.ForEachRow(rows, []any{&taskDTO.ID, &taskDTO.CreatedAt, &taskDTO.FinishedAt, &taskDTO.SummaryTime}, func() error {
		task := entities.TaskSummary{
			ID:        taskDTO.ID,
			CreatedAt: taskDTO.CreatedAt.Format(time.RFC3339),
		}

		if taskDTO.FinishedAt != nil {
			task.FinishedAt = taskDTO.FinishedAt.Format(time.RFC3339)
		}

		if taskDTO.SummaryTime != nil {
			task.SummaryTime = fmt.Sprintf("%dh%dm", int(taskDTO.SummaryTime.Hours()), int(taskDTO.SummaryTime.Minutes()))
		}

		tasks = append(tasks, task)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("repositories: task: getReportSummaryTime: forEachRow: %w", err)
	}

	if len(tasks) == 0 {
		return nil, entities.ErrorNoAnyTasksForThisUser
	}

	return tasks, nil
}

func (r *TaskRepo) buildGetReportSummaryTimeWhereSortStatement(sort entities.TaskSort) squirrel.And {
	if sort.StartTime == "" && sort.EndTime == "" {
		return nil
	}

	statement := squirrel.And{}

	if sort.StartTime != "" {
		statement = append(statement, squirrel.GtOrEq{
			"created_at": sort.StartTime,
		})
	}

	if sort.EndTime != "" {
		statement = append(statement, squirrel.LtOrEq{
			"finished_at": sort.EndTime,
		})
	}

	return statement
}
