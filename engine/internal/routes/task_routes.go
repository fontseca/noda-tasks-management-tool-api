package routes

import (
	"noda/api/handler"
	"noda/engine/internal/injector"
	"noda/engine/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func InitializeForTask(r *chi.Mux) {
	s := injector.TaskService()
	h := handler.NewTaskHandler(s)

	/* For logged in user.  */

	r.Get("/user/tasks", middleware.WithBearerAuthorization(h.RetrieveTasksFromUser))
	r.Post("/user/tasks", nil)
	r.Get("/user/tasks/{task_id}", nil)
	r.Patch("/user/tasks/{task_id}", nil)
	r.Delete("/user/tasks/{task_id}", nil)

	/* For administrators (only for development purpose).  */

	r.Get("/tasks", h.RetrieveAll)
	r.Post("/tasks", nil)
	r.Get("/tasks/{task_id}", h.RetrieveTaskByID)
}
