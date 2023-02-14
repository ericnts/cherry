package mate

import (
	"fmt"
	"github.com/ericnts/log"
	"reflect"
	"sync"
)

var (
	emptyValues = make([]reflect.Value, 0)
)

func newResourcePool(resourceType ResourceType) *resourcePool {
	return &resourcePool{
		resourceType: resourceType,
		single:       make(map[reflect.Type]*ResourceElement),
		pools:        make(map[reflect.Type]*sync.Pool),
	}
}

type resourcePool struct {
	resourceType ResourceType
	single       map[reflect.Type]*ResourceElement
	pools        map[reflect.Type]*sync.Pool
}

func (r *resourcePool) bind(single bool, f interface{}) {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()
	fKind := fType.Kind()

	if single {
		if fKind != reflect.Ptr {
			log.Panicf("单例模式只能绑定指针类型 %s.%s", fType.PkgPath(), fType.Name())
		}
		r.single[fType] = &ResourceElement{resourceType: r.resourceType, Value: reflect.ValueOf(f)}
		return
	}

	var outType reflect.Type
	switch fKind {
	case reflect.Func:
		if fType.NumIn() != 0 {
			log.Panicf("构造方法不应有参数 %s.%s", fType.PkgPath(), fType.Name())
		}
		if fType.NumOut() != 1 {
			log.Panicf("构造方法只能拥有一个返回值 %s.%s", fType.PkgPath(), fType.Name())
		}
		outType = fType.Out(0)
		if outType.Kind() != reflect.Ptr {
			log.Panicf("构造方法需要返回指针类型 %s.%s", fType.PkgPath(), fType.Name())
		}
	case reflect.Ptr:
		outType = fType
	default:
		log.Panicf("只能绑定构造方法或指针类型 %s.%s", fType.PkgPath(), fType.Name())
	}

	r.pools[outType] = &sync.Pool{
		New: func() interface{} {
			var outValue reflect.Value
			if fKind == reflect.Func {
				outValue = fValue.Call(emptyValues)[0]
			} else {
				outValue = reflect.New(fType)
			}
			element := &ResourceElement{resourceType: r.resourceType, Value: outValue}
			element.registerCall(outValue)
			allFieldsFromValue(outValue, func(fieldValue reflect.Value) {
				if element.registerWork(fieldValue) {
					return
				}
				for rType, pool := range App.poolMap {
					if r.resourceType <= rType {
						if pool.di(fieldValue, element) {
							return
						}
					}
				}
			})
			return element
		},
	}
	return
}

func (r *resourcePool) get(t reflect.Type) (*ResourceElement, bool) {
	if com, ok := r.single[t]; ok {
		return com, true
	} else if t.Kind() == reflect.Interface {
		for objT, p := range r.single {
			if objT.Implements(t) {
				return p, true
			}
		}
	}

	var pool *sync.Pool
	if t.Kind() == reflect.Interface {
		for objT, p := range r.pools {
			if objT.Implements(t) {
				pool = p
				break
			}
		}
	} else {
		pool = r.pools[t]
	}
	if pool == nil {
		return nil, false
	}
	e := pool.Get()
	if e == nil {
		return nil, false
	}
	return e.(*ResourceElement), true
}

func (r *resourcePool) di(value reflect.Value, element *ResourceElement) bool {
	if !value.IsZero() {
		return true
	}
	valueType := value.Type()
	newElement, ok := r.get(valueType)
	if !ok {
		return false
	}
	if !value.CanSet() {
		panic(fmt.Sprintf("结构体属性必须为公开访问的: %v", value.Type().String()))
	}
	value.Set(newElement.Value)
	element.appendElement(newElement)
	return true
}

func (r *resourcePool) put(element *ResourceElement) {
	if element.resourceType != r.resourceType {
		return
	}
	t := element.Value.Type()
	if _, ok := r.single[t]; ok {
		return
	}
	pool := r.pools[t]
	if pool == nil {
		return
	}
	pool.Put(element)
}

func (r *resourcePool) singleBooting(app *Application) {
	type boot interface {
		Booting(*Application)
	}
	for _, com := range r.single {
		item, ok := com.Value.Interface().(boot)
		if !ok {
			continue
		}
		item.Booting(app)
	}
}
