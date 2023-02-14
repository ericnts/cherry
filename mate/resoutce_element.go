package mate

import (
	"reflect"
)

var (
	emptyResourceElement *ResourceElement
)

type ResourceElement struct {
	resourceType ResourceType
	workers      []reflect.Value
	calls        []BeginRequest
	Value        reflect.Value
}

func (r *ResourceElement) Interface() interface{} {
	return r.Value.Interface()
}

func (r *ResourceElement) registerWork(value reflect.Value) bool {
	if value.Kind() == reflect.Interface && workerType.AssignableTo(value.Type()) && value.CanSet() {
		//如果是运行时对象
		r.workers = append(r.workers, value)
		return true
	}
	return false
}

func (r *ResourceElement) registerCall(value reflect.Value) {
	if value.IsNil() {
		return
	}
	if br, ok := value.Interface().(BeginRequest); ok {
		r.calls = append(r.calls, br)
	}
}

func (r *ResourceElement) setWork(worker Worker) {
	if worker == nil {
		return
	}
	workerValue := reflect.ValueOf(worker)
	for i := 0; i < len(r.workers); i++ {
		r.workers[i].Set(workerValue)
	}
	for i := 0; i < len(r.calls); i++ {
		r.calls[i].BeginRequest(worker)
	}
	worker.AddFreeElement(r)
}

func (r *ResourceElement) appendElement(element *ResourceElement) {
	if element == nil {
		return
	}
	r.workers = append(r.workers, element.workers...)
	r.calls = append(r.calls, element.calls...)
}
