package injector

import (
	"noda"
	"noda/api/repository"
	"noda/api/service"
	"sync"
)

var (
	groupOnce    sync.Once
	groupService *service.GroupService
)

func GroupService() *service.GroupService {
	if nil == groupService {
		groupOnce.Do(func() {
			db := noda.Database()
			r := repository.NewGroupRepository(db)
			groupService = service.NewGroupService(r)
		})
	}
	return groupService
}
