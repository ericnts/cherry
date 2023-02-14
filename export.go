package cherry

import (
	"context"
	"errors"
	"fmt"
	"github.com/ericnts/cherry/mate"
	"reflect"
)

var (
	app = mate.App
)

func NewApplication() *mate.Application {
	return app
}

func AfterRun(f func(*mate.Application)) {
	app.After(f)
}

func Prepare(f func(*mate.Application)) {
	app.Prepare(mate.PLMiddle, f)
}

func PrepareWithLevel(level mate.PrepareLevel, f func(*mate.Application)) {
	app.Prepare(level, f)
}

func DoPrepare() {
	app.DoPrepare()
}

func BindApplication(f interface{}) {
	app.Prepare(mate.PLMiddle, func(app *mate.Application) {
		app.BindApplication(f)
	})
}

func BindService(f interface{}) {
	app.Prepare(mate.PLMiddle, func(app *mate.Application) {
		app.BindService(f)
	})
}

func BindFactory(f interface{}) {
	app.Prepare(mate.PLMiddle, func(app *mate.Application) {
		app.BindFactory(f)
	})
}

func BindRepository(f interface{}) {
	app.Prepare(mate.PLMiddle, func(app *mate.Application) {
		app.BindRepository(f)
	})
}

func BindInfrastructure(single bool, com interface{}) {
	app.Prepare(mate.PLMiddle, func(app *mate.Application) {
		app.BindInfrastructure(single, com)
	})
}

func CallApplication(c context.Context, f interface{}) error {
	return call(c, mate.RTApplication, f)
}

func CallService(c context.Context, f interface{}) error {
	return call(c, mate.RTService, f)
}

func CallRepository(c context.Context, f interface{}) error {
	return call(c, mate.RTRepository, f)
}

func CallInfrastructure(c context.Context, f interface{}) error {
	return call(c, mate.RTInfrastructure, f)
}

func Call(c context.Context, f interface{}) error {
	app.DoPrepare()

	inTypes, err := mate.ParseCallFunc(f)
	if err != nil {
		return err
	}
	args := make([]reflect.Value, 0, len(inTypes))
	worker := mate.NewWorker(c)
	for i := range inTypes {
		element, ok := app.GetAnonymous(worker, inTypes[i])
		if !ok {
			return errors.New(fmt.Sprintf("没有找到实例%d", i))
		}
		args = append(args, element.Value)
	}
	returnValue := reflect.ValueOf(f).Call(args)
	worker.Free()
	if len(returnValue) > 0 && !returnValue[0].IsNil() {
		returnInterface := returnValue[0].Interface()
		err, _ := returnInterface.(error)
		return err
	}
	return nil
}

func call(c context.Context, resourceType mate.ResourceType, f interface{}) error {
	app.DoPrepare()

	inType, err := mate.ParseCallFunc(f)
	if err != nil {
		return err
	}
	worker := mate.NewWorker(c)
	element, ok := app.Get(worker, resourceType, inType[0])
	if !ok {
		return errors.New("没有找到实例")
	}
	returnValue := reflect.ValueOf(f).Call([]reflect.Value{element.Value})
	worker.Free()
	if len(returnValue) > 0 && !returnValue[0].IsNil() {
		returnInterface := returnValue[0].Interface()
		err, _ := returnInterface.(error)
		return err
	}
	return nil
}

func LoadApplication(c context.Context, target any) error {
	return load(c, mate.RTApplication, target)
}

func LoadService(c context.Context, target any) error {
	return load(c, mate.RTService, target)
}

// Load 只负责装在对象，不负责回收
func load(c context.Context, resourceType mate.ResourceType, target any) error {
	app.DoPrepare()

	targetType := reflect.TypeOf(target)
	if targetType.Kind() != reflect.Ptr && targetType.Kind() != reflect.Interface {
		return errors.New("the target must be point or interface")
	}

	worker := mate.NewWorker(c)
	element, ok := app.Get(worker, resourceType, targetType)
	if !ok {
		return errors.New("没有找到实例")
	}
	reflect.ValueOf(target).Elem().Set(element.Value.Elem())
	return nil
}
