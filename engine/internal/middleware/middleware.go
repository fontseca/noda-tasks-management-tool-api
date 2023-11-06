package middleware

import (
	"context"
	"errors"
	"log"
	"net/http"
	"noda"
	"noda/api/data/types"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	noda.EmitError(w, noda.ErrTargetNotFound)
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

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		value := r.Header.Get("Authorization")

		if value == "" {
			w.Header().Set("WWW-Authenticate", "Bearer realm=\"access to users\"")
			noda.EmitError(w, noda.ErrMissingAuthorizationHeader)
			return
		}

		token := strings.Split(value, " ")[1]
		t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil // TODO: Please use a better way to get a secret.
		})

		if err != nil {
			var e = noda.ErrJSONWebToken.Clone()
			switch {
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, jwt.ErrInvalidKey):
				e.SetDetails("Key is invalid.")
			case errors.Is(err, jwt.ErrInvalidKeyType):
				e.SetDetails("Key is of invalid type.")
			case errors.Is(err, jwt.ErrTokenMalformed):
				e.SetDetails("This token is not properly formed.")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				e.SetDetails("This token signature is invalid.")
			case errors.Is(err, jwt.ErrTokenExpired):
				e.SetDetails("This token has expired.").
					SetHint("Try singing in again.")
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				e.SetDetails("This token is not valid yet.")
			case errors.Is(err, jwt.ErrTokenInvalidClaims):
				e.SetDetails("This token has invalid claims.")
			case errors.Is(err, jwt.ErrInvalidType):
				e.SetDetails("Invalid type for claim.")
			}
			noda.EmitError(w, e)
			return
		}

		if !t.Valid {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims := t.Claims.(jwt.MapClaims)
		id, err := uuid.Parse(claims["user_id"].(string))
		if err != nil {
			noda.EmitError(w, noda.ErrCorruptedClaim)
			return
		}
		ctx := context.WithValue(r.Context(), types.ContextKey{}, types.JWTPayload{
			UserID:   id,
			UserRole: types.Role(claims["user_role"].(float64))})
		r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	})
}

func AdminPrivileges(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.Context().Value(types.ContextKey{})
		if got == nil {
			log.Println("in function `WithAdminRole', got nil value for `r.Context().Value(types.ContextKey{})'")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jwtPayload := got.(types.JWTPayload)
		if jwtPayload.UserRole != types.RoleAdmin {
			noda.EmitError(w, noda.ErrNoEnoughRights)
			return
		}
		next.ServeHTTP(w, r)
	})
}
