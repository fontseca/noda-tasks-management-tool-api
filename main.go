package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"noda/data/types"
	"noda/failure"
	"noda/global"
	"noda/handler"
	"noda/repository"
	"noda/service"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// mustGetEnv tries to get an env var or exists.
func mustGetEnv(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if "" == value {
		log.Fatalf("could not load env var: %s", key)
	}
	return value
}

// withAuthorization returns a middleware that performs JWT-based authorization.
// It verifies the token's validity and parses its claims. If the token is
// invalid or malformed, it responds with an appropriate error. If the token is
// valid, it extracts user information from the claims and adds it to the request
// context.
func withAuthorization(next http.HandlerFunc) http.HandlerFunc {
	secret := global.Secret()
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.TrimSpace(r.Header.Get("Authorization"))
		if "" == authorization {
			w.Header().Set("WWW-Authenticate", "Bearer realm=\"access to system\"")
			failure.EmitError(w, failure.ErrMissingAuthorizationHeader)
			return
		}

		tokenStr := strings.Split(authorization, " ")[1]
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) { return secret, nil })
		if err != nil {
			var e = failure.ErrJSONWebToken.Clone()
			switch {
			default:
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, jwt.ErrInvalidKey):
				_ = e.SetDetails("Key is invalid.")
			case errors.Is(err, jwt.ErrInvalidKeyType):
				_ = e.SetDetails("Key is of invalid type.")
			case errors.Is(err, jwt.ErrTokenMalformed):
				_ = e.SetDetails("This token is not properly formed.")
			case errors.Is(err, jwt.ErrTokenSignatureInvalid):
				_ = e.SetDetails("This token signature is invalid.")
			case errors.Is(err, jwt.ErrTokenExpired):
				_ = e.
					SetDetails("This token has expired.").
					SetHint("Try singing in again.")
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				_ = e.SetDetails("This token is not valid yet.")
			case errors.Is(err, jwt.ErrTokenInvalidClaims):
				_ = e.SetDetails("This token has invalid claims.")
			case errors.Is(err, jwt.ErrInvalidType):
				_ = e.SetDetails("Invalid type for claim.")
			}
			failure.EmitError(w, e)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		id, err := uuid.Parse(claims["user_uuid"].(string))
		if err != nil {
			failure.EmitError(w, failure.ErrCorruptedClaim)
			return
		}
		ctx := context.WithValue(r.Context(), types.ContextKey{}, types.JWTPayload{
			UserID:   id,
			UserRole: types.Role(claims["user_role"].(float64))})
		r = r.Clone(ctx)
		next.ServeHTTP(w, r)
	}
}

// withAdminPrivileges returns a middleware that checks if the user has admin
// privileges.
func withAdminPrivileges(next http.HandlerFunc) http.HandlerFunc {
	return withAuthorization(func(w http.ResponseWriter, r *http.Request) {
		got := r.Context().Value(types.ContextKey{})
		if got == nil {
			log.Println("in function `WithAdminRole', got nil value for `r.Context().Value(types.ContextKey{})'")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		jwtPayload := got.(types.JWTPayload)
		if jwtPayload.UserRole != types.RoleAdmin {
			failure.EmitError(w, failure.ErrNoEnoughRights)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// withHeader returns a middleware that sets the given key-value pair in the
// HTTP response header.
func withHeader(key, value string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		})
	}
}

// withAllowedContentTypes returns a middleware that enforces a whitelist of allowed
// request Content-Types. If the Content-Type is not in the whitelist, it
// responds with a statusCode Unsupported Media Type (415).
func withAllowedContentTypes(contentTypes ...string) func(http.Handler) http.Handler {
	allowedContentTypes := make(map[string]struct{}, len(contentTypes))
	for _, contentType := range contentTypes {
		allowedContentTypes[strings.TrimSpace(strings.ToLower(contentType))] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if 0 == r.ContentLength {
				next.ServeHTTP(w, r)
				return
			}

			requestContentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
			if i := strings.Index(requestContentType, ";"); i > -1 {
				requestContentType = requestContentType[:i]
			}

			if _, ok := allowedContentTypes[requestContentType]; ok {
				next.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(http.StatusUnsupportedMediaType)
		})
	}
}

// responseRecorder inherits from an http.ResponseWriter and records some of its
// mutations for later inspection.
type responseRecorder struct {
	http.ResponseWriter
	status  int // the status code passed to WriteHeader
	written int // number of bytes written in the body
}

func newResponseRecorder(w http.ResponseWriter) *responseRecorder {
	var (
		status  = http.StatusOK
		written = 0
	)

	if r, ok := w.(*responseRecorder); ok {
		status = r.statusCode()
		written = r.bodySize()
	}

	return &responseRecorder{w, status, written}
}

func (w *responseRecorder) statusCode() int {
	if 0 == w.status {
		return http.StatusOK
	}

	return w.status
}

func (w *responseRecorder) bodySize() int {
	return w.written
}

func (w *responseRecorder) WriteHeader(code int) {
	w.status = code

	contentType := w.ResponseWriter.Header().Get("Content-Type")

	_, ok := w.ResponseWriter.(*responseRecorder)
	if !ok && "text/plain; charset=utf-8" == contentType {
		w.ResponseWriter.Header().Set("Content-Type", "application/json")
	}

	w.ResponseWriter.WriteHeader(code)
}

func (w *responseRecorder) Write(b []byte) (n int, err error) {
	notFound := bytes.Equal(b, []byte("404 page not found\n"))
	notAllowed := bytes.Equal(b, []byte("Method Not Allowed\n"))

	_, ok := w.ResponseWriter.(*responseRecorder)
	if !ok && notFound || notAllowed {
		return 0, nil
	}

	w.written = len(b)
	return w.ResponseWriter.Write(b)
}

// withNotFoundHandler is a middleware that handles any 404 not found scenario
// with a predefined response.
func withNotFoundHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := newResponseRecorder(w)
		next.ServeHTTP(recorder, r)
		if http.StatusNotFound == recorder.statusCode() {
			failure.EmitError(w, failure.ErrTargetNotFound, true)
		}
	})
}

// withNotFoundHandler is a middleware that responds with a predefined response
// when a 405 status code is received.
func withMethodNotAllowedHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := newResponseRecorder(w)
		next.ServeHTTP(recorder, r)
		if http.StatusMethodNotAllowed == recorder.statusCode() {
			failure.EmitError(w, failure.ErrNotAllowed, true)
		}
	})
}

// withRequestLoggerTo returns a middleware that logs every request made to the server
// to all the writers passed as parameter. If no target is passed, then the log entry
// is discarded.
//
// The log messages use the Common Log Format with the request latency as an extension.
// Example logs:
//
// 127.0.0.1:53582 - - [17/Apr/2024:07:41:30 -0600] "GET /me HTTP/1.1" 200 319 2.943217ms
// 127.0.0.1:41256 - - [17/Apr/2024:08:58:32 -0600] "GET /you HTTP/1.1" 404 120 113.534µs
func withRequestLoggerTo(targets ...io.Writer) func(http.Handler) http.Handler {
	writer := io.Discard

	if 0 < len(targets) {
		writer = io.MultiWriter(targets...)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := newResponseRecorder(w)
			starts := time.Now()

			next.ServeHTTP(recorder, r)

			latency := time.Since(starts)
			if time.Second < latency {
				latency = latency.Truncate(time.Microsecond)
			}

			size := recorder.bodySize()
			sizeStr := strconv.Itoa(size)
			if 0 == size {
				sizeStr = "-"
			}

			_, err := fmt.Fprintf(writer, "%s - - [%v] \"%s %s %s\" %d %s %s\n",
				r.RemoteAddr,
				time.Now().Format("02/Jan/2006:15:04:05 -0700"),
				r.Method,
				r.URL.Path,
				r.Proto,
				recorder.statusCode(),
				sizeStr,
				latency,
			)
			if nil != err {
				log.Print(err)
			}
		})
	}
}

// registerMiddlewares applies a series of middleware functions to a ServeMux and
// returns an http.Handler that chains the provided middleware functions to the
// ServeMux parameter.
func registerMiddlewares(mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) http.Handler {
	var h http.Handler = mux
	for _, middleware := range middlewares {
		if nil != middleware {
			h = middleware(h)
		}
	}
	return h
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	failLogFile, err := os.OpenFile("fail.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		log.Fatalf("could not create/open file: %v", err)
	}
	defer failLogFile.Close()

	log.SetOutput(io.MultiWriter(os.Stderr, failLogFile))

	var (
		serverPort   = mustGetEnv("SERVER_PORT")
		dbConnString = mustGetEnv("DB_CONN_STRING")
	)

	db, err := sql.Open("postgres", dbConnString)
	if nil != err {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if nil != err {
		log.Fatalf("could not ping database: %v", err)
	}

	mux := http.NewServeMux()

	var (
		userRepository = repository.NewUserRepository(db)
		userService    = service.NewUserService(userRepository)
		userHandler    = handler.NewUserHandler(userService)
	)

	mux.Handle("GET /me", withAuthorization(userHandler.HandleRetrievalOfLoggedInUser))
	mux.Handle("PATCH /me", withAuthorization(userHandler.HandleUpdateForLoggedUser))
	mux.Handle("DELETE /me", withAuthorization(userHandler.HandleRemovalOfLoggedUser))
	mux.Handle("GET /me/settings", withAuthorization(userHandler.HandleRetrievalOfLoggedUserSettings))
	mux.Handle("GET /me/settings/{setting_key}", withAuthorization(userHandler.HandleRetrievalOfOneSettingOfLoggedUser))
	mux.Handle("PUT /me/settings/{setting_key}", withAuthorization(userHandler.HandleUpdateOneSettingForLoggedUser))

	mux.Handle("GET /users", withAdminPrivileges(userHandler.HandleUsersRetrieval))
	mux.Handle("GET /users/{user_uuid}", withAdminPrivileges(userHandler.HandleRetrievalOfUserByID))
	mux.Handle("GET /users/search", withAdminPrivileges(userHandler.HandleUsersSearch))
	mux.Handle("DELETE /users/{user_uuid}", withAdminPrivileges(userHandler.HandleUserDeletion))
	mux.Handle("PUT /users/{user_uuid}/block", withAdminPrivileges(userHandler.HandleBlockUser))
	mux.Handle("DELETE /users/{user_uuid}/block", withAdminPrivileges(userHandler.HandleUnblockUser))
	mux.Handle("GET /users/blocked", withAdminPrivileges(userHandler.HandleBlockedUsersRetrieval))
	mux.Handle("PUT /users/{user_uuid}/make_admin", withAdminPrivileges(userHandler.HandleAdminPromotion))
	mux.Handle("DELETE /users/{user_uuid}/make_admin", withAdminPrivileges(userHandler.HandleDegradeAdminToUser))

	var (
		authenticationService = service.NewAuthenticationService(userService)
		authenticationHandler = handler.NewAuthenticationHandler(authenticationService)
	)

	mux.HandleFunc("POST /signup", authenticationHandler.HandleSignUp)
	mux.HandleFunc("POST /login", authenticationHandler.HandleSignIn)

	var (
		groupRepository = repository.NewGroupRepository(db)
		groupService    = service.NewGroupService(groupRepository)
		groupHandler    = handler.NewGroupHandler(groupService)
	)

	mux.Handle("GET /me/groups", withAuthorization(groupHandler.HandleGroupsRetrieval))
	mux.Handle("POST /me/groups", withAuthorization(groupHandler.HandleGroupCreation))
	mux.Handle("GET /me/groups/{group_uuid}", withAuthorization(groupHandler.HandleRetrieveGroupByID))
	mux.Handle("PATCH /me/groups/{group_uuid}", withAuthorization(groupHandler.HandleGroupUpdate))
	mux.Handle("DELETE /me/groups/{group_uuid}", withAuthorization(groupHandler.HandleGroupDeletion))

	var (
		listRepository = repository.NewListRepository(db)
		listService    = service.NewListService(listRepository)
		listHandler    = handler.NewListHandler(listService)
	)

	mux.Handle("POST /me/lists", withAuthorization(listHandler.HandleScatteredListCreation))
	mux.Handle("GET /me/lists", withAuthorization(listHandler.HandleRetrievalOfLists))
	mux.Handle("GET /me/lists/{list_uuid}", withAuthorization(listHandler.HandleScatteredListRetrievalByID))
	mux.Handle("PATCH /me/lists/{list_uuid}", withAuthorization(listHandler.HandlePartialUpdateOfScatteredList))
	mux.Handle("DELETE /me/lists/{list_uuid}", withAuthorization(listHandler.HandleScatteredListDeletion))
	mux.Handle("POST /me/groups/{group_uuid}/lists", withAuthorization(listHandler.HandleGroupedListCreation))
	mux.Handle("GET /me/groups/{group_uuid}/lists", withAuthorization(listHandler.HandleGroupedListsRetrieval))
	mux.Handle("GET /me/groups/{group_uuid}/lists/{list_uuid}", withAuthorization(listHandler.HandleGroupedListRetrievalByID))
	mux.Handle("PATCH /me/groups/{group_uuid}/lists/{list_uuid}", withAuthorization(listHandler.HandlePartialUpdateOfGroupedList))
	mux.Handle("DELETE /me/groups/{group_uuid}/lists/{list_uuid}", withAuthorization(listHandler.HandleGroupedListDeletion))

	serverLogFile, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if nil != err {
		log.Fatalf("could not create/open file: %v", err)
	}
	defer serverLogFile.Close()

	decorated := registerMiddlewares(
		mux,
		withNotFoundHandler,
		withMethodNotAllowedHandler,
		withRequestLoggerTo(os.Stdout, serverLogFile),
		withHeader("Access-Control-Allow-Credentials", "true"),
		withHeader("Access-Control-Allow-Headers", "*"),
		withHeader("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS"),
		withHeader("Access-Control-Allow-Origin", "*"),
		withHeader("Content-Type", "application/json"),
		withAllowedContentTypes("application/json"),
	)

	listener, err := net.Listen("tcp", ":"+serverPort)
	if nil != err {
		log.Fatalf("could not listen on port %s: %v", serverPort, err)
	}
	defer listener.Close()

	server := http.Server{
		Handler:                      decorated,
		DisableGeneralOptionsHandler: false,
		ReadTimeout:                  5 * time.Second,
		ReadHeaderTimeout:            5 * time.Second,
		WriteTimeout:                 5 * time.Second,
		IdleTimeout:                  120 * time.Second,
		MaxHeaderBytes:               1 << 21,
	}

	log.Fatal(server.Serve(listener))
}
