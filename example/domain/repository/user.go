package repository

import (
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/mate"
)

type UserRepo interface {
	mate.Repo[*entity.User]

	GetByName(username string) (*entity.User, error)
}
