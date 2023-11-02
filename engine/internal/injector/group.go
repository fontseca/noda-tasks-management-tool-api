package injector

import (
	"noda/api/repository"
	"noda/api/service"
	"noda/database"
	"sync"
)

var (
	groupOnce    sync.Once
	groupService *service.GroupService
)

func GroupService() *service.GroupService {
	if nil == groupService {
		groupOnce.Do(func() {
			db := database.Get()
			r := repository.NewGroupRepository(db)
			groupService = service.NewGroupService(r)
		})
	}
	return groupService
}
