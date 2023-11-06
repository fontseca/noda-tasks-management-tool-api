package routes

import (
	"noda/engine/internal/injector"
	"noda/engine/internal/middleware"
	"noda/handler"

	"github.com/go-chi/chi/v5"
)

func InitializeForTasks(router chi.Router) {
	s := injector.TaskService()
	h := handler.NewTaskHandler(s)

	/* For logged in user.  */

	router.
		With(middleware.Authorization).
		Group(func(r chi.Router) {
			r.Get("/me/tasks", h.RetrieveTasksFromUser)
			r.Post("/me/tasks", nil)
			r.Get("/me/tasks/{task_id}", nil)
			r.Patch("/me/tasks/{task_id}", nil)
			r.Delete("/me/tasks/{task_id}", nil)
		})
}
