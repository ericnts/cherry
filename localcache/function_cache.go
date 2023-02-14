package localcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ericnts/log"
	"gorm.io/gorm"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

var Function IFunctionCache = &FunctionCache{}

func init() {
	go func() {
		// 创建一个计时器
		timeTickerChan := time.Tick(time.Hour)
		for {
			Function.Release(context.TODO())
			<-timeTickerChan
		}
	}()
}

type IFunctionCache interface {
	Call(ctx context.Context, timeOut time.Duration, function interface{}, args ...interface{}) (interface{}, error)
	Release(ctx context.Context)
}

type FunctionCache struct {
	cacheMap sync.Map
}

type CacheItem struct {
	TimeOut time.Time
	Data    interface{}
}

func (c *FunctionCache) Call(ctx context.Context, timeOut time.Duration, function interface{}, args ...interface{}) (interface{}, error) {
	funcValue := reflect.ValueOf(function)
	if funcValue.Kind() != reflect.Func {
		return nil, errors.New("function is not function")
	}

	funcType := funcValue.Type()
	if funcType.IsVariadic() {
		if len(args) < funcType.NumIn()-1 {
			return nil, fmt.Errorf("argument %d length doesn't equal to provide length %d \n", funcValue.Type().NumIn(), len(args))
		}
	} else if funcType.NumIn() != len(args) {
		return nil, fmt.Errorf("argument %d length doesn't equal to provide length %d \n", funcValue.Type().NumIn(), len(args))
	}

	for i := range args {
		var reqType reflect.Type
		if funcType.IsVariadic() && i >= funcType.NumIn()-1 {
			reqType = funcType.In(funcType.NumIn() - 1).Elem()
		} else {
			reqType = funcType.In(i)
		}
		argType := reflect.TypeOf(args[i])
		if args[i] == nil {
			switch funcType.In(i).Kind() {
			case reflect.Map, reflect.Ptr, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
				continue
			}
			return nil, fmt.Errorf("第%v个参数类型不匹配，要求：%s，提供：nil", i+1, reqType.String())
		} else if reqType != argType {
			return nil, fmt.Errorf("第%v个参数类型不匹配，要求：%s，提供：%s", i+1, reqType.String(), argType.String())
		}
	}

	if funcType.NumOut() != 2 {
		return nil, errors.New("返回参数必须为两位，且最后一位是err")
	}

	if !funcType.Out(funcType.NumOut() - 1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, errors.New("last output must be error")
	}
	if !funcType.Out(funcType.NumOut() - 1).Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return nil, errors.New("last output must be error")
	}

	key, ok := getKey(funcValue, args...)
	if ok {
		value, ok := c.cacheMap.Load(key)
		if ok {
			item := value.(*CacheItem)
			if item.TimeOut.After(time.Now()) {
				return item.Data, nil
			}
		}
	} else {
		log.Warnf("无法使用缓存, %v", key)
	}

	argValues := make([]reflect.Value, 0, len(args))
	for i := 0; i < len(args); i++ {
		if args[i] == nil {
			argValues = append(argValues, reflect.New(funcValue.Type().In(i)).Elem())
		} else {
			argValues = append(argValues, reflect.ValueOf(args[i]))
		}
	}

	resultValues := funcValue.Call(argValues)
	if err := resultValues[1].Interface(); err != nil {
		return nil, err.(error)
	}
	result := resultValues[0].Interface()
	if ok {
		c.cacheMap.Store(key, &CacheItem{TimeOut: time.Now().Add(timeOut), Data: result})
	}
	return result, nil
}

func (c *FunctionCache) Release(ctx context.Context) {
	c.cacheMap.Range(func(key, value interface{}) bool {
		if value != nil {
			if value.(*CacheItem).TimeOut.Before(time.Now()) {
				c.cacheMap.Delete(key)
			}
		}
		return true
	})
}

func getKey(funcValue reflect.Value, args ...interface{}) (string, bool) {
	keys := make([]interface{}, 0, len(args)+1)
	functionName := runtime.FuncForPC(funcValue.Pointer()).Name()
	keys = append(keys, functionName)
	for i := range args {
		switch args[i].(type) {
		case *gorm.DB:
			continue
		default:
			if bytes, err := json.Marshal(args[i]); err != nil || len(bytes) > 1024 {
				log.WithError(err).Errorf("序列化失败，%+v", args[i])
				return functionName, false
			}
			keys = append(keys, args[i])
		}
	}
	if len(keys) == 1 {
		return functionName, false
	}
	bytes, err := json.Marshal(keys)
	if err != nil {
		log.WithError(err).Errorf("序列化失败，%+v", keys)
		return functionName, false
	}
	return strings.ReplaceAll(string(bytes), ":", "-"), true
}
