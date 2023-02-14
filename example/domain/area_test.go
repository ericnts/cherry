package domain

import (
	"context"
	"fmt"
	"github.com/ericnts/cherry"
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/example/domain/entity"
	"github.com/ericnts/log"
	"sync"
	"testing"
	"time"
)

var poStructMap sync.Map

func TestLoad(t *testing.T) {
	userService := new(UserService)
	cherry.LoadService(context.TODO(), userService)
	log.Infof("s")
}

func TestFn(t *testing.T) {
	err := A(func(befor func(), after func()) {

	})
	log.Error(err)
}

func A(bs ...func(befor func(), after func())) error {
	befor := func() {
		log.Info("befor")
	}
	after := func() {
		log.Info("after")
	}
	for i := range bs {
		bs[i](befor, after)
	}
	return nil
}

func TestPo(t *testing.T) {
	value, ok := poStructMap.Load("1")
	if !ok {
		value = base.PoStruct{Name: "zhangsan"}
	}
	log.Infof(value.(base.PoStruct).Name)
}

func TestArea(t *testing.T) {
	cherry.CallService(context.TODO(), func(s *AreaService) {
		err := s.Transaction(func() error {
			area := new(entity.Area)
			area.ID = "53ba1112ae6943a7b80d101f8b843eaf"
			area.SetParentID("c230f463040740059294cccd23f1d5bc")
			update, err := s.Repo.Update(area)
			log.Infof("%v", update)
			log.Error(err)
			return err
		})
		log.Error(err)
	})
}

func TestDelArea(t *testing.T) {
	cherry.CallService(context.TODO(), func(s *UserService) {
		_, err := s.Delete(new(entity.User), "1")
		log.Info(err)
	})
}

func TestUserService_GetByName2(t *testing.T) {

	type OfficeRes struct {
		base.Page

		ID        string        `op:"lt"`
		Name      string        `op:"gt"`
		Remark    string        `op:"une"`
		UpdatedAt string        `op:"like"`
		Code      int           `op:"in"`
		Delat     []interface{} `op:"glt"`
		CreateBy  string        `op:"-"`
	}

	ofc := &OfficeRes{
		Page:      base.Page{},
		ID:        "123",
		Name:      "123",
		Remark:    "123",
		Delat:     []interface{}{time.Now().Add(-1 * time.Hour), time.Now()},
		UpdatedAt: time.Now().Format(time.Layout),
		Code:      1000,
		CreateBy:  "张三",
	}

	qyStr, params := base.GetQuery(ofc)
	fmt.Printf("query :%s\nparams :%+v", qyStr, params)
}
