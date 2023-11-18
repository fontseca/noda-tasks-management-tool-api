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
	"noda"
	"noda/data/types"
	"regexp"
	"strconv"
	"strings"
)

func extractQueryParameter(r *http.Request, key, fallback string) string {
	k := strings.Trim(r.URL.Query().Get(key), " \t\n")
	if strings.Compare(k, "") == 0 {
		return fallback
	}
	return k
}

func parsePagination(w http.ResponseWriter, r *http.Request) *types.Pagination {
	page, err := strconv.ParseInt(extractQueryParameter(r, "page", "1"), 10, 64)
	if err != nil {
		err, _ := err.(*strconv.NumError)
		var e = noda.ErrBadQueryParameter.Clone()
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, strconv.ErrSyntax):
			noda.EmitError(w, e.SetDetails("Value for parameter \"page\" is not a valid decimal number."))
		case errors.Is(err, strconv.ErrRange):
			var details = fmt.Sprintf("Parameter \"page\" has value %s, which is out of raneg for signed 64 bits numbers.", err.Num)
			noda.EmitError(w, e.SetDetails(details))
		}
		return nil
	}

	agg := noda.AggregateDetails{}

	if page <= 0 {
		agg.Append("The parameter \"page\" must be a positive number.")
	}

	rpp, err := strconv.ParseInt(extractQueryParameter(r, "rpp", "10"), 10, 64)
	if err != nil {
		err, _ := err.(*strconv.NumError)
		var e = noda.ErrBadQueryParameter.Clone()
		switch {
		default:
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		case errors.Is(err, strconv.ErrSyntax):
			noda.EmitError(w, e.SetDetails("Value for parameter \"rpp\" is not a valid decimal number."))
		case errors.Is(err, strconv.ErrRange):
			var details = fmt.Sprintf("Parameter \"rpp\" has value %s, which is out of raneg for signed 64 bits numbers.", err.Num)
			noda.EmitError(w, e.SetDetails(details))
		}
		return nil
	}

	if rpp <= 0 {
		agg.Append("The parameter \"rpp\" must be a positive number.")
	}

	if !agg.Has() {
		return &types.Pagination{
			Page: page,
			RPP:  rpp,
		}
	}

	noda.EmitError(w, noda.ErrBadQueryParameter.Clone().SetDetails(agg.Error()))
	return nil
}

func extractSorting(w http.ResponseWriter, r *http.Request) string {
	sortBy := extractQueryParameter(r, "sort_by", "")
	if len(r.URL.Query()["sort_by"]) > 1 {
		noda.EmitError(w, noda.ErrMultipleValuesForQueryParameter.
			Clone().
			FormatDetails("sort_by"))
		return "?"
	}
	if "" == sortBy {
		return sortBy
	}
	matched, err := regexp.MatchString(`^(?:(?:\+|-)[_a-zA-Z][_a-zA-Z0-9]+)$`, sortBy)
	if err != nil {
		log.Println(err)
	}
	if matched {
		return sortBy
	}
	details, _ := json.Marshal([]string{
		"Must start with either one plus sign (+) or one minus sign (-).",
		"Must contain one or more word characters (alphanumeric characters and underscores).",
	})
	noda.EmitError(w, noda.ErrQueryParameterNotParsed.Clone().SetDetails(string(details)))
	return "?"
}

func decodeJSONRequestBody(w http.ResponseWriter, r *http.Request, out any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var decoder = json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var err = decoder.Decode(&out)
	if nil != err {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		default:
			return err
		case errors.As(err, &syntaxError):
			var details = fmt.Sprintf("Body contains ill-formed JSON at position %d.",
				syntaxError.Offset)
			return errors.New(details)
		case errors.As(err, &unmarshalTypeError):
			var details = fmt.Sprintf("Body contains an invalid value for the %q field at position %d.",
				unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return errors.New(details)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("Body contains ill-formed JSON.")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			var field = strings.TrimPrefix(err.Error(), "json: unknown field ")
			var details = fmt.Sprintf("Body contains an unknown field: %s.", field)
			return errors.New(details)
		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty.")
		case err.Error() == "http: request body too large":
			return errors.New("Body must not be larger than 1MB.")
		}
	}
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("Body must only contain a single JSON object.")
	}
	return nil
}

func parsePathParameterToUUID(r *http.Request, parameter string) (uuid.UUID, error) {
	var key = chi.URLParam(r, parameter)
	id, err := uuid.Parse(key)
	if nil != err {
		switch {
		default:
			log.Println(err)
			return uuid.Nil, err
		case strings.Contains(err.Error(), "invalid UUID format"):
			return uuid.Nil, noda.ErrInvalidUUIDFormat
		case strings.Contains(err.Error(), "invalid UUID length"):
			return uuid.Nil, noda.ErrInvalidUUIDLength
		}
	}
	return id, nil
}

func extractUserPayload(r *http.Request) (userID uuid.UUID, userRole types.Role) {
	payload := r.Context().Value(types.ContextKey{}).(types.JWTPayload)
	return payload.UserID, payload.UserRole
}

func redirect(w http.ResponseWriter, r *http.Request, to string) {
	var (
		scheme = "http://"
		host   = r.Host
	)
	if nil != r.TLS {
		scheme = "https://"
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s%s", scheme, host, to))
	w.WriteHeader(http.StatusSeeOther)
}
