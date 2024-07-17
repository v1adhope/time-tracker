package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

func getUserIDByOffset(driver *postgresql.Postgres, offset uint64) string {
	id := ""
	sql, args, _ := driver.Builder.Select("user_id").From("users").Limit(1).Offset(offset).ToSql()
	driver.Pool.QueryRow(context.Background(), sql, args...).Scan(&id)

	return id
}

func TestCreatePositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type user struct {
		Surname        string `json:"surname"`
		Name           string `json:"name"`
		Patronymic     string `json:"patronymic"`
		Address        string `json:"address"`
		PassportNumber string `json:"passportNumber"`
	}

	testCases := []struct {
		key  string
		user user
	}{
		{
			key: "case 1",
			user: user{
				Surname:        "Bode",
				Name:           "Rogers",
				Patronymic:     "Robertovich",
				Address:        "1123 Ola Brook",
				PassportNumber: "6666 888888",
			},
		},
		{
			key: "case 2",
			user: user{
				Surname:        "Ondricka",
				Name:           "Coby",
				Patronymic:     "Victorovich",
				Address:        "9312 Weber Neck",
				PassportNumber: "6666 666666",
			},
		},
		{
			key: "case 3",
			user: user{
				Surname:        "Shanahan",
				Name:           "Timothy",
				Patronymic:     "Jacobson",
				Address:        "78510 Howard Street",
				PassportNumber: "8888 666666",
			},
		},
	}

	for _, tc := range testCases {
		bytesBody, _ := json.Marshal(&tc.user)
		req, _ := http.NewRequest("POST", "/v1/users/", strings.NewReader(string(bytesBody)))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, tc.key)

		user := user{}
		json.NewDecoder(recorder.Body).Decode(&user)
		assert.Equal(t, tc.user, user, tc.key)
	}
}

func TestCreateNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type user struct {
		Surname        string `json:"surname"`
		Name           string `json:"name"`
		Patronymic     string `json:"patronymic"`
		Address        string `json:"address"`
		PassportNumber string `json:"passportNumber"`
	}

	testCases := []struct {
		key   string
		input user
	}{
		{
			key: "wrong passport",
			input: user{
				Surname:        "Bode",
				Name:           "Rogers",
				Patronymic:     "Robertovich",
				Address:        "1123 Ola Brook",
				PassportNumber: "6666 8888889",
			},
		},
		{
			key: "forgot field",
			input: user{
				Surname:        "Ondricka",
				Name:           "Coby",
				Address:        "9312 Weber Neck",
				PassportNumber: "3333 333333",
			},
		},
	}

	for _, tc := range testCases {
		bytesBody, _ := json.Marshal(&tc.input)
		req, _ := http.NewRequest("POST", "/v1/users/", strings.NewReader(string(bytesBody)))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, tc.key)
	}
}

func TestUpdatePositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type user struct {
		ID             string `json:"-"`
		Surname        string `json:"surname"`
		Name           string `json:"name"`
		Patronymic     string `json:"patronymic"`
		Address        string `json:"address"`
		PassportNumber string `json:"passportNumber"`
	}

	testCases := []user{
		{
			ID:             getUserIDByOffset(postgres, 0),
			Surname:        "Sporer",
			Name:           "Lemuel",
			Patronymic:     "Schultz",
			Address:        "3042 Nicolas Summit",
			PassportNumber: "4444 664656",
		},
		{
			ID:             getUserIDByOffset(postgres, 1),
			Surname:        "Reilly",
			Name:           "Clara",
			Patronymic:     "Mohr",
			Address:        "1720 Schmeler Road",
			PassportNumber: "4424 664656",
		},
		{
			ID:             getUserIDByOffset(postgres, 2),
			Surname:        "Farrell",
			Name:           "Nona",
			Patronymic:     "Wyman-Lockman",
			Address:        "706 Willms Ranch",
			PassportNumber: "4324 664656",
		},
	}

	for _, tc := range testCases {
		userJSON, _ := json.Marshal(&tc)
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/users/%s", tc.ID), strings.NewReader(string(userJSON)))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, tc.ID)
	}
}

func TestUpdateNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type user struct {
		id             string
		Surname        string `json:"surname"`
		Name           string `json:"name"`
		Patronymic     string `json:"patronymic"`
		Address        string `json:"address"`
		PassportNumber string `json:"passportNumber"`
	}

	testCases := []struct {
		key          string
		input        user
		expectedCode int
	}{
		{
			key: "wrong passport number (5 numbers instead 6)",
			input: user{
				id:             getUserIDByOffset(postgres, 0),
				Surname:        "Sporer",
				Name:           "Lemuel",
				Patronymic:     "Schultz",
				Address:        "3042 Nicolas Summit",
				PassportNumber: "4444 66465",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong passport number (number constist not number)",
			input: user{
				id:             getUserIDByOffset(postgres, 1),
				Surname:        "Reilly",
				Name:           "Clara",
				Patronymic:     "Mohr",
				Address:        "1720 Schmeler Road",
				PassportNumber: "4424 66465s",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong passport number (serial consist 5 numbers insetad 4)",
			input: user{
				id:             getUserIDByOffset(postgres, 2),
				Surname:        "Farrell",
				Name:           "Nona",
				Patronymic:     "Wyman-Lockman",
				Address:        "706 Willms Ranch",
				PassportNumber: "43241 664656",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong uuid",
			input: user{
				id:             "1ef442cb-6f00-9f30-4cb600d007df",
				Surname:        "Farrell",
				Name:           "Nona",
				Patronymic:     "Wyman-Lockman",
				Address:        "706 Willms Ranch",
				PassportNumber: "4324 664656",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "passport number has already exist",
			input: user{
				id:             getUserIDByOffset(postgres, 4),
				Surname:        "Farrell",
				Name:           "Nona",
				Patronymic:     "Wyman-Lockman",
				Address:        "706 Willms Ranch",
				PassportNumber: "3333 333333",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "user doesn't exist",
			input: user{
				id:             "1ef442cb-db96-6f00-9f30-4cb600d007df",
				Surname:        "Sporer",
				Name:           "Lemuel",
				Patronymic:     "Schultz",
				Address:        "3042 Nicolas Summit",
				PassportNumber: "4444 664659",
			},
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		bytesBody, _ := json.Marshal(&tc.input)
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/users/%s", tc.input.id), strings.NewReader(string(bytesBody)))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
	}
}

func TestDeletePositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	testCases := []struct {
		id string
	}{
		{
			id: getUserIDByOffset(postgres, 0),
		},
		{
			id: getUserIDByOffset(postgres, 1),
		},
		{
			id: getUserIDByOffset(postgres, 2),
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/users/%s", tc.id), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, tc.id)
	}
}

func TestDeleteNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	testCases := []struct {
		key          string
		id           string
		expectedCode int
	}{
		{
			key:          "not uuid case 1",
			id:           "1",
			expectedCode: http.StatusBadRequest,
		},
		{
			key:          "not uuid case 2",
			id:           "slfa",
			expectedCode: http.StatusBadRequest,
		},
		{
			key:          "not uuid case 3",
			id:           "1ef442cb-bf1b-6c40-a618029ec695",
			expectedCode: http.StatusBadRequest,
		},
		{
			key:          "uuid doesn't exist case 1",
			id:           "1ef442cb-bf1b-6c40-add5-a618029ec695",
			expectedCode: http.StatusNoContent,
		},
		{
			key:          "uuid doesn't exist case 2",
			id:           "1ef442cb-db96-6f00-9f30-4cb600d007df",
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/users/%s", tc.id), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
	}
}

func TestGetAllPositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type user struct {
		Surname        string `json:"surname"`
		Name           string `json:"name"`
		Patronymic     string `json:"patronymic"`
		Address        string `json:"address"`
		PassportNumber string `json:"passportNumber"`
	}
	type input struct {
		bySurname        string
		byName           string
		byPatronymic     string
		byAddress        string
		byPassportNumber string

		limit  string
		offset string
	}

	testCases := []struct {
		key      string
		input    input
		expected []user
	}{
		{
			key: "no query test",
			expected: []user{
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
			},
		},
		{
			key: "pagination only",
			input: input{
				limit:  "2",
				offset: "1",
			},
			expected: []user{
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
			},
		},
		{
			key: "filter only equal values",
			input: input{
				bySurname:        "eq:Runolfsdottir",
				byName:           "eq:Violette",
				byPatronymic:     "eq:Johns",
				byAddress:        "eq:52265 Parker Crossroad",
				byPassportNumber: "eq:3333 666666",
			},
			expected: []user{
				{
					Surname:        "Runolfsdottir",
					Name:           "Violette",
					Patronymic:     "Johns",
					Address:        "52265 Parker Crossroad",
					PassportNumber: "3333 666666",
				},
			},
		},
		{
			key: "filter only ilike values",
			input: input{
				bySurname:        "ilike:Rip",
				byName:           "ilike:tri",
				byPatronymic:     "ilike:ck",
				byAddress:        "ilike:Jefferson",
				byPassportNumber: "ilike:5555",
			},
			expected: []user{
				{
					Surname:        "Rippin",
					Name:           "Katrine",
					Patronymic:     "Block",
					Address:        "985 N Jefferson Street",
					PassportNumber: "5555 124041",
				},
			},
		},
		{
			key: "pagination with filter",
			input: input{
				bySurname: "ilike:u",
				limit:     "2",
				offset:    "1",
			},
			expected: []user{
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
			},
		},
	}

	for _, tc := range testCases {
		query := url.Values{}

		query.Set("limit", tc.input.limit)
		query.Set("offset", tc.input.offset)
		query.Set("surname", tc.input.bySurname)
		query.Set("name", tc.input.byName)
		query.Set("patronymic", tc.input.byPatronymic)
		query.Set("address", tc.input.byAddress)
		query.Set("passportNumber", tc.input.byPassportNumber)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/?%s", query.Encode()), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, tc.key)

		users := make([]user, 0)
		json.NewDecoder(recorder.Body).Decode(&users)
		assert.Equal(t, tc.expected, users, tc.key)
	}
}

func TestGetAllNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type input struct {
		bySurname        string
		byName           string
		byPatronymic     string
		byAddress        string
		byPassportNumber string

		limit  string
		offset string
	}

	testCases := []struct {
		key          string
		input        input
		expectedCode int
	}{
		{
			key: "with limit 0",
			input: input{
				limit: "0",
			},
			expectedCode: http.StatusNoContent,
		},
		{
			key: "wrong surname filter",
			input: input{
				bySurname: "notEq:Funk",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong name filter",
			input: input{
				byName: "like:Theresia",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong patronymic filter",
			input: input{
				byPatronymic: "notLike:Cummerata-Thompson",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong address filter",
			input: input{
				byPatronymic: "lte:53636",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "wrong passportNumber filter",
			input: input{
				byPassportNumber: "=:53636",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "offset not uint64",
			input: input{
				offset: "twoonethree",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		query := url.Values{}

		query.Set("limit", tc.input.limit)
		query.Set("offset", tc.input.offset)
		query.Set("surname", tc.input.bySurname)
		query.Set("name", tc.input.byName)
		query.Set("patronymic", tc.input.byPatronymic)
		query.Set("address", tc.input.byAddress)
		query.Set("passportNumber", tc.input.byPassportNumber)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/?%s", query.Encode()), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
	}
}

func TestInfoPositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type user struct {
		Surname    string `json:"surname"`
		Name       string `json:"name"`
		Patronymic string `json:"patronymic"`
		Address    string `json:"address"`
	}
	type input struct {
		passportSeries string
		passportNumber string
	}

	testCases := []struct {
		key      string
		input    input
		expected user
	}{
		{
			key: "case 1",
			input: input{
				passportSeries: "3333",
				passportNumber: "333333",
			},
			expected: user{
				Surname:    "Funk",
				Name:       "Theresia",
				Patronymic: "Cummerata-Thompson",
				Address:    "53636 Gabrielle Mount",
			},
		},
	}

	for _, tc := range testCases {
		quaery := url.Values{}

		quaery.Set("passportSeries", tc.input.passportSeries)
		quaery.Set("passportNumber", tc.input.passportNumber)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/info?%s", quaery.Encode()), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, tc.key)

		user := user{}
		json.NewDecoder(recorder.Body).Decode(&user)
		assert.Equal(t, tc.expected, user, tc.key)
	}
}

func TestInfoNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type input struct {
		passportSeries string
		passportNumber string
	}

	testCase := []struct {
		key   string
		input input
	}{
		{
			key: "case 0",
			input: input{
				passportSeries: "1247",
				passportNumber: "95829s",
			},
		},
		{
			key: "case 1",
			input: input{
				passportSeries: "11111",
				passportNumber: "958295",
			},
		},
		{
			key: "case 2",
			input: input{
				passportSeries: "1111",
				passportNumber: "95829",
			},
		},
	}

	for _, tc := range testCase {
		quaery := url.Values{}

		quaery.Set("passportSeries", tc.input.passportSeries)
		quaery.Set("passportNumber", tc.input.passportNumber)

		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/users/info?%s", quaery.Encode()), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code, tc.key)
	}
}
