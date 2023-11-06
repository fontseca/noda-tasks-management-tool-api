package injector

import (
	"noda"
	"noda/repository"
	"noda/service"
	"sync"
)

var (
	taskOnce    sync.Once
	taskService *service.TaskService
)

func TaskService() *service.TaskService {
	if taskService == nil {
		taskOnce.Do(func() {
			rep := repository.NewTaskRepository(noda.Database())
			taskService = service.NewTaskService(rep)
		})
	}
	return taskService
}
