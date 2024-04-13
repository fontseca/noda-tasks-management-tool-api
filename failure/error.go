package failure

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

/* URL details.  */

var (
	ErrTargetNotFound = &Error{
		code:    ErrorCode("U0001"),
		message: "Target not found.",
		details: "Could not find the resource requested by the given URL.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrBadQueryParameter = &Error{
		code:    ErrorCode("U0002"),
		message: "Query parameter failure.",
		details: "",
		hint:    "",
		status:  http.StatusBadRequest,
	}
	ErrMultipleValuesForQueryParameter = &Error{
		code:    ErrorCode("U0003"),
		message: "Multiple values for query parameter.",
		details: "Too much values for query parameter: %q.",
		hint:    "Provide only one value for this query parameter.",
		status:  http.StatusBadRequest,
	}
	ErrQueryParameterNotParsed = &Error{
		code:    ErrorCode("U0004"),
		message: "Could not parse parameter.",
		details: "Could not parse query parameter: %q.",
		hint:    "Provide only one value for this query parameter.",
		status:  http.StatusBadRequest,
	}
	ErrInvalidUUIDFormat = &Error{
		code:    ErrorCode("U0005"),
		message: "Error parsing path parameter.",
		details: "Invalid UUID format.",
		hint:    "",
		status:  http.StatusBadRequest,
	}
	ErrInvalidUUIDLength = &Error{
		code:    ErrorCode("U0006"),
		message: "Error parsing path parameter.",
		details: "Invalid UUID length.",
		hint:    "",
		status:  http.StatusBadRequest,
	}
)

/* Authentication details.  */

var (
	ErrMissingAuthorizationHeader = &Error{
		code:    ErrorCode("A0001"),
		message: "Authorization refused.",
		details: "Missing \"Authorization\" header in request.",
		hint:    "Check HTTP headers in request.",
		status:  http.StatusUnauthorized,
	}
	ErrNoEnoughRights = &Error{
		code:    ErrorCode("A0002"),
		message: "Authorization refused.",
		details: "Insufficient rights to access this resource.",
		hint:    "",
		status:  http.StatusUnauthorized,
	}
	ErrJSONWebToken = &Error{
		code:    ErrorCode("A0003"),
		message: "JSON Web Token failure.",
		details: "",
		hint:    "",
		status:  http.StatusUnauthorized,
	}
	ErrCorruptedClaim = &Error{
		code:    ErrorCode("A0004"),
		message: "JSON Web Token failure.",
		details: "One claim in JWT seems to be corrupted.",
		hint:    "",
		status:  http.StatusUnauthorized,
	}
)

/* Service details.  */

var (
	ErrTooLong = &Error{
		code:    ErrorCode("S0001"),
		message: "Request did not meet validation.",
		details: "Field %q is too long for %s. Maximum name length must be %d.",
		hint:    "",
		status:  http.StatusBadRequest,
	}
	ErrPasswordTooLong = &Error{
		code:    ErrorCode("S0002"),
		message: "Request did not meet validation.",
		details: "The length of this password exceeds 72 bytes.",
		hint:    "",
		status:  http.StatusBadRequest,
	}
)

/* Request details.  */

var (
	ErrMalformedRequest = &Error{
		code:    ErrorCode("RQ001"),
		message: "Bad JSON in request body.",
		details: "",
		hint:    "",
		status:  http.StatusBadRequest,
	}
	ErrBadRequest = &Error{
		code:    ErrorCode("RQ002"),
		message: "Bad request made.",
		details: "",
		hint:    "Check fields in request body object.",
		status:  http.StatusBadRequest,
	}
	ErrPasswordRestrictions = &Error{
		code:    ErrorCode("RQ003"),
		message: "Password restrictions not met.",
		details: "",
		hint:    "",
		status:  http.StatusBadRequest,
	}
	ErrSelfOperation = &Error{
		code:    ErrorCode("RQ004"),
		message: "Refused to perform self operation.",
		details: "You cannot perform this operation on the logged in user.",
		hint:    "",
		status:  http.StatusBadRequest,
	}
)

/* Repository details.  */

var (
	ErrUserNotFound = &Error{
		code:    ErrorCode("R0001"),
		message: "Not found.",
		details: "Could not find any user with this ID.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrUserNoLongerExists = &Error{
		code:    ErrorCode("R0008"),
		message: "Not found.",
		details: "This user account no longer exists.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrGroupNotFound = &Error{
		code:    ErrorCode("R002"),
		message: "Not found.",
		details: "Could not find any group with this ID.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrListNotFound = &Error{
		code:    ErrorCode("R0003"),
		message: "Not found.",
		details: "Could not find any list with this ID.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrTaskNotFound = &Error{
		code:    ErrorCode("R0008"),
		message: "Not found.",
		details: "Could not find any task with this ID.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrSettingNotFound = &Error{
		code:    ErrorCode("R0004"),
		message: "Not found.",
		details: "Could not find any user setting with this ID.",
		hint:    "",
		status:  http.StatusNotFound,
	}
	ErrSameEmail = &Error{
		code:    ErrorCode("R0005"),
		message: "Conflicting email address.",
		details: "This email address is already registered.",
		hint:    "Try using another one.",
		status:  http.StatusBadRequest,
	}
	ErrIncorrectPassword = &Error{
		code:    ErrorCode("R0006"),
		message: "Signing in failed.",
		details: "This password does not match with the one that's expected.",
		hint:    "Try using another one or recover it.",
		status:  http.StatusBadRequest,
	}
	ErrUserBlocked = &Error{
		code:    ErrorCode("R0007"),
		message: "Authentication refused.",
		details: "This user account has been blocked.",
		hint:    "",
		status:  http.StatusForbidden,
	}
	ErrDeadlineExceeded = errors.New("context deadline exceeded")
)

type ErrorCode string

type Error struct {
	status  int
	code    ErrorCode
	message string
	details string
	hint    string
}

func (e *Error) Error() string {
	return e.details
}

func (e *Error) Clone() *Error {
	return &Error{
		code:    e.code,
		message: e.message,
		details: e.details,
		hint:    e.hint,
		status:  e.status,
	}
}

func (e *Error) Details() string {
	return e.details
}

func (e *Error) SetDetails(details string) *Error {
	e.details = strings.Trim(details, " \n\t")
	return e
}

func (e *Error) FormatDetails(a ...any) *Error {
	e.details = fmt.Sprintf(strings.Trim(e.details, " \n\t"), a...)
	return e
}

func (e *Error) Status() int {
	return e.status
}

func (e *Error) SetStatus(status int) *Error {
	e.status = status
	return e
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) SetMessage(message string) *Error {
	e.message = strings.Trim(message, " \n\t")
	return e
}

func (e *Error) Hint() string {
	return e.hint
}

func (e *Error) SetHint(hint string) *Error {
	e.hint = strings.Trim(hint, " \n\t")
	return e
}

type AggregateDetails struct {
	details []string
}

func (a *AggregateDetails) Error() string {
	data, err := json.Marshal(a.details)
	if nil != err {
		log.Println(err)
		return ""
	}
	return string(data)
}

func (a *AggregateDetails) Append(detail string) {
	a.details = append(a.details, detail)
}

func (a *AggregateDetails) Has() bool {
	return len(a.details) > 0
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

type errorBody struct {
	Code    ErrorCode `json:"error_code"`
	Message string    `json:"message"`
	Details any       `json:"details,omitempty"`
	Hint    string    `json:"hint,omitempty"`
}

func EmitError(w http.ResponseWriter, e *Error) {
	var response = &errorBody{
		Code:    e.code,
		Message: e.message,
		Hint:    e.hint,
	}
	var buf = []byte(e.details)
	err := json.Unmarshal(buf, &response.Details)
	if nil != err {
		var s *json.SyntaxError
		if errors.As(err, &s) {
			/* A normal string is expected.  */
			response.Details = &e.details
		}
	}
	res, err := json.Marshal(response)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(e.status)
	w.Write(res)
}

type nilParameterError struct {
	f string
	p string
}

func (np nilParameterError) Error() string {
	return fmt.Sprintf("parameter %q on function %q cannot be uuid.Nil or nil", np.p, np.f)
}

func NewNilParameterError(funcName string, parameter string) error {
	return &nilParameterError{f: funcName, p: parameter}
}
