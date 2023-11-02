package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"noda/api/data/transfer"
	"noda/api/service"
	"noda/failure"
	"strings"
)

type GroupHandler struct {
	s *service.GroupService
}

func NewGroupHandler(service *service.GroupService) *GroupHandler {
	return &GroupHandler{service}
}

func (h *GroupHandler) HandleGroupCreation(w http.ResponseWriter, r *http.Request) {
	var group = new(transfer.GroupCreation)
	var err = decodeJSONRequestBody(w, r, group)
	if nil != err {
		var mr *malformedRequest
		if errors.As(err, &mr) {
			failure.Emit(w, mr.status, mr.message, mr.details)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	var validationError = group.Validate()
	if nil != validationError {
		failure.Emit(w, http.StatusBadRequest, "validation did not succeed", validationError.Dump())
		return
	}
	userID, _ := extractUserPayload(r)
	insertedID, err := h.s.SaveGroup(userID, group)
	if nil != err {
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, failure.ErrNotFound):
			failure.Emit(w, http.StatusNotFound, "not found", "this is user account no longer exists")
		case errors.Is(err, failure.ErrDeadlineExceeded):
			w.WriteHeader(http.StatusInternalServerError)
		case strings.Contains(err.Error(), "name too long"):
			failure.Emit(w, http.StatusBadRequest, "bad request", err)
		}
		return
	}
	var result = map[string]string{"insertedID": insertedID}
	data, err := json.Marshal(result)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
