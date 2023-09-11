package failure

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func NewResponse(message string, details any) *Response {
	return &Response{
		Message: message,
		Details: details,
	}
}

func Emit(
	w http.ResponseWriter,
	status int,
	message string,
	details any,
) {
	response := NewResponse(message, details)
	res, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	w.Write(res)
}
