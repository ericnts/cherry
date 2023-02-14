package localcache

import (
	"context"
	"github.com/ericnts/log"
	"sync"
	"testing"
	"time"
)

type Obj struct {
	Name string
}

func TestDBCache_Get(t *testing.T) {
	officeID := "officeID"
	robotID := "robotID"
	//log.Error(DB.Del(context.TODO(), handler1))
	result, err := Function.Call(context.TODO(), time.Second, handler1, officeID, robotID)
	log.Info(result)
	log.Error(err)

	result, err = Function.Call(context.TODO(), time.Second, handler1, officeID, robotID)
	log.Info(result)
	log.Error(err)

	time.Sleep(time.Second)

	result, err = Function.Call(context.TODO(), time.Second, handler1, officeID, robotID)
	log.Info(result)
	log.Error(err)
}

func handler1(officeID, robotID string) (*Obj, error) {
	log.Info("方法被调用了")
	return &Obj{
		Name: "asdf",
	}, nil

}

func TestNormal(t *testing.T) {
	date := time.Now().AddDate(0, 0, -3)
	log.Info(date)
	date = date.AddDate(0, 0, -31)
	log.Info(date)
	date = date.AddDate(0, 0, -31)
	log.Info(date)
}

func TestFunctionCache_Call(t *testing.T) {
	re, err := Function.Call(context.TODO(), time.Second, S, "s", "b", 1)
	log.Error(err)
	log.Info(re)
}

func S(s string, a ...string) (string, error) {
	for i := range a {
		s += a[i]
	}
	return s, nil
}

func TestFunctionName(t *testing.T) {
	s := new(sync.Map)
	for i := 0; i < 10; i++ {
		s.Store(i, i)
	}

	s.Range(func(key, value interface{}) bool {
		if value.(int)%2 == 0 {
			s.Delete(key)
		}
		return true
	})
	log.Info("------------------")
	s.Range(func(key, value interface{}) bool {
		log.Infof("%v--%v", key, value)
		return true
	})
}
