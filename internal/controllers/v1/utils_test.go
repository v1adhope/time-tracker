package v1_test

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/v1adhope/time-tracker/internal/configs"
	v1 "github.com/v1adhope/time-tracker/internal/controllers/v1"
	"github.com/v1adhope/time-tracker/internal/entities"
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

	seeding(mainCtx, repos)

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

func seeding(ctx context.Context, repos *repositories.Repos) {
	users := []entities.User{
		{
			Surname:        "Funk",
			Name:           "Theresia",
			Patronymic:     "Cummerata-Thompson",
			Address:        "53636 Gabrielle Mount",
			PassportNumber: "3333 333333",
		},
		{
			Surname:        "Runolfsdottir",
			Name:           "Violette",
			Patronymic:     "Johns",
			Address:        "52265 Parker Crossroad",
			PassportNumber: "3333 666666",
		},
		{
			Surname:        "McCullough",
			Name:           "Jessie",
			Patronymic:     "Waelchi",
			Address:        "8020 Dach Pine",
			PassportNumber: "3333 444444",
		},
		{
			Surname:        "Rippin",
			Name:           "Katrine",
			Patronymic:     "Block",
			Address:        "985 N Jefferson Street",
			PassportNumber: "5555 124041",
		},
		{
			Surname:        "Schulist",
			Name:           "Kailee",
			Patronymic:     "Fritsch",
			Address:        "5303 Church View",
			PassportNumber: "2515 692797",
		},
	}

	for no, user := range users {
		userWithID, _ := repos.User.Create(ctx, user)
		users[no].ID = userWithID.ID
	}

	for i := 0; i < 3; i++ {
		task, _ := repos.Task.Create(ctx, users[3].ID)
		repos.Task.SetFinishedAt(ctx, task.ID)
	}

	for i := 0; i < 2; i++ {
		repos.Task.Create(ctx, users[3].ID)
	}

	for i := 0; i < 4; i++ {
		task, _ := repos.Task.Create(ctx, users[2].ID)
		repos.Task.SetFinishedAt(ctx, task.ID)
	}

	for i := 0; i < 2; i++ {
		repos.Task.Create(ctx, users[4].ID)
	}
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
