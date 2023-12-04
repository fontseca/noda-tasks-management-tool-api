package handler

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"noda/data/transfer"
	"noda/data/types"
	"noda/mocks"
	"strconv"
	"testing"
)

func TestUserHandler_HandleUsersRetrieval(t *testing.T) {
	const (
		method  = "GET"
		target  = "/users"
		routine = "Fetch"
	)

	t.Run("success", func(t *testing.T) {
		var (
			pagination = types.Pagination{Page: 1, RPP: 10}
			search     = "x"
			sortExpr   = "+first_name"
			values     = url.Values{
				"search":  []string{search},
				"sort_by": []string{sortExpr},
				"page":    []string{strconv.FormatInt(pagination.Page, 10)},
				"rpp":     []string{strconv.FormatInt(pagination.RPP, 10)},
			}
			serviceResult = &types.Result[transfer.User]{
				Page:      pagination.Page,
				RPP:       pagination.RPP,
				Payload:   make([]*transfer.User, 2),
				Retrieved: 1,
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = string(marshal(t, serviceResult))
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		var s = mocks.NewUserServiceMock()
		s.On(routine, &pagination, search, sortExpr).Return(serviceResult, nil)
		var recorder = httptest.NewRecorder()
		NewUserHandler(s).HandleUsersRetrieval(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedResponseBody, string(responseBody))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			unexpected         = errors.New("unexpected error")
			expectedStatusCode = http.StatusInternalServerError
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		var s = mocks.NewUserServiceMock()
		s.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewUserHandler(s).HandleUsersRetrieval(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestUserHandler_HandleBlockedUsersRetrieval(t *testing.T) {
	const (
		method  = "GET"
		target  = "/users/blocked"
		routine = "FetchBlocked"
	)

	t.Run("success", func(t *testing.T) {
		var (
			pagination = types.Pagination{Page: 1, RPP: 10}
			search     = "x"
			sortExpr   = "+first_name"
			values     = url.Values{
				"search":  []string{search},
				"sort_by": []string{sortExpr},
				"page":    []string{strconv.FormatInt(pagination.Page, 10)},
				"rpp":     []string{strconv.FormatInt(pagination.RPP, 10)},
			}
			serviceResult = &types.Result[transfer.User]{
				Page:      pagination.Page,
				RPP:       pagination.RPP,
				Payload:   make([]*transfer.User, 2),
				Retrieved: 1,
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = string(marshal(t, serviceResult))
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		var s = mocks.NewUserServiceMock()
		s.On(routine, &pagination, search, sortExpr).Return(serviceResult, nil)
		var recorder = httptest.NewRecorder()
		NewUserHandler(s).HandleBlockedUsersRetrieval(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedResponseBody, string(responseBody))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			unexpected         = errors.New("unexpected error")
			expectedStatusCode = http.StatusInternalServerError
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		var s = mocks.NewUserServiceMock()
		s.On(routine, mock.Anything, mock.Anything, mock.Anything).Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewUserHandler(s).HandleBlockedUsersRetrieval(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}
