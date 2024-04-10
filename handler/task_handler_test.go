package handler

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"noda/data/transfer"
	"noda/data/types"
	"noda/mocks"
	"testing"
	"time"
)

func TestTaskHandler_HandleCreateTask(t *testing.T) {
	const (
		method        = "POST"
		target        = "/me/lists/{list_id}/tasks"
		serviceMethod = "Save"
	)
	var (
		listID   = uuid.New()
		creation = &transfer.TaskCreation{
			Title:       "Title",
			Headline:    "Headline",
			Description: "Description",
			Priority:    types.TaskPriorityHigh,
			Status:      types.TaskStatusIncomplete,
			DueDate:     time.Time{},
			RemindAt:    time.Time{},
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			insertedID           = uuid.New()
			requestBody          = marshal(t, creation)
			expectedStatusCode   = http.StatusCreated
			expectedResponseBody = marshal(t, JSON{"inserted_id": insertedID.String()})
		)
		var recorder = httptest.NewRecorder()
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var m = mocks.NewTaskServiceMock()
		m.On(serviceMethod, userID, listID, creation).Return(insertedID, nil)
		NewTaskHandler(m).HandleCreateTask(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
	})

	t.Run("got a service error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var requestBody = marshal(t, creation)
		var expectedStatusCode = http.StatusInternalServerError
		var recorder = httptest.NewRecorder()
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var m = mocks.NewTaskServiceMock()
		m.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything).Return(uuid.Nil, unexpected)
		NewTaskHandler(m).HandleCreateTask(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody))
	})
}

func TestTaskHandler_HandleCreateTaskForTodayList(t *testing.T) {
	const (
		method        = "POST"
		target        = "/me/tasks"
		serviceMethod = "Save"
	)
	var (
		creation = &transfer.TaskCreation{
			Title:       "Title",
			Headline:    "Headline",
			Description: "Description",
			Priority:    types.TaskPriorityHigh,
			Status:      types.TaskStatusIncomplete,
			DueDate:     time.Time{},
			RemindAt:    time.Time{},
		}
	)

	t.Run("success", func(t *testing.T) {
		var (
			insertedID           = uuid.New()
			requestBody          = marshal(t, creation)
			expectedStatusCode   = http.StatusCreated
			expectedResponseBody = marshal(t, JSON{"inserted_id": insertedID.String()})
		)
		var recorder = httptest.NewRecorder()
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		var m = mocks.NewTaskServiceMock()
		m.On(serviceMethod, userID, uuid.Nil, creation).Return(insertedID, nil)
		NewTaskHandler(m).HandleCreateTaskForTodayList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
	})

	t.Run("got a service error", func(t *testing.T) {
		var unexpected = errors.New("unexpected error")
		var requestBody = marshal(t, creation)
		var expectedStatusCode = http.StatusInternalServerError
		var recorder = httptest.NewRecorder()
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		var m = mocks.NewTaskServiceMock()
		m.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything).Return(uuid.Nil, unexpected)
		NewTaskHandler(m).HandleCreateTaskForTodayList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody))
	})
}
