package server

import (
	"noda/repository"
	"noda/service"
	"sync"
)

var (
	authOnce              sync.Once
	authenticationService *service.AuthenticationService
)

func getAuthenticationService() *service.AuthenticationService {
	if nil == authenticationService {
		authOnce.Do(func() {
			userService := getUserService()
			authenticationService = service.NewAuthenticationService(userService)
		})
	}
	return authenticationService
}

var (
	groupOnce    sync.Once
	groupService service.GroupService
)

func getGroupService() service.GroupService {
	if nil == groupService {
		groupOnce.Do(func() {
			r := repository.NewGroupRepository(getDatabase())
			groupService = service.NewGroupService(r)
		})
	}
	return groupService
}

var (
	listOnce    sync.Once
	listService service.ListService
)

func getListService() service.ListService {
	if nil == listService {
		listOnce.Do(func() {
			r := repository.NewListRepository(getDatabase())
			listService = service.NewListService(r)
		})
	}
	return listService
}

var (
	taskOnce    sync.Once
	taskService *service.TaskService
)

func getTaskService() *service.TaskService {
	if nil == taskService {
		taskOnce.Do(func() {
			r := repository.NewTaskRepository(getDatabase())
			taskService = service.NewTaskService(r)
		})
	}
	return taskService
}

var (
	userOnce    sync.Once
	userService service.UserService
)

func getUserService() service.UserService {
	if nil == taskService {
		userOnce.Do(func() {
			r := repository.NewUserRepository(getDatabase())
			userService = service.NewUserService(r)
		})
	}
	return userService
}
