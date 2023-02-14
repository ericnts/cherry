package persistence

import (
	"context"
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/po"
	"github.com/ericnts/cherry/util"
	"github.com/ericnts/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestArea_DeleteAssociation(t *testing.T) {
	err := orm.DB.AutoMigrate(
		new(entity.User),
		new(entity.Area),
	)
	assert.Nil(t, err)
	err = cherry.CallRepository(context.Background(), func(r *AreaRepository) {
		r.UnscopedDeleteByID("area_id")
		user := po.User{}
		user.ID = "user_id"
		user.Username = util.CreateUUID()

		area := new(entity.Area)
		area.ID = "area_id"
		area.Users = append(area.Users, user)
		_, err := r.Create(area)
		assert.Nil(t, err)

		area2 := new(entity.Area)
		area2.ID = "area_id2"
		area2.ParentIDs = "area_id"
		_, err = r.Create(area2)
		assert.Nil(t, err)

		err = r.DeleteAssociationByID("Users", "area_id")
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}

func TestArea_Create(t *testing.T) {
	_ = orm.DB.AutoMigrate(
		new(entity.Area),
	)
	err := cherry.CallRepository(context.Background(), func(r *AreaRepository) {
		_, err := r.DeleteByID("53ba1112ae6943a7b80d101f8b843eaf")
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}
