package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"noda/api/data/types"
	"noda/failure"
	"regexp"
	"strconv"
	"strings"
)

func parseQueryParameter(r *http.Request, key, fallback string) string {
	k := r.URL.Query().Get(key)
	if strings.Compare(k, "") == 0 {
		k = fallback
	}
	return k
}

func parsePagination(w http.ResponseWriter, r *http.Request) *types.Pagination {
	page, err := strconv.ParseInt(parseQueryParameter(r, "page", "1"), 10, 64)
	if err != nil {
		err, _ := err.(*strconv.NumError)
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, strconv.ErrSyntax):
			failure.Emit(w, http.StatusBadRequest, "query parameter failure", "\"page\" is not a valid decimal number")
		case errors.Is(err, strconv.ErrRange):
			failure.Emit(w, http.StatusBadRequest, "query parameter failure",
				fmt.Errorf("%q has as value %s, which is out of range for a signed 64 bits number", "page", err.Num))
		}
		return nil
	}

	agg := failure.Aggregation{}

	if page <= 0 {
		agg.Append(errors.New("\"page\" must be a positive number"))
	}

	rpp, err := strconv.ParseInt(parseQueryParameter(r, "rpp", "10"), 10, 64)
	if err != nil {
		err, _ := err.(*strconv.NumError)
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, strconv.ErrSyntax):
			failure.Emit(w, http.StatusBadRequest, "query parameter failure", "\"rpp\" is not a valid decimal number")
		case errors.Is(err, strconv.ErrRange):
			failure.Emit(w, http.StatusBadRequest, "query parameter failure",
				fmt.Errorf("%q has as value %s, which is out of range for a signed 64 bits number", "rpp", err.Num))
		}
		return nil
	}

	if rpp <= 0 {
		agg.Append(errors.New("\"rpp\" must be a positive number"))
	}

	if !agg.Has() {
		return &types.Pagination{
			Page: page,
			RPP:  rpp,
		}
	}

	failure.Emit(w, http.StatusBadRequest, "query parameter failure", agg.Dump())
	return nil
}

func parseSorting(w http.ResponseWriter, r *http.Request) string {
	sortBy := parseQueryParameter(r, "sort_by", "")
	if len(r.URL.Query()["sort_by"]) > 1 {
		failure.Emit(w, http.StatusBadRequest, "too much values for query parameter: \"sort_by\"",
			"please provide only one parameter value for key \"sort_by\"")
		return ""
	}
	matched, err := regexp.MatchString(`^(?:(?:\+|-)[_a-zA-Z][_a-zA-Z0-9]+)$`, sortBy)
	if err != nil {
		log.Println(err)
	}
	if matched {
		return sortBy
	}
	failure.Emit(
		w, http.StatusBadRequest,
		"could not parse query parameter: \"sort_by\"",
		[]string{
			"must start with either one plus sign (+) or one minus sign (-)",
			"must contain one or more word characters (alphanumeric characters and underscores)",
		})
	return ""
}

type malformedRequest struct {
	status           int
	message, details string
}

func (mr *malformedRequest) Error() string {
	return mr.details
}

func decodeJSONRequestBody(w http.ResponseWriter, r *http.Request, out any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var decoder = json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var err = decoder.Decode(&out)
	if nil != err {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		const message = "bad request body"
		switch {
		default:
			return err
		case errors.As(err, &syntaxError):
			var details = fmt.Sprintf("request body contains ill-formed JSON at position %d",
				syntaxError.Offset)
			return &malformedRequest{
				status:  http.StatusBadRequest,
				message: message,
				details: details,
			}
		case errors.As(err, &unmarshalTypeError):
			var details = fmt.Sprintf("request body contains an invalid value for the %q field at position %d",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{
				status:  http.StatusBadRequest,
				message: message,
				details: details,
			}
		case errors.Is(err, io.ErrUnexpectedEOF):
			var details = "request body contains ill-formed JSON"
			return &malformedRequest{
				status:  http.StatusBadRequest,
				message: message,
				details: details,
			}
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			var field = strings.TrimPrefix(err.Error(), "json: unknown field ")
			var details = fmt.Sprintf("request body contains unknown field %s", field)
			return &malformedRequest{
				status:  http.StatusBadRequest,
				message: message,
				details: details,
			}
		case errors.Is(err, io.EOF):
			return &malformedRequest{
				status:  http.StatusBadRequest,
				message: message,
				details: "request body must not be empty",
			}
		case err.Error() == "http: request body too large":
			return &malformedRequest{
				status:  http.StatusBadRequest,
				message: message,
				details: "request body must not be larger than 1MB",
			}
		}
	}
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return &malformedRequest{
			status:  http.StatusBadRequest,
			message: "bad request body",
			details: "request body must only contain a single JSON object",
		}
	}
	return nil
}

func parsePathParameterToUUID(w http.ResponseWriter, r *http.Request, parameter string) (uuid.UUID, error) {
	var key = chi.URLParam(r, parameter)
	id, err := uuid.Parse(key)
	if nil != err {
		var message = fmt.Sprintf("failure with %q", key)
		failure.Emit(w, http.StatusBadRequest, message, err)
		return uuid.Nil, err
	}
	return id, nil
}
