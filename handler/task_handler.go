package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"noda/failure"

	"github.com/google/uuid"

	"noda/data/transfer"
	"noda/service"
)

type TaskHandler struct {
	s service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{s: service}
}

func (h *TaskHandler) doCreateTask(belongsToAList bool, w http.ResponseWriter, r *http.Request) {
	var task = new(transfer.TaskCreation)
	var err = parseRequestBody(w, r, task)
	if nil != err {
		failure.EmitError(w, failure.ErrMalformedRequest.Clone().SetDetails(err.Error()))
		return
	}
	err = task.Validate()
	if nil != err {
		failure.EmitError(w, failure.ErrBadRequest.Clone().SetDetails(err.Error()))
		return
	}
	var userID, _ = extractUserPayload(r)
	var listID uuid.UUID
	if belongsToAList {
		listID = parseParameterToUUID(w, r, "list_id")
		if didNotParse(listID) {
			return
		}
	}
	insertedTaskID, err := h.s.Save(userID, listID, task)
	if gotAndHandledServiceError(w, err) {
		return
	}
	var result = map[string]string{"inserted_id": insertedTaskID.String()}
	data, err := json.Marshal(result)
	if nil != err {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}

func (h *TaskHandler) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	h.doCreateTask(true, w, r)
}

func (h *TaskHandler) HandleCreateTaskForTodayList(w http.ResponseWriter, r *http.Request) {
	h.doCreateTask(false, w, r)
}
