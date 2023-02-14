package repository

import (
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/mate"
)

type AreaRepo interface {
	mate.Repo[*entity.Area]
}
