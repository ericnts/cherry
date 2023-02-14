package persistence

import (
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/repository"
	"github.com/ericnts/cherry/mate"
)

var _ repository.UserRepo = (*UserRepository)(nil)

func init() {
	cherry.BindRepository(func() *UserRepository {
		return new(UserRepository)
	})
}

type UserRepository struct {
	mate.Repository[*entity.User]
}

func (u *UserRepository) GetByName(username string) (*entity.User, error) {
	res := new(entity.User)
	err := u.DB().First(res, "username = ?", username).Error
	return res, err
}
