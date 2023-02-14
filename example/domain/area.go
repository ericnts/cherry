package domain

import (
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/repository"
	"github.com/ericnts/cherry/mate"
)

func init() {
	cherry.BindService(func() *AreaService {
		return new(AreaService)
	})
}

type AreaService struct {
	mate.Resource

	mate.Service[repository.AreaRepo, *entity.Area]
}
