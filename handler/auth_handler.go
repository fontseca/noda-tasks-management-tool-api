package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"noda/data/transfer"
	"noda/failure"
	"noda/service"
)

type AuthenticationHandler struct {
	s service.AuthenticationService
}

func NewAuthenticationHandler(s service.AuthenticationService) *AuthenticationHandler {
	return &AuthenticationHandler{s}
}

func (h *AuthenticationHandler) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	next := &transfer.UserCreation{}
	var err = parseRequestBody(w, r, next)
	if nil != err {
		failure.EmitError(w, failure.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	err = next.Validate()
	if nil != err {
		failure.EmitError(w, failure.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	insertedID, err := h.s.SignUp(next)
	if err != nil {
		var (
			a *failure.AggregateDetails
			e *failure.Error
		)
		if errors.As(err, &a) {
			failure.EmitError(w, failure.ErrPasswordRestrictions.Clone().SetDetails(a.Error()))
		} else if errors.As(err, &e) {
			failure.EmitError(w, e)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	var payload = map[string]string{"user_uuid": insertedID.String()}
	data, err := json.Marshal(payload)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *AuthenticationHandler) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	credentials := &transfer.UserCredentials{}
	var err = parseRequestBody(w, r, credentials)
	if nil != err {
		failure.EmitError(w, failure.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	res, err := h.s.SignIn(credentials)
	if err != nil {
		var e *failure.Error
		if errors.As(err, &e) {
			switch {
			default:
				failure.EmitError(w, e)
			case errors.Is(e, failure.ErrUserNotFound):
				failure.EmitError(w, e.
					Clone().
					SetDetails("Could not find any user with the email %q.").
					FormatDetails(credentials.Email).
					SetHint("Verify email address or use another one."))
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	data, err := json.Marshal(res)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
