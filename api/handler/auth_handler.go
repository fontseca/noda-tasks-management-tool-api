package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"noda/api/data/transfer"
	"noda/api/service"
	"noda/failure"
)

type AuthenticationHandler struct {
	s *service.AuthenticationService
}

func NewAuthenticationHandler(s *service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{s}
}

func (h *AuthenticationHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	next := &transfer.UserCreation{}
	if err := json.NewDecoder(r.Body).Decode(next); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := next.Validate(); err != nil {
		failure.Emit(w, http.StatusBadRequest, "validations did not succeed", err.Dump())
		return
	}

	insertedID, err := h.s.SignUp(next)
	if err != nil {
		var passwdErrors *failure.Aggregation
		switch {
		case errors.As(err, &passwdErrors):
			failure.Emit(w, http.StatusBadRequest, "password restrictions not met", passwdErrors.Dump())
			return
		case errors.Is(err, failure.ErrSameEmail):
			failure.Emit(w, http.StatusBadRequest,
				"signing up failed", failure.ErrSameEmail)
			return
		case errors.Is(err, failure.ErrPasswordTooLong):
			failure.Emit(w, http.StatusBadRequest,
				"signing up failed", failure.ErrPasswordTooLong)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": insertedID.String(),
	})
}

func (h *AuthenticationHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	credentials := &transfer.UserCredentials{}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(credentials); err != nil {
		// TODO: Catch all different errors
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := h.s.SignIn(credentials)
	if err != nil {
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		case errors.Is(err, failure.ErrUserNotFound):
			failure.Emit(w, http.StatusNotFound,
				"signing in failed", fmt.Sprintf("could not find any user with email %q", credentials.Email))
			return
		case errors.Is(err, failure.ErrIncorrectPassword):
			failure.Emit(w, http.StatusBadRequest,
				"signing in failed", failure.ErrIncorrectPassword)
			return
		case errors.Is(err, failure.ErrUserBlocked):
			failure.Emit(w, http.StatusForbidden,
				"authentication refused", "this account has been blocked")
			return
		}
	}

	json.NewEncoder(w).Encode(res)
}
