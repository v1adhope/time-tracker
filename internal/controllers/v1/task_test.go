package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskStartPositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	testCases := []struct {
		key    string
		userID string
	}{
		{
			key:    "case 1",
			userID: getUserID(postgres, 0),
		},
		{
			key:    "case 2",
			userID: getUserID(postgres, 1),
		},
		{
			key:    "case 3",
			userID: getUserID(postgres, 2),
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/tasks/start/%s", tc.userID), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code, tc.key)
	}
}

func TestTaskStartNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	testCases := []struct {
		key          string
		userID       string
		expectedCode int
	}{
		{
			key:          "there's no user with that id",
			userID:       "1ef44bc9-5c77-6c90-a41a-1f1b0b522ea3",
			expectedCode: http.StatusNoContent,
		},
		{
			key:          "not correct type of id case 1",
			userID:       "5a55-6710-a9cf-ddslfjkjjj2e",
			expectedCode: http.StatusBadRequest,
		},
		{
			key:          "not correct type of id case 2",
			userID:       "1",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/tasks/start/%s", tc.userID), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
	}
}

func TestTaskEndPositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	testCases := []struct {
		key    string
		userID string
	}{
		{
			key:    "case 1",
			userID: getTaskID(postgres, 0),
		},
		{
			key:    "case 2",
			userID: getTaskID(postgres, 1),
		},
		{
			key:    "case 3",
			userID: getTaskID(postgres, 2),
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/end/%s", tc.userID), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, tc.key)
	}
}

func TestTaskEndNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	testCases := []struct {
		key          string
		userID       string
		expectedCode int
	}{
		{
			key:          "there's no user with that id",
			userID:       "1ef44bc9-5c77-6c90-a41a-1f1b0b522ea3",
			expectedCode: http.StatusNoContent,
		},
		{
			key:          "case 1 not correct type of id",
			userID:       "5a55-6710-a9cf-ddslfjkjjj2e",
			expectedCode: http.StatusBadRequest,
		},
		{
			key:          "case 2 not correct type of id",
			userID:       "1",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/end/%s", tc.userID), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
	}
}

func TestTaskSummaryTimePositive(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type input struct {
		id string
		// startTime string
		// endTime   string
	}
	type task struct {
		CreatedAt   string `json:"createdAt"`
		FinishedAt  string `json:"finishedAt"`
		SummaryTime string `json:"summaryTime"`
	}

	testCases := []struct {
		key         string
		input       input
		expectedLen int
	}{
		{
			key: "case 1",
			input: input{
				id: getUserID(postgres, 3),
			},
			expectedLen: 5,
		},
		{
			key: "case 2",
			input: input{
				id: getUserID(postgres, 2),
			},
			expectedLen: 4,
		},
		{
			key: "case 3",
			input: input{
				id: getUserID(postgres, 4),
			},
			expectedLen: 2,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/summary-time/%s", tc.input.id), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code, tc.key)

		tasks := make([]task, 0)
		json.NewDecoder(recorder.Body).Decode(&tasks)
		assert.Equal(t, tc.expectedLen, len(tasks))

		for _, task := range tasks {
			startTime, err := time.Parse(time.RFC3339, task.CreatedAt)
			assert.NoError(t, err, tc.key)

			if task.FinishedAt != "" {
				finishedTime, err := time.Parse(time.RFC3339, task.FinishedAt)
				assert.NoError(t, err, tc.key)

				diff := finishedTime.Sub(startTime)
				assert.Equal(t, fmt.Sprintf("%dh%dm", int(diff.Hours()), int(diff.Minutes())), task.SummaryTime)
			}
		}
	}
}

func TestTaskSummaryTimeNegative(t *testing.T) {
	postgres, handler := prepare()
	defer postgres.Close()

	type input struct {
		id string
		// startTime string
		// endTime   string
	}

	testCases := []struct {
		key          string
		input        input
		expectedCode int
	}{
		{
			key: "there's in no user with that id",
			input: input{
				id: "1ef44ce4-6afb-6da0-9e4e-6ea3cb7df39c",
			},
			expectedCode: http.StatusNoContent,
		},
		{
			key: "case 1 wrong id type",
			input: input{
				id: "1ef44ce4-6da0-9e4e-6ea3cb7df39c",
			},
			expectedCode: http.StatusBadRequest,
		},
		{
			key: "case 2 wrong id type",
			input: input{
				id: "2",
			},
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/summary-time/%s", tc.input.id), nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
	}
}
