package mate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync/atomic"
	"unicode"
)

func parsePoolFunc(f interface{}) (outType reflect.Type, e error) {
	fType := reflect.TypeOf(f)
	if fType.Kind() != reflect.Func {
		e = errors.New("it's not a func")
		return
	}
	if fType.NumOut() != 1 {
		e = errors.New("return must be an object pointer")
		return
	}
	outType = fType.Out(0)
	if outType.Kind() != reflect.Ptr {
		e = errors.New("return must be an object pointer")
		return
	}
	return
}

func ParseCallFunc(f interface{}) (inTypes []reflect.Type, e error) {
	ftype := reflect.TypeOf(f)
	if ftype.Kind() != reflect.Func {
		e = errors.New("It's not a func")
		return
	}
	for i := 0; i < ftype.NumIn(); i++ {
		inType := ftype.In(i)
		inTypes = append(inTypes, inType)
		if inType.Kind() != reflect.Ptr && inType.Kind() != reflect.Interface {
			e = errors.New(fmt.Sprintf("The %d pointer parameter must be a service object", i))
			return
		}
	}
	if len(inTypes) == 0 {
		e = errors.New("The pointer parameter must be a service object")
		return
	}

	if ftype.NumOut() == 0 {
		return
	}

	if ftype.NumOut() != 1 {
		e = errors.New("The return value must be of the error type")
		return
	}
	outType := ftype.Out(0)
	if outType.Kind() != reflect.Interface {
		e = errors.New("The return value must be of the error type")
		return
	}
	if _, ok := outType.MethodByName("Error"); !ok {
		e = errors.New("The return value must be of the error type")
		return
	}
	return
}

func allFieldsFromValue(val reflect.Value, call func(reflect.Value)) {
	destVal := indirect(val)
	destType := destVal.Type()
	if destType.Kind() != reflect.Struct && destType.Kind() != reflect.Interface {
		return
	}
	for index := 0; index < destVal.NumField(); index++ {
		if destType.Field(index).Anonymous {
			allFieldsFromValue(destVal.Field(index).Addr(), call)
			continue
		}
		val := destVal.Field(index)
		kind := val.Kind()
		if kind != reflect.Ptr && kind != reflect.Interface {
			continue
		}
		call(val)
	}
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func parseMethodName(methodName string) (method, path string) {
	if strings.HasPrefix(methodName, "Get") {
		method = "GET"
		path = methodName[3:]
	} else if strings.HasPrefix(methodName, "Post") {
		method = "POST"
		path = methodName[4:]
	} else if strings.HasPrefix(methodName, "Add") {
		method = "POST"
		path = methodName[3:]
	} else if strings.HasPrefix(methodName, "Put") {
		method = "PUT"
		path = methodName[3:]
	} else if strings.HasPrefix(methodName, "Edit") {
		method = "PUT"
		path = methodName[4:]
	} else if strings.HasPrefix(methodName, "Modify") {
		method = "PUT"
		path = methodName[6:]
	} else if strings.HasPrefix(methodName, "Delete") {
		method = "DELETE"
		path = methodName[6:]
	} else if strings.HasPrefix(methodName, "Del") {
		method = "DELETE"
		path = methodName[3:]
	} else {
		return
	}
	var itemIndex atomic.Int64
	path = parsePath(path, &itemIndex)
	return
}

func parsePath(path string, itemIndex *atomic.Int64) string {
	if len(path) == 0 {
		return ""
	}
	builder := &strings.Builder{}
	ps := strings.Split(path, "Of")
	if len(ps) > 1 {
		for i := len(ps); i > 0; i-- {
			builder.WriteString(parsePath(ps[i-1], itemIndex))
		}
		return builder.String()
	}
	for len(path) != 0 {
		index := strings.Index(path, "Item")
		if index == 0 {
			itemIndex.Add(1)
			if itemIndex.Load() > 1 {
				builder.WriteString(fmt.Sprintf("/:id%d", itemIndex.Load()))
			} else {
				builder.WriteString("/:id")
			}
			path = path[4:]
			continue
		}
		if index == -1 {
			index = len(path)
		}
		currentPath := path[:index]
		builder.WriteString("/")
		lowIndex := len(currentPath)
		for i, r := range []rune(currentPath) {
			if unicode.IsLower(r) {
				lowIndex = i
				break
			}
		}
		if lowIndex == 0 {
			builder.WriteString(currentPath)
		} else if lowIndex == 1 {
			builder.WriteString(strings.ToLower(currentPath[:1]))
			builder.WriteString(currentPath[1:])
		} else if lowIndex == len(currentPath) {
			builder.WriteString(strings.ToLower(currentPath))
		} else {
			builder.WriteString(strings.ToLower(currentPath[:lowIndex-1]))
			builder.WriteString(currentPath[lowIndex-1:])
		}

		path = path[index:]
	}
	return builder.String()
}
