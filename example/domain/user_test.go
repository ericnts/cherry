package domain

import (
	"context"
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/cherry/example/domain/po"
	_ "github.com/ericnts/cherry/example/infrastructure"
	"github.com/ericnts/cherry/util"
	"github.com/ericnts/orm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_Update(t *testing.T) {
	err := orm.DB.AutoMigrate(
		new(entity.User),
		new(entity.Area),
	)
	assert.Nil(t, err)
	err = cherry.CallService(context.Background(), func(r *UserService) {
		_, err := r.Repo.UnscopedDeleteByID("user_id")
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

func TestUserService_Create(t *testing.T) {
	type args struct {
		user *entity.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			args: args{
				user: &entity.User{
					User: po.User{
						RecordPO: base.RecordPO{
							ID: "qwer",
						},
						Username: "zhangsan1",
						Password: "zhangsan",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cherry.CallService(context.Background(), func(s *UserService) {
				got, err := s.Create(tt.args.user)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				assert.True(t, got > 0)
			})
			assert.Nil(t, err)
		})
	}
}

func TestUserService_GetByID(t *testing.T) {
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
			args: args{
				id: "a84f92bc57a04aa4ab82b11297c33daa",
			},
			want: &entity.User{
				User: po.User{
					Username: "zhangsan",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cherry.CallService(context.Background(), func(s *UserService) {
				got, err := s.GetByID(tt.args.id)
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

func TestUserService_GetByName(t *testing.T) {
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
			args: args{
				username: "zhangsan",
			},
			want: &entity.User{
				User: po.User{
					Username: "zhangsan",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cherry.CallService(context.Background(), func(s *UserService) {
				got, err := s.GetByName(tt.args.username)
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

func TestUserService_Update(t *testing.T) {
	_ = cherry.CallService(context.Background(), func(s *UserService) {
		_, err := s.Repo.UnscopedDeleteByID("user1", "user2")
		assert.Nil(t, err)
		user := &entity.User{
			User: po.User{
				RecordPO: base.RecordPO{
					ID: "user1",
				},
				No:         1,
				CategoryID: "officeID",
				Username:   "zhangsan",
				Password:   "zhangsan",
			},
		}
		_, err = s.Create(user)
		assert.Nil(t, err)
		user2 := &entity.User{
			User: po.User{
				RecordPO: base.RecordPO{
					ID: "user2",
				},
				No:         2,
				CategoryID: "officeID",
				Username:   "lisi",
				Password:   "zhangsan",
			},
		}
		_, err = s.Create(user2)
		assert.Nil(t, err)
		_, err = s.Update(user)
		assert.Nil(t, err)
	})

}
