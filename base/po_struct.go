package base

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/ericnts/cherry/exception"
)

var poStructMap sync.Map

type PoField struct {
	Column  string
	Comment string
}

type PoStruct struct {
	Name               string
	UpdateAssociations []string
	DeleteAssociations []string
	FieldMap           map[string]PoField
	Indexes            []IndexCheck
	Verify             *VerifyCheck
}

type IndexCheck struct {
	Key     string
	Fields  []string
	Columns []string
	Err     error
}

type VerifyCheck struct {
	Fields  []string
	Columns []string
}

type Index struct {
	Key     string
	Columns []string
	Values  []interface{}
	Err     error
}

type Verify struct {
	Columns []string
	Values  []interface{}
}

func GetUpdateAssociation(po PO) []string {
	poStruct := getPoStruct(reflect.ValueOf(po).Elem())
	return poStruct.UpdateAssociations
}

func GetDeleteAssociation(po PO) []string {
	poStruct := getPoStruct(reflect.ValueOf(po).Elem())
	return poStruct.DeleteAssociations
}

func GetIndexes(po PO) []Index {
	poValue := reflect.ValueOf(po).Elem()
	poStruct := getPoStruct(poValue)
	if len(poStruct.Indexes) == 0 {
		return nil
	}
	res := make([]Index, 0, len(poStruct.Indexes))
	for _, indexCheck := range poStruct.Indexes {
		values := make([]interface{}, len(indexCheck.Fields))
		for i, fieldKey := range indexCheck.Fields {
			values[i] = poValue.FieldByName(fieldKey).Interface()
		}
		res = append(res, Index{
			Key:     indexCheck.Key,
			Columns: indexCheck.Columns,
			Values:  values,
			Err:     indexCheck.Err,
		})
	}
	return res
}

func GetVerify(po PO) *Verify {
	poValue := reflect.ValueOf(po).Elem()
	poStruct := getPoStruct(poValue)
	if nil == poStruct.Verify || len(poStruct.Verify.Columns) == 0 {
		return nil
	}
	values := make([]interface{}, len(poStruct.Verify.Columns))
	for j, fieldKey := range poStruct.Verify.Fields {
		val := poValue.FieldByName(fieldKey)
		if reflect.Invalid != val.Kind() {
			values[j] = val.Interface()
		}
	}
	return &Verify{
		Columns: poStruct.Verify.Columns,
		Values:  values,
	}
}

func getPoStruct(poValue reflect.Value) PoStruct {
	poStructKey := poValue.Type().String()
	value, ok := poStructMap.Load(poStructKey)
	if ok {
		return value.(PoStruct)
	}
	poStruct := PoStruct{FieldMap: make(map[string]PoField)}
	verifyCheck := new(VerifyCheck)
	allFieldsFromValue(poValue, func(value reflect.Value, field reflect.StructField) {
		tag := field.Tag
		ormTag, ok := tag.Lookup("gorm")
		if !ok {
			return
		}
		var column, comment string
		for _, ormItem := range strings.Split(ormTag, ";") {
			if strings.Index(ormItem, "column:") == 0 {
				column = strings.TrimSpace(ormItem[7:])
			} else if strings.Index(ormItem, "comment:") == 0 {
				comment = strings.TrimSpace(ormItem[8:])
			}
		}
		if len(column) > 0 {
			if len(comment) == 0 {
				comment = column
			}
			poStruct.FieldMap[field.Name] = PoField{
				Column:  column,
				Comment: comment,
			}
		}

		if checkTag, ok := tag.Lookup("cherry"); ok {
			for _, checkItem := range strings.Split(checkTag, ";") {
				if strings.Index(checkItem, "auto") == 0 && strings.Index(ormTag, "many2many") != -1 {
					// 保存多对多关联关系
					if checkItem == "auto" || checkItem == "autoUpdate" {
						poStruct.UpdateAssociations = append(poStruct.UpdateAssociations, field.Name)
					}
					if checkItem == "auto" || checkItem == "autoDelete" {
						poStruct.DeleteAssociations = append(poStruct.DeleteAssociations, field.Name)
					}
				} else if checkItem == "verify" {
					if len(column) > 0 {
						verifyCheck.Columns = append(verifyCheck.Columns, column)
						verifyCheck.Fields = append(verifyCheck.Fields, field.Name)
					}
				} else if strings.Index(checkItem, "index") == 0 {
					var indexKey string
					if items := strings.Split(checkItem, ":"); len(items) == 2 {
						indexKey = strings.TrimSpace(items[1])
					}
					if len(indexKey) == 0 {
						indexKey = field.Name
					}
					poStruct.Indexes = append(poStruct.Indexes, IndexCheck{Key: indexKey})
				}
			}
		}
	})

	if len(verifyCheck.Fields) > 0 {
		poStruct.Verify = verifyCheck
	}

	for i, index := range poStruct.Indexes {
		var comments []string
		for _, fieldName := range strings.Split(index.Key, "_") {
			commentIgnore := strings.Index(fieldName, "#") == 0
			if commentIgnore {
				fieldName = fieldName[1:]
			}
			poField, ok := poStruct.FieldMap[fieldName]
			if !ok {
				continue
			}
			if !commentIgnore {
				comments = append(comments, poField.Comment)
			}
			poStruct.Indexes[i].Fields = append(poStruct.Indexes[i].Fields, fieldName)
			poStruct.Indexes[i].Columns = append(poStruct.Indexes[i].Columns, poField.Column)
		}
		if len(comments) == 0 {
			poStruct.Indexes[i].Err = exception.Custom(exception.DataRepeat, "数据保存失败，数据重复", errors.New(index.Key))
		} else {
			poStruct.Indexes[i].Err = exception.Custom(exception.DataRepeat, fmt.Sprintf("数据保存失败，%s已存在", strings.Join(comments, "或")), errors.New(index.Key))
		}
	}

	value = poStruct
	poStructMap.Store(poStructKey, value)

	return poStruct
}

func allFieldsFromValue(val reflect.Value, call func(reflect.Value, reflect.StructField)) {
	destVal := val
	for destVal.Kind() == reflect.Ptr || destVal.Kind() == reflect.Interface {
		destVal = destVal.Elem()
	}
	destType := destVal.Type()
	for index := 0; index < destVal.NumField(); index++ {
		if destType.Field(index).Anonymous {
			allFieldsFromValue(destVal.Field(index).Addr(), call)
			continue
		}
		call(destVal.Field(index), destType.Field(index))
	}
}
