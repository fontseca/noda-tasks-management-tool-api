package failure

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/http"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrGroupNotFound     = errors.New("group not found")
	ErrListNotFound      = errors.New("list not found")
	ErrSettingNotFound   = errors.New("user setting not found")
	ErrSameEmail         = errors.New("the given email address is already registered")
	ErrIncorrectPassword = errors.New("the given password does not match with stored password")
	ErrPasswordTooLong   = errors.New("the given password length exceeds 72 bytes")
	ErrUserBlocked       = errors.New("this user has been blocked")
	ErrDeadlineExceeded  = errors.New("context deadline exceeded")
)

type Aggregation struct {
	errors []string
}

func NewAggregation() *Aggregation {
	return &Aggregation{}
}

func (a *Aggregation) Error() string {
	str := ""
	for _, err := range a.errors {
		str += err + "\n"
	}
	return str
}

func (a *Aggregation) Append(err error) {
	a.errors = append(a.errors, err.Error())
}

func (a *Aggregation) Dump() []string {
	return a.errors
}

func (a *Aggregation) Has() bool {
	return len(a.errors) > 0
}

func PQErrorToString(err *pq.Error) string {
	if err.Hint == "" {
		err.Hint = "(none)"
	}
	if err.Detail == "" {
		err.Detail = "(none)"
	}
	return fmt.Sprintf("postgres driver failed with error \033[1;31m%s\033[0m (%s):\n"+
		"  message: \033[0;33m%s\033[0m\n"+
		"   detail: %s\n"+
		"     hint: %s",
		err.Code, err.Code.Name(), err.Message, err.Detail, err.Hint)
}

type Response struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func NewResponse(message string, details any) *Response {
	err, ok := details.(error)
	if ok {
		return &Response{
			Message: message,
			Details: err.Error(),
		}
	}

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
