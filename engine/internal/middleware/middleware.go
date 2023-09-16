package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"noda/api/data/types"
	"noda/failure"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	failure.Emit(w, http.StatusNotFound,
		"Target not found.", fmt.Sprintf("Could not find resource `%s'.", r.URL))
}

func LetOptionsPassThrough(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func WithBearerAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Header.Get("Authorization")

		if value == "" {
			w.Header().Set("WWW-Authenticate", "Bearer realm=\"access to users\"")
			failure.Emit(w, http.StatusUnauthorized, "bad authorization request", "no Authorization header provided")
			return
		}

		token := strings.Split(value, " ")[1]
		t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		})

		if err != nil {
			details := ""
			switch {
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, jwt.ErrInvalidKey):
				details = jwt.ErrInvalidKey.Error() // key is invalid
			case errors.Is(err, jwt.ErrInvalidKeyType):
				details = jwt.ErrInvalidKeyType.Error() // key is of invalid type
			case errors.Is(err, jwt.ErrTokenMalformed):
				details = "token is not properly formed"
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				details = jwt.ErrTokenSignatureInvalid.Error() // token signature is invalid
			case errors.Is(err, jwt.ErrTokenExpired):
				details = "token has expired: sign in again" // token is expired
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				details = jwt.ErrTokenNotValidYet.Error() // token is not valid yet
			case errors.Is(err, jwt.ErrTokenInvalidClaims):
				details = jwt.ErrTokenInvalidClaims.Error() // token has invalid claims
			case errors.Is(err, jwt.ErrInvalidType):
				details = jwt.ErrInvalidType.Error() // invalid type for claim
			}
			failure.Emit(w, http.StatusUnauthorized, "jwt failure", details)
			return
		}

		if !t.Valid {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims := t.Claims.(jwt.MapClaims)
		id, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			failure.Emit(w, http.StatusUnauthorized, "jwt failure", "a claim in jwt seems corrupted")
			return
		}
		ctx := context.WithValue(r.Context(), types.ContextKey{}, id)
		r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	})
}
