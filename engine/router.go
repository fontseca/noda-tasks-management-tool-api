package engine

import (
	"encoding/json"
	"net/http"
	"noda/api/handler"
	"noda/engine/internal/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// for _, route := range router.Routes() {
// 	for method := range route.Handlers {
// 		fmt.Println(method, route.Pattern)
// 	}
// }

// chi.Walk(router, func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
// 	fmt.Printf("%s %s\n", method, route)
// 	return nil
// })

type serviceHandlers struct {
	TaskHandler *handler.TaskHandler
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	answer := struct {
		Message string `json:"message"`
	}{
		Message: "not found",
	}
	res, err := json.Marshal(answer)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(res)
}

func initializeRouter(handlers *serviceHandlers) *chi.Mux {
	r := chi.NewRouter()
	setUpMiddlewares(r)
	r.NotFound(notFoundHandler)
	routes.InitializeForTask(r, handlers.TaskHandler)
	return r
}

func setUpMiddlewares(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))
}
