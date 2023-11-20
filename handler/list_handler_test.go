package handler

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"noda"
	"noda/data/model"
	"noda/data/transfer"
	"noda/data/types"
	"strconv"
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
	var args = o.Called(ownerID, groupID, pagination, needle, sortBy)
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
	return args.Error(0)
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
	var args = o.Called(ownerID, groupID, listID, up)
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

func TestListHandler_HandleScatteredListRetrievalByID(t *testing.T) {
	var listID = uuid.New()
	const (
		method = "GET"
		target = "/me/lists/{list_id}"
	)

	t.Run("success", func(t *testing.T) {
		var (
			list = &model.List{
				ID:          listID,
				OwnerID:     userID,
				GroupID:     uuid.Nil,
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
		withPathParameters(&request, parameters{"list_id": list.ID.String()})
		var s = new(listServiceMock)
		s.On("FindListByID", list.OwnerID, uuid.Nil, list.ID).Return(list, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Empty(t, result.Header, "No header is expected, but got: %d.", len(result.Header))
		assert.Empty(t, result.Cookies(), "No cookie is expected, but got: %d.", len(result.Cookies()))
	})

	t.Run("got an expected service error", func(t *testing.T) {
		var expectedError = noda.ErrUserNotFound
		var expectedStatusCode = expectedError.Status()
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var s = new(listServiceMock)
		s.On("FindListByID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedError.Details())
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var expectedStatusCode = http.StatusInternalServerError
		var unexpected = errors.New("unexpected error")
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var s = new(listServiceMock)
		s.On("FindListByID", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListRetrievalByID(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandleGroupedListsRetrieval(t *testing.T) {
	var groupID = uuid.New()
	const (
		method = "GET"
		target = "/me/groups/{group_id}/lists"
	)

	t.Run("success", func(t *testing.T) {
		var (
			pagination = types.Pagination{Page: 1, RPP: 10}
			search     = "a"
			sortExpr   = "+name"
			values     = url.Values{
				"search":  []string{search},
				"sort_by": []string{sortExpr},
				"page":    []string{strconv.FormatInt(pagination.Page, 10)},
				"rpp":     []string{strconv.FormatInt(pagination.RPP, 10)},
			}
			serviceResult = &types.Result[model.List]{
				Page:      pagination.Page,
				RPP:       pagination.RPP,
				Payload:   make([]*model.List, 2),
				Retrieved: 1,
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = string(marshal(t, serviceResult))
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var s = new(listServiceMock)
		s.On("FindGroupedLists", userID, groupID, &pagination, search, sortExpr).
			Return(serviceResult, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListsRetrieval(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedResponseBody, string(responseBody))
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, result.Header, "No header is expected, but got: %d.", len(result.Header))
		assert.Empty(t, result.Cookies(), "No cookie is expected, but got: %d.", len(result.Cookies()))
	})

	t.Run("could not parse pagination: negative number", func(t *testing.T) {
		var (
			values               = url.Values{"page": []string{"-100"}}
			expectedStatusCode   = http.StatusBadRequest
			expectedResponseBody = "The parameter \\\"page\\\" must be a positive number."
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindGroupedLists")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListsRetrieval(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedResponseBody)
	})

	t.Run("could not parse sort expression", func(t *testing.T) {
		var (
			values               = url.Values{"sort_by": []string{"foo"}}
			expectedStatusCode   = http.StatusBadRequest
			expectedResponseBody = "[\"Must start with either one plus sign (+) or one minus sign (-).\",\"Must contain one or more word characters (alphanumeric characters and underscores).\"]"
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindGroupedLists")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListsRetrieval(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedResponseBody)
	})

	t.Run("parsing \"group_id\" failed: UUID is too short", func(t *testing.T) {
		var (
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID length."
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindGroupedLists")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListsRetrieval(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("got an expected service error", func(t *testing.T) {
		var (
			expectedError      = noda.ErrUserNotFound
			expectedStatusCode = expectedError.Status()
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var s = new(listServiceMock)
		s.On("FindGroupedLists", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListsRetrieval(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedError.Details())
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			unexpected         = errors.New("unexpected error")
			expectedStatusCode = http.StatusInternalServerError
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String()})
		var s = new(listServiceMock)
		s.On("FindGroupedLists", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListsRetrieval(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandleRetrievalOfLists(t *testing.T) {
	const (
		method = "GET"
		target = "/me/lists"
	)

	t.Run("success for scattered lists, where all=anything", func(t *testing.T) {
		var (
			pagination = types.Pagination{Page: 1, RPP: 10}
			search     = "a"
			sortExpr   = "+name"
			values     = url.Values{
				"search":  []string{search},
				"sort_by": []string{sortExpr},
				"page":    []string{strconv.FormatInt(pagination.Page, 10)},
				"rpp":     []string{strconv.FormatInt(pagination.RPP, 10)},
			}
			serviceResult = &types.Result[model.List]{
				Page:      pagination.Page,
				RPP:       pagination.RPP,
				Payload:   make([]*model.List, 2),
				Retrieved: 1,
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = string(marshal(t, serviceResult))
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.On("FindScatteredLists", userID, &pagination, search, sortExpr).
			Return(serviceResult, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleRetrievalOfLists(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedResponseBody, string(responseBody))
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, result.Header, "No header is expected, but got: %d.", len(result.Header))
		assert.Empty(t, result.Cookies(), "No cookie is expected, but got: %d.", len(result.Cookies()))
	})

	t.Run("success for all lists (scattered and grouped), where all=true", func(t *testing.T) {
		var (
			pagination = types.Pagination{Page: 1, RPP: 10}
			search     = "a"
			sortExpr   = "+name"
			values     = url.Values{
				"search":  []string{search},
				"sort_by": []string{sortExpr},
				"page":    []string{strconv.FormatInt(pagination.Page, 10)},
				"rpp":     []string{strconv.FormatInt(pagination.RPP, 10)},
				"all":     []string{"true"},
			}
			serviceResult = &types.Result[model.List]{
				Page:      pagination.Page,
				RPP:       pagination.RPP,
				Payload:   make([]*model.List, 2),
				Retrieved: 1,
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = string(marshal(t, serviceResult))
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.On("FindLists", userID, &pagination, search, sortExpr).
			Return(serviceResult, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleRetrievalOfLists(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedResponseBody, string(responseBody))
		assert.Equal(t, expectedStatusCode, result.StatusCode)
	})

	t.Run("could not parse pagination: negative number", func(t *testing.T) {
		var (
			values               = url.Values{"page": []string{"-100"}}
			expectedStatusCode   = http.StatusBadRequest
			expectedResponseBody = "The parameter \\\"page\\\" must be a positive number."
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindScatteredLists")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleRetrievalOfLists(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedResponseBody)
	})

	t.Run("could not parse sort expression", func(t *testing.T) {
		var (
			values               = url.Values{"sort_by": []string{"foo"}}
			expectedStatusCode   = http.StatusBadRequest
			expectedResponseBody = "[\"Must start with either one plus sign (+) or one minus sign (-).\",\"Must contain one or more word characters (alphanumeric characters and underscores).\"]"
		)
		var request = httptest.NewRequest(method, target+"?"+values.Encode(), nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.AssertNotCalled(t, "FindScatteredLists")
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleRetrievalOfLists(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedResponseBody)
	})

	t.Run("got an expected service error", func(t *testing.T) {
		var (
			expectedError      = noda.ErrUserNotFound
			expectedStatusCode = expectedError.Status()
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.On("FindScatteredLists", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleRetrievalOfLists(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedError.Details())
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			unexpected         = errors.New("unexpected error")
			expectedStatusCode = http.StatusInternalServerError
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		var s = new(listServiceMock)
		s.On("FindScatteredLists", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleRetrievalOfLists(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandlePartialUpdateOfListByID(t *testing.T) {
	var groupID, listID = uuid.New(), uuid.New()
	const (
		method        = "PATCH"
		target        = "/me/groups/{group_id}/lists/{list_id}"
		serviceMethod = "UpdateList"
	)

	t.Run("success", func(t *testing.T) {
		var (
			up                 = &transfer.ListUpdate{Name: "new list name", Description: "new list description"}
			expectedStatusCode = http.StatusNoContent
			requestBody        = marshal(t, up)
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, userID, groupID, listID, up).Return(true, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandlePartialUpdateOfGroupedList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = string(extractResponseBody(t, response.Body))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, responseBody, "No response body is expected.")
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
	})

	t.Run("body = {}? take me to the already existent list", func(t *testing.T) {
		var (
			requestBody        = []byte(" { } ")
			expectedStatusCode = http.StatusSeeOther
			recorder           = httptest.NewRecorder()
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.AssertNotCalled(t, serviceMethod)
		NewListHandler(s).HandlePartialUpdateOfGroupedList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = string(extractResponseBody(t, response.Body))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, responseBody, "No response body is expected.")
	})

	t.Run("could not decode JSON body: il-formed JSON", func(t *testing.T) {
		var (
			requestBody            = []byte("{")
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Body contains ill-formed JSON."
			request                = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
			recorder               = httptest.NewRecorder()
			s                      = new(listServiceMock)
		)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		s.AssertNotCalled(t, serviceMethod)
		NewListHandler(s).HandlePartialUpdateOfGroupedList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = string(extractResponseBody(t, response.Body))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, responseBody, expectedInResponseBody)
	})

	t.Run("parsing \"group_id\" failed: UUID is too short", func(t *testing.T) {
		var (
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID length."
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": "", "list_id": listID.String()})
		var m = new(listServiceMock)
		m.AssertNotCalled(t, serviceMethod)
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandlePartialUpdateOfGroupedList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("parsing \"list_id\" failed: UUID is too short", func(t *testing.T) {
		var (
			expectedStatusCode     = http.StatusBadRequest
			expectedInResponseBody = "Invalid UUID length."
		)
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": "x"})
		var m = new(listServiceMock)
		m.AssertNotCalled(t, serviceMethod)
		var recorder = httptest.NewRecorder()
		NewListHandler(m).HandlePartialUpdateOfGroupedList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
	})

	t.Run("got an expected service error", func(t *testing.T) {
		var (
			requestBody        = marshal(t, JSON{"name": "n", "description": "d"})
			expectedError      = noda.ErrUserNotFound
			expectedStatusCode = expectedError.Status()
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandlePartialUpdateOfGroupedList(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedError.Details())
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			requestBody        = marshal(t, JSON{"name": "n", "description": "d"})
			unexpected         = errors.New("unexpected error")
			expectedStatusCode = http.StatusInternalServerError
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(false, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandlePartialUpdateOfGroupedList(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandlePartialUpdateOfScatteredList(t *testing.T) {
	var listID = uuid.New()
	const (
		method        = "PATCH"
		target        = "/me/lists/{list_id}"
		serviceMethod = "UpdateList"
	)

	t.Run("success", func(t *testing.T) {
		var (
			up                 = &transfer.ListUpdate{Name: "new list name", Description: "new list description"}
			expectedStatusCode = http.StatusNoContent
			requestBody        = marshal(t, up)
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, userID, uuid.Nil, listID, up).Return(true, nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandlePartialUpdateOfScatteredList(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = string(extractResponseBody(t, response.Body))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, responseBody, "No response body is expected.")
		assert.Empty(t, response.Cookies(), "No cookie is expected, but got: %d.", len(response.Cookies()))
		assert.Empty(t, response.Header, "No header is expected, but got: %d.", len(response.Header))
	})

	t.Run("got an expected service error", func(t *testing.T) {
		var (
			requestBody        = marshal(t, JSON{"name": "n", "description": "d"})
			expectedError      = noda.ErrUserNotFound
			expectedStatusCode = expectedError.Status()
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, userID, uuid.Nil, listID, mock.AnythingOfType("*transfer.ListUpdate")).
			Return(false, expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandlePartialUpdateOfScatteredList(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Contains(t, string(responseBody), expectedError.Details())
	})

	t.Run("got an unexpected service error", func(t *testing.T) {
		var (
			requestBody        = marshal(t, JSON{"name": "n", "description": "d"})
			unexpected         = errors.New("unexpected error")
			expectedStatusCode = http.StatusInternalServerError
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, userID, uuid.Nil, listID, mock.AnythingOfType("*transfer.ListUpdate")).
			Return(false, unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandlePartialUpdateOfScatteredList(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandleGroupedListDeletion(t *testing.T) {
	var listID, groupID = uuid.New(), uuid.New()
	const (
		method        = "DELETE"
		target        = "/me/groups/{group_id}/lists/{list_id}"
		serviceMethod = "DeleteList"
	)

	t.Run("success", func(t *testing.T) {
		var expectedStatusCode = http.StatusNoContent
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, userID, groupID, listID).Return(nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListDeletion(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, responseBody)
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
		s.AssertNotCalled(t, serviceMethod)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListDeletion(recorder, request)
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
		s.AssertNotCalled(t, serviceMethod)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListDeletion(recorder, request)
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
		s.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything).
			Return(expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListDeletion(recorder, request)
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
		s.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything).
			Return(unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleGroupedListDeletion(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestListHandler_HandleScatteredListDeletion(t *testing.T) {
	var listID, groupID = uuid.New(), uuid.Nil
	const (
		method        = "DELETE"
		target        = "/me/lists/{list_id}"
		serviceMethod = "DeleteList"
	)

	t.Run("success", func(t *testing.T) {
		var expectedStatusCode = http.StatusNoContent
		var request = httptest.NewRequest(method, target, nil)
		withLoggedUser(&request)
		withPathParameters(&request, parameters{"group_id": groupID.String(), "list_id": listID.String()})
		var s = new(listServiceMock)
		s.On(serviceMethod, userID, groupID, listID).Return(nil)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListDeletion(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, responseBody)
		assert.Empty(t, result.Header, "No header is expected, but got: %d.", len(result.Header))
		assert.Empty(t, result.Cookies(), "No cookie is expected, but got: %d.", len(result.Cookies()))
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
		s.AssertNotCalled(t, serviceMethod)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListDeletion(recorder, request)
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
		s.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything).
			Return(expectedError)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListDeletion(recorder, request)
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
		s.On(serviceMethod, mock.Anything, mock.Anything, mock.Anything).
			Return(unexpected)
		var recorder = httptest.NewRecorder()
		NewListHandler(s).HandleScatteredListDeletion(recorder, request)
		var result = recorder.Result()
		defer result.Body.Close()
		var responseBody = extractResponseBody(t, result.Body)
		assert.Equal(t, expectedStatusCode, result.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}
