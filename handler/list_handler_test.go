package handler

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"testing"
)

type mockListService struct {
	mock.Mock
}

func (o *mockListService) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	var args = o.Called(ownerID, groupID, next)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) FindListByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	var args = o.Called(ownerID, groupID, listID)
	var arg1 = args.Get(0)
	if nil != arg1 {
		list = arg1.(*model.List)
	}
	return list, args.Error(1)
}

func (o *mockListService) GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) FindLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (lists *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		lists = arg1.(*types.Result[model.List])
	}
	return lists, args.Error(1)
}

func (o *mockListService) FindGroupedLists(ownerID, groupID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *mockListService) FindScatteredLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *mockListService) DeleteList(ownerID, groupID, listID uuid.UUID) error {
	var args = o.Called(ownerID, groupID, listID)
	return args.Error(1)
}

func (o *mockListService) DuplicateList(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error) {
	var args = o.Called(ownerID, listID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *mockListService) ConvertToScatteredList(ownerID, listID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *mockListService) MoveList(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, targetGroupID)
	return args.Bool(0), args.Error(1)
}

func (o *mockListService) UpdateList(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error) {
	var args = o.Called(ownerID, groupID, listID)
	return args.Bool(0), args.Error(1)
}

func TestListHandler_HandleGroupedListCreation(t *testing.T) {
	var groupID = uuid.New()
	const (
		method = "POST"
		target = "/me/groups/{group_id}/lists"
	)

	t.Run("success", func(t *testing.T) {
		var (
			requestBody          = marshal(t, JSON{"name": "list name", "description": "list description"})
			insertedID           = uuid.New()
			expectedStatusCode   = http.StatusCreated
			expectedResponseBody = marshal(t, JSON{"insertedID": insertedID.String()})
			next                 = &transfer.ListCreation{Name: "list name", Description: "list description"}
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withUserLoggedUser(&request)
		withPathParameter(&request, "group_id", groupID.String())
		var m = new(mockListService)
		m.On("SaveList", userID, groupID, next).Return(insertedID, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleGroupedListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
	})

	t.Run("could not decode JSON body: il-formed JSON", func(t *testing.T) {
		var (
			requestBody            = []byte("{")
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Body contains ill-formed JSON."
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var m = new(mockListService)
		m.AssertNotCalled(t, "SaveList")
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleGroupedListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("list creation validation failed on required fields", func(t *testing.T) {
		var (
			requestBody            = []byte("{}")
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "[\"Validation for \\\"name\\\" failed on: required.\",\"Validation for \\\"description\\\" failed on: required.\"]"
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var m = new(mockListService)
		m.AssertNotCalled(t, "SaveList")
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleGroupedListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("parsing \"group_id\" failed: UUID is too short", func(t *testing.T) {
		var (
			requestBody            = marshal(t, JSON{"name": "n", "description": "d"})
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID length."
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withUserLoggedUser(&request)
		withPathParameter(&request, "group_id", "x")
		var m = new(mockListService)
		m.AssertNotCalled(t, "SaveList")
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleGroupedListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("parsing \"group_id\" failed: invalid UUID format", func(t *testing.T) {
		var (
			requestBody            = marshal(t, JSON{"name": "n", "description": "d"})
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID format."
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withUserLoggedUser(&request)
		withPathParameter(&request, "group_id", "a0e2240b-8f5b-4b1e-88a9-c6d9284a6afX")
		var m = new(mockListService)
		m.AssertNotCalled(t, "SaveList")
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleGroupedListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("got a service error", func(t *testing.T) {
		var (
			requestBody        = marshal(t, JSON{"name": "n", "description": "d"})
			expectedStatusCode = http.StatusInternalServerError
			unexpected         = errors.New("unexpected error")
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withUserLoggedUser(&request)
		withPathParameter(&request, "group_id", groupID.String())
		var m = new(mockListService)
		m.On("SaveList", mock.Anything, mock.Anything, mock.AnythingOfType("*transfer.ListCreation")).
			Return(uuid.Nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleGroupedListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}
