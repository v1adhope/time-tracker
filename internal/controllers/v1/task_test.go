package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskStartPositive(t *testing.T) {
	postgres, handler := prepare()
	t.Cleanup(func() {
		defer postgres.Close()
	})

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
		t.Run(tc.key, func(t *testing.T) {
			req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/tasks/start/%s", tc.userID), nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusCreated, recorder.Code, tc.key)
		})
	}
}

func TestTaskStartNegative(t *testing.T) {
	postgres, handler := prepare()
	t.Cleanup(func() {
		defer postgres.Close()
	})

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
		t.Run(tc.key, func(t *testing.T) {
			req, _ := http.NewRequest("POST", fmt.Sprintf("/v1/tasks/start/%s", tc.userID), nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
		})
	}
}

func TestTaskEndPositive(t *testing.T) {
	postgres, handler := prepare()
	t.Cleanup(func() {
		defer postgres.Close()
	})

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
		t.Run(tc.key, func(t *testing.T) {
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/end/%s", tc.userID), nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code, tc.key)
		})
	}
}

func TestTaskEndNegative(t *testing.T) {
	postgres, handler := prepare()
	t.Cleanup(func() {
		defer postgres.Close()
	})

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
		t.Run(tc.key, func(t *testing.T) {
			req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/tasks/end/%s", tc.userID), nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
		})
	}
}

func TestTaskSummaryTimePositive(t *testing.T) {
	postgres, handler := prepare()
	t.Cleanup(func() {
		defer postgres.Close()
	})

	type input struct {
		id        string
		startTime string
		endTime   string
	}
	type task struct {
		CreatedAt   string `json:"createdAt"`
		FinishedAt  string `json:"finishedAt"`
		SummaryTime string `json:"summaryTime"`
	}

	testCases := []struct {
		key      string
		input    input
		expected []task
	}{
		{
			key: "tasks' 4th user",
			input: input{
				id: getUserID(postgres, 3),
			},
			expected: []task{
				{
					CreatedAt: "2024-12-16T09:08:25Z",
				},
				{
					CreatedAt: "2024-08-11T11:25:00Z",
				},
				{
					CreatedAt:   "2024-03-11T11:25:00Z",
					FinishedAt:  "2024-05-11T09:08:25Z",
					SummaryTime: "1461h43m",
				},
				{
					CreatedAt:   "2024-04-16T09:08:25Z",
					FinishedAt:  "2024-05-16T09:08:25Z",
					SummaryTime: "720h0m",
				},
				{
					CreatedAt:   "2024-01-16T09:08:25Z",
					FinishedAt:  "2024-01-16T16:10:00Z",
					SummaryTime: "7h2m",
				},
			},
		},
		{
			key: "tasks' 3rd user",
			input: input{
				id: getUserID(postgres, 2),
			},
			expected: []task{
				{
					CreatedAt:   "2024-03-16T00:08:25Z",
					FinishedAt:  "2024-03-24T00:00:00Z",
					SummaryTime: "191h52m",
				},
				{
					CreatedAt:   "2024-05-18T11:00:00Z",
					FinishedAt:  "2024-05-20T09:08:25Z",
					SummaryTime: "46h8m",
				},
				{
					CreatedAt:   "2024-01-16T07:00:25Z",
					FinishedAt:  "2024-01-16T09:08:25Z",
					SummaryTime: "2h8m",
				},
				{
					CreatedAt:   "2024-11-16T07:08:25Z",
					FinishedAt:  "2024-11-16T09:08:25Z",
					SummaryTime: "2h0m",
				},
			},
		},
		{
			key: "tasks' 5 user",
			input: input{
				id: getUserID(postgres, 4),
			},
			expected: []task{
				{
					CreatedAt: "2024-12-16T09:08:25Z",
				},
				{
					CreatedAt: "2024-08-11T11:25:00Z",
				},
			},
		},
		{
			key: "with start time",
			input: input{
				id:        getUserID(postgres, 3),
				startTime: "2024-04-01T00:00:00Z",
			},
			expected: []task{
				{
					CreatedAt: "2024-12-16T09:08:25Z",
				},
				{
					CreatedAt: "2024-08-11T11:25:00Z",
				},
				{
					CreatedAt:   "2024-04-16T09:08:25Z",
					FinishedAt:  "2024-05-16T09:08:25Z",
					SummaryTime: "720h0m",
				},
			},
		},
		{
			key: "with end time",
			input: input{
				id:      getUserID(postgres, 3),
				endTime: "2024-05-12T00:00:00Z",
			},
			expected: []task{
				{
					CreatedAt:   "2024-03-11T11:25:00Z",
					FinishedAt:  "2024-05-11T09:08:25Z",
					SummaryTime: "1461h43m",
				},
				{
					CreatedAt:   "2024-01-16T09:08:25Z",
					FinishedAt:  "2024-01-16T16:10:00Z",
					SummaryTime: "7h2m",
				},
			},
		},
		{
			key: "with start and end time",
			input: input{
				id:        getUserID(postgres, 3),
				startTime: "2024-02-01T00:00:00Z",
				endTime:   "2024-05-12T00:00:00Z",
			},
			expected: []task{
				{
					CreatedAt:   "2024-03-11T11:25:00Z",
					FinishedAt:  "2024-05-11T09:08:25Z",
					SummaryTime: "1461h43m",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			query := url.Values{}

			query.Set("startTime", tc.input.startTime)
			query.Set("endTime", tc.input.endTime)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/summary-time/%s?%s", tc.input.id, query.Encode()), nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code, tc.key)

			tasks := make([]task, 0)
			json.NewDecoder(recorder.Body).Decode(&tasks)

			assert.Equal(t, tc.expected, tasks, tc.key)
		})
	}
}

func TestTaskSummaryTimeNegative(t *testing.T) {
	postgres, handler := prepare()
	t.Cleanup(func() {
		postgres.Close()
	})

	type input struct {
		id        string
		startTime string
		endTime   string
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
		{
			key: "miss start time",
			input: input{
				id:        getUserID(postgres, 3),
				startTime: "2025-02-01T00:00:00Z",
			},
			expectedCode: http.StatusNoContent,
		},
		{
			key: "miss end time",
			input: input{
				id:      getUserID(postgres, 3),
				endTime: "2023-02-01T00:00:00Z",
			},
			expectedCode: http.StatusNoContent,
		},
		{
			key: "miss start and end time",
			input: input{
				id:        getUserID(postgres, 3),
				startTime: "2022-02-01T00:00:00Z",
				endTime:   "2023-02-01T00:00:00Z",
			},
			expectedCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			query := url.Values{}

			query.Set("startTime", tc.input.startTime)
			query.Set("endTime", tc.input.endTime)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/tasks/summary-time/%s?%s", tc.input.id, query.Encode()), nil)
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedCode, recorder.Code, tc.key)
		})
	}
}
