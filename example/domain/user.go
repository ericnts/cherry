package domain

import (
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/repository"
	"github.com/ericnts/cherry/mate"
)

func init() {
	cherry.BindService(func() *UserService {
		return new(UserService)
	})
}

type UserService struct {
	mate.Resource

	mate.Service[repository.UserRepo, *entity.User]
}

func (s *UserService) GetByName(username string) (*entity.User, error) {
	return s.Repo.GetByName(username)
}
