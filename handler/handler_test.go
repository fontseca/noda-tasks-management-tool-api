package handler

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"noda/data/types"
	"os"
	"testing"
)

type JSON map[string]string

func marshal(t *testing.T, v any) []byte {
	data, err := json.Marshal(v)
	if nil != err {
		t.Fatalf("got error while marshalling: %v", err)
	}
	return data
}

func extractResponseBody(t *testing.T, body io.Reader) []byte {
	b, err := io.ReadAll(body)
	if nil != err {
		t.Fatalf("could not read response body: %v", err)
	}
	return b
}

var userID = uuid.New()

func withLoggedUser(request **http.Request) {
	var ctx = context.WithValue((*request).Context(), types.ContextKey{}, types.JWTPayload{
		UserID:   userID,
		UserRole: types.RoleUser,
	})
	*request = (*request).Clone(ctx)
}

type parameters map[string]string

func withPathParameters(request **http.Request, params parameters) {
	for key, value := range params {
		(*request).SetPathValue(key, value)
	}
}

func beQuiet() func() {
	null, _ := os.Open(os.DevNull)
	sout := os.Stdout
	serr := os.Stderr
	os.Stdout = null
	os.Stderr = null
	log.SetOutput(null)
	return func() {
		defer null.Close()
		os.Stdout = sout
		os.Stderr = serr
		log.SetOutput(os.Stderr)
	}
}
