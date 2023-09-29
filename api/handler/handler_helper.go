package handler

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"noda/api/data/types"
	"noda/failure"
	"strconv"
)

func GetQueryParameter(r *http.Request, key, fallback string) string {
	k := r.URL.Query().Get(key)
	if k == "" {
		k = fallback
	}
	return k
}

func ParsePagination(w http.ResponseWriter, r *http.Request) *types.Pagination {
	page, err := strconv.ParseInt(GetQueryParameter(r, "page", "1"), 10, 64)
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

	rpp, err := strconv.ParseInt(GetQueryParameter(r, "rpp", "10"), 10, 64)
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

	maxValidBeforeOverflow := (math.MaxInt64 / rpp) - 1
	if page > maxValidBeforeOverflow {
		page = maxValidBeforeOverflow
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
