package handler

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"testing"
	"time"
)

type listServiceMock struct {
	mock.Mock
}

func (o *listServiceMock) SaveList(ownerID, groupID uuid.UUID, next *transfer.ListCreation) (insertedID uuid.UUID, err error) {
	var args = o.Called(ownerID, groupID, next)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *listServiceMock) FindListByID(ownerID, groupID, listID uuid.UUID) (list *model.List, err error) {
	var args = o.Called(ownerID, groupID, listID)
	var arg1 = args.Get(0)
	if nil != arg1 {
		list = arg1.(*model.List)
	}
	return list, args.Error(1)
}

func (o *listServiceMock) GetTodayListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *listServiceMock) GetTomorrowListID(ownerID uuid.UUID) (listID uuid.UUID, err error) {
	var args = o.Called(ownerID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *listServiceMock) FindLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (lists *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		lists = arg1.(*types.Result[model.List])
	}
	return lists, args.Error(1)
}

func (o *listServiceMock) FindGroupedLists(ownerID, groupID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *listServiceMock) FindScatteredLists(ownerID uuid.UUID, pagination *types.Pagination, needle, sortBy string) (result *types.Result[model.List], err error) {
	var args = o.Called(ownerID, pagination, needle, sortBy)
	var arg1 = args.Get(0)
	if nil != arg1 {
		result = arg1.(*types.Result[model.List])
	}
	return result, args.Error(1)
}

func (o *listServiceMock) DeleteList(ownerID, groupID, listID uuid.UUID) error {
	var args = o.Called(ownerID, groupID, listID)
	return args.Error(1)
}

func (o *listServiceMock) DuplicateList(ownerID, listID uuid.UUID) (replicaID uuid.UUID, err error) {
	var args = o.Called(ownerID, listID)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (o *listServiceMock) ConvertToScatteredList(ownerID, listID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID)
	return args.Bool(0), args.Error(1)
}

func (o *listServiceMock) MoveList(ownerID, listID, targetGroupID uuid.UUID) (ok bool, err error) {
	var args = o.Called(ownerID, listID, targetGroupID)
	return args.Bool(0), args.Error(1)
}

func (o *listServiceMock) UpdateList(ownerID, groupID, listID uuid.UUID, up *transfer.ListUpdate) (ok bool, err error) {
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
			next                 = &transfer.ListCreation{Name: "list name", Description: "list description"}
			requestBody          = marshal(t, JSON{"name": next.Name, "description": next.Description})
			insertedID           = uuid.New()
			expectedStatusCode   = http.StatusCreated
			expectedResponseBody = marshal(t, JSON{"insertedID": insertedID.String()})
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var m = new(listServiceMock)
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
		var m = new(listServiceMock)
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
		var m = new(listServiceMock)
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
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": "x"})
		var m = new(listServiceMock)
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
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": "a0e2240b-8f5b-4b1e-88a9-c6d9284a6afX"})
		var m = new(listServiceMock)
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
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var m = new(listServiceMock)
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

func TestListHandler_HandleScatteredListCreation(t *testing.T) {
	const (
		method = "POST"
		target = "/me/lists"
	)
	var next = &transfer.ListCreation{Name: "list name", Description: "list description"}

	t.Run("success", func(t *testing.T) {
		var (
			requestBody          = marshal(t, JSON{"name": next.Name, "description": next.Description})
			insertedID           = uuid.New()
			expectedStatusCode   = http.StatusCreated
			expectedResponseBody = marshal(t, JSON{"insertedID": insertedID.String()})
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		var m = new(listServiceMock)
		m.On("SaveList", userID, uuid.Nil, next).Return(insertedID, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleScatteredListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
	})

	t.Run("got a service error", func(t *testing.T) {
		var (
			requestBody        = marshal(t, JSON{"name": next.Name, "description": next.Description})
			expectedStatusCode = http.StatusInternalServerError
			unexpected         = errors.New("unexpected error")
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		var m = new(listServiceMock)
		m.On("SaveList", userID, uuid.Nil, next).Return(uuid.Nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandleScatteredListCreation(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandleGroupedListRetrievalByID(t *testing.T) {
	var listID, groupID = uuid.New(), uuid.New()
	const (
		method = "GET"
		target = "/me/groups/{group_id}/lists/{list_id}"
	)

	t.Run("success", func(t *testing.T) {
		var (
			list = &model.List{
				ID:          listID,
				OwnerID:     userID,
				GroupID:     groupID,
				Name:        "list name",
				Description: "list description",
				UpdatedAt:   time.Now(),
				CreatedAt:   time.Now(),
				ArchivedAt:  nil,
				IsArchived:  false,
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = marshal(t, list)
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": list.GroupID.String(), "list_id": list.ID.String()})
		var s = new(listServiceMock)
		s.On("FindListByID", list.OwnerID, list.GroupID, list.ID).Return(list, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Empty(t, result.Header, "No header is expected, but got: %d.", len(result.Header))
		assert.Empty(t, result.Cookies(), "No cookie is expected, but got: %d.", len(result.Cookies()))
	})

	t.Run("parsing \"group_id\" failed: UUID is too short", func(t *testing.T) {
		var (
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID length."
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": "x"})
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindListByID")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("parsing \"list_id\" failed: UUID is too short", func(t *testing.T) {
		var (
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID length."
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": "x"})
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindListByID")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("got an expected service error", func(t *testing.T) {
		var (
			expectedError          = noda.ErrUserNotFound
			expectedStatusCode     = expectedError.Status()
			expectedInResponseBody = expectedError.Details()
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On("FindListByID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			expectedStatusCode = http.StatusInternalServerError
			unexpected         = errors.New("unexpected error")
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On("FindListByID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}
