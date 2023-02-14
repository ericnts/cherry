package persistence

import (
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/repository"
	"github.com/ericnts/cherry/mate"
)

var _ repository.AreaRepo = (*AreaRepository)(nil)

func init() {
	cherry.BindRepository(func() *AreaRepository {
		return new(AreaRepository)
	})
}

type AreaRepository struct {
	mate.Resource

	mate.Repository[*entity.Area]
}
