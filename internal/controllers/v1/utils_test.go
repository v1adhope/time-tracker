package v1_test

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/configs"
	v1 "github.com/v1adhope/time-tracker/internal/controllers/v1"
	"github.com/v1adhope/time-tracker/internal/usecases"
	"github.com/v1adhope/time-tracker/internal/usecases/repositories"
	"github.com/v1adhope/time-tracker/pkg/logger"
	"github.com/v1adhope/time-tracker/pkg/postgresql"
)

func prepare() (*postgresql.Postgres, *gin.Engine) {
	cfg, err := configs.Build("../../../.env")
	if err != nil {
		log.Fatal(err)
	}

	appLog := logger.New(cfg.Logger.LogLevel)

	mainCtx := context.Background()

	postgres, err := postgresql.Build(mainCtx, cfg.Postgres)
	if err != nil {
		log.Fatal("can't get postgres pool")
	}

	postgres.Migrate("../../../migrations")

	repos := repositories.New(postgres)

	seeding(mainCtx, postgres)

	usecases := usecases.New(repos)

	if err := v1.RegisterCustomValidations(); err != nil {
		log.Fatal("can't register custom validations")
	}

	handler := gin.New()

	v1.Handle(&v1.Router{
		Handler:  handler,
		Usecases: usecases,
		Log:      appLog,
	})

	return postgres, handler
}

func seeding(ctx context.Context, postgres *postgresql.Postgres) {
	sql, args, _ := postgres.Builder.Insert("users").
		Columns("surname", "name", "patronymic", "address", "passport_number").
		Values("Funk", "Theresia", "Cummerata-Thompson", "53636 Gabrielle Mount", "3333 333333").
		Values("Runolfsdottir", "Violette", "Johns", "52265 Parker Crossroad", "3333 666666").
		Values("McCullough", "Jessie", "Waelchi", "8020 Dach Pine", "3333 444444").
		Values("Rippin", "Katrine", "Block", "985 N Jefferson Street", "5555 124041").
		Values("Schulist", "Kailee", "Fritsch", "5303 Church View", "2515 692797").
		ToSql()

	postgres.Pool.Exec(ctx, sql, args...)

	sql, args, _ = postgres.Builder.Insert("tasks").
		Columns("created_at", "finished_at", "user_id").
		Values("2024-01-16T09:08:25Z", "2024-01-16T16:10:00Z", getUserID(postgres, 3)).
		Values("2024-03-11T11:25:00Z", "2024-05-11T09:08:25Z", getUserID(postgres, 3)).
		Values("2024-04-16T09:08:25Z", "2024-05-16T09:08:25Z", getUserID(postgres, 3)).
		Values("2024-12-16T09:08:25Z", nil, getUserID(postgres, 3)).
		Values("2024-08-11T11:25:00Z", nil, getUserID(postgres, 3)).
		Values("2024-11-16T07:08:25Z", "2024-11-16T09:08:25Z", getUserID(postgres, 2)).
		Values("2024-05-18T11:00:00Z", "2024-05-20T09:08:25Z", getUserID(postgres, 2)).
		Values("2024-01-16T07:00:25Z", "2024-01-16T09:08:25Z", getUserID(postgres, 2)).
		Values("2024-03-16T00:08:25Z", "2024-03-24T00:00:00Z", getUserID(postgres, 2)).
		Values("2024-12-16T09:08:25Z", nil, getUserID(postgres, 4)).
		Values("2024-08-11T11:25:00Z", nil, getUserID(postgres, 4)).
		ToSql()

	postgres.Pool.Exec(ctx, sql, args...)
}

func getID(driver *postgresql.Postgres, table, column string, offset uint64) string {
	id := ""
	sql, args, _ := driver.Builder.Select(column).From(table).Limit(1).Offset(offset).ToSql()
	driver.Pool.QueryRow(context.Background(), sql, args...).Scan(&id)

	return id
}

func getUserID(driver *postgresql.Postgres, offset uint64) string {
	return getID(driver, "users", "user_id", offset)
}

func getTaskID(driver *postgresql.Postgres, offset uint64) string {
	return getID(driver, "tasks", "task_id", offset)
}
