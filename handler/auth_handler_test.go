package handler

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"noda/data/transfer"
	"noda/data/types"
	"noda/mocks"
	"testing"
	"time"
)

func TestAuthenticationHandler_HandleSignUp(t *testing.T) {
	const (
		method  = "POST"
		target  = "/signup"
		routine = "SignUp"
	)
	var creation = &transfer.UserCreation{
		FirstName:  "John",
		MiddleName: "Alexander",
		LastName:   "Doe",
		Surname:    "Doe",
		Password:   "C}Ryf6P.%'g@$D+;7A,(b",
		Email:      "es09911@zbock.com",
	}

	t.Run("success", func(t *testing.T) {
		var (
			inserted             = uuid.New()
			expectedStatusCode   = http.StatusCreated
			expectedResponseBody = marshal(t, JSON{"user_id": inserted.String()})
			requestBody          = marshal(t, creation)
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.On(routine, creation).Return(inserted, nil)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignUp(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	})

	t.Run("could not parse JSON body", func(t *testing.T) {
		var (
			expectedInResponseBody = "Body contains ill-formed JSON."
			expectedStatusCode     = http.StatusBadRequest
			requestBody            = []byte(" { ")
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.AssertNotCalled(t, routine)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignUp(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	})

	t.Run("could not validate JSON body", func(t *testing.T) {
		var (
			expectedInResponseBody = "[\"Validation for \\\"first_name\\\" failed on: required.\",\"Validation for \\\"last_name\\\" failed on: required.\",\"Validation for \\\"email\\\" failed on: required.\",\"Validation for \\\"password\\\" failed on: required.\"]"
			expectedStatusCode     = http.StatusBadRequest
			requestBody            = []byte(" { } ")
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.AssertNotCalled(t, routine)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignUp(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	})

	t.Run("got service error", func(t *testing.T) {
		var (
			expectedStatusCode = http.StatusInternalServerError
			unexpected         = errors.New("unexpected error")
			requestBody        = marshal(t, creation)
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.On(routine, creation).Return(uuid.Nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignUp(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}

func TestAuthenticationHandler_HandleSignIn(t *testing.T) {
	const (
		method  = "POST"
		target  = "/login"
		routine = "SignIn"
	)
	var credentials = &transfer.UserCredentials{
		Email:    "es09911@zbock.com",
		Password: "C}Ryf6P.%'g@$D+;7A,(b",
	}

	t.Run("success", func(t *testing.T) {
		var (
			tokenPayload = &types.TokenPayload{
				ID:       "",
				Issuer:   "issuer",
				Token:    "token",
				Subject:  "subject",
				IssuedAt: time.Now(),
				Expires: types.TokenExpires{
					Within: 3600,
					At:     time.Now().Add(1 * time.Hour),
					Unit:   "s",
				},
			}
			expectedStatusCode   = http.StatusOK
			expectedResponseBody = marshal(t, tokenPayload)
			requestBody          = marshal(t, credentials)
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.On(routine, credentials).Return(tokenPayload, nil)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignIn(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, string(expectedResponseBody), string(responseBody))
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	})

	t.Run("could not parse JSON body", func(t *testing.T) {
		var (
			expectedInResponseBody = "Body contains ill-formed JSON."
			expectedStatusCode     = http.StatusBadRequest
			requestBody            = []byte(" { ")
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.AssertNotCalled(t, routine)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignIn(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Contains(t, string(responseBody), expectedInResponseBody)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
	})

	t.Run("got service error", func(t *testing.T) {
		var (
			expectedStatusCode = http.StatusInternalServerError
			unexpected         = errors.New("unexpected error")
			requestBody        = marshal(t, credentials)
		)
		var request = httptest.NewRequest(method, target, bytes.NewReader(requestBody))
		var s = mocks.NewAuthenticationServiceMock()
		s.On(routine, credentials).Return(nil, unexpected)
		var recorder = httptest.NewRecorder()
		NewAuthenticationHandler(s).HandleSignIn(recorder, request)
		var response = recorder.Result()
		defer response.Body.Close()
		var responseBody = extractResponseBody(t, response.Body)
		assert.Equal(t, expectedStatusCode, response.StatusCode)
		assert.Empty(t, string(responseBody), "No response body is expected.")
	})
}
