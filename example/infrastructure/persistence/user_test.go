package persistence

import (
	"context"
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/po"
	"github.com/ericnts/cherry/util"
	"github.com/ericnts/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Delete(t *testing.T) {
	err := orm.DB.AutoMigrate(
		new(entity.User),
		new(entity.Area),
	)
	assert.Nil(t, err)
	err = cherry.CallRepository(context.Background(), func(r *UserRepository) {
		_, err := r.UnscopedDeleteByID("user_id")
		assert.Nil(t, err)
		user := new(entity.User)
		user.ID = "user_id"
		user.Username = util.CreateUUID()

		var area po.Area
		area.ID = "area_id"
		user.Areas = append(user.Areas, area)

		_, err = r.Create(user)
		assert.Nil(t, err)

		user.SetUsername(util.CreateUUID())
		_, err = r.DeleteByID(user.ID)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}

func TestUser_DeleteAssociation(t *testing.T) {
	err := orm.DB.AutoMigrate(
		new(entity.User),
		new(entity.Area),
	)
	assert.Nil(t, err)
	err = cherry.CallRepository(context.Background(), func(r *UserRepository) {
		_, err := r.UnscopedDeleteByID("user_id")
		assert.Nil(t, err)
		user := new(entity.User)
		user.ID = "user_id"
		user.Username = util.CreateUUID()

		var area po.Area
		area.ID = "area_id"
		user.Areas = append(user.Areas, area)

		_, err = r.Create(user)
		assert.Nil(t, err)

		user.SetUsername(util.CreateUUID())
		_, err = r.Update(user)
		assert.Nil(t, err)
	})
	assert.Nil(t, err)
}

func TestUser_GetByID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "case1",
			args: args{id: "a84f92bc57a04aa4ab82b11297c33daa"},
			want: &entity.User{
				User: po.User{
					RecordPO: base.RecordPO{
						ID: "a84f92bc57a04aa4ab82b11297c33daa",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cherry.CallRepository(context.Background(), func(r *UserRepository) {
				got, err := r.GetByID(tt.args.id)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.ID != tt.want.ID {
					t.Errorf("GetByID() got = %v, want %v", got, tt.want)
				}
			})
			assert.Nil(t, err)
		})
	}
}

func TestUser_GetByName(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    *entity.User
		wantErr bool
	}{
		{
			name: "case1",
			args: args{username: "zhangsan"},
			want: &entity.User{
				User: po.User{
					Username: "zhangsan",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cherry.CallRepository(context.Background(), func(r *UserRepository) {
				got, err := r.GetByName(tt.args.username)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if got.Username != tt.want.Username {
					t.Errorf("GetByName() got = %v, want %v", got, tt.want)
				}
			})
			assert.Nil(t, err)
		})
	}
}
