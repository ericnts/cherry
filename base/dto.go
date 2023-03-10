package base

import (
	"fmt"
	"github.com/ericnts/cherry/util"
	"reflect"

	"github.com/mailru/easyjson/buffer"
)

type Query interface {
	Query() (string, []interface{})
	GetOrder() string
	GetPreloads() []interface{}
	GetOmits() []string
}

type Preload struct {
	Query string
	Args  []interface{}
}

type EmptyQuery struct {
	Order    string        `json:"order" form:"order"` //排序
	Preloads []interface{} `json:"-"`                  //预加载
	Omits    []string      `json:"-"`                  //忽略的字段
}

func (e *EmptyQuery) Query() (string, []interface{}) {
	return "", nil
}

func (e *EmptyQuery) GetOrder() string {
	return e.Order
}

func (e *EmptyQuery) GetPreloads() []interface{} {
	return e.Preloads
}

func (e *EmptyQuery) GetOmits() []string {
	return e.Omits
}

type PageQuery interface {
	Query

	Limit() int
	Offset() int
}

type Page struct {
	EmptyQuery

	PageNO   int `json:"pageNo" form:"pageNo"`     //分页编号
	PageSize int `json:"pageSize" form:"pageSize"` //分页大小
}

func (p *Page) Limit() int {
	return p.PageSize
}

func (p *Page) Offset() int {
	return (p.PageNO - 1) * p.PageSize
}

func GetQuery(query Query) (string, []interface{}) {
	ref := reflect.TypeOf(query).Elem()
	rev := reflect.ValueOf(query).Elem()

	buf := buffer.Buffer{}
	params := make([]interface{}, 0)

	joinStr := " and "
	filedNub := ref.NumField()
	for i := 0; i < filedNub; i++ {
		fieldKey := ref.Field(i)
		fieldVal := rev.Field(i)
		_, zeroIgnore := fieldKey.Tag.Lookup(zeroValIgnore)
		if fieldKey.Type.Kind() != reflect.Ptr {
			if (fieldKey.Type.Kind() == reflect.Struct) || (!zeroIgnore && fieldVal.IsZero()) {
				continue
			}
		} else {
			if fieldVal.IsNil() {
				continue
			}
		}

		optionVal, _ := fieldKey.Tag.Lookup(operationTag)
		columnName := util.UnderscoreName(fieldKey.Name)

		switch optionVal {
		case opIgnore:
			continue
		case opLike:
			buf.AppendString(fmt.Sprintf("%s like ?", columnName))
			fieldVal = reflect.ValueOf(fmt.Sprintf("%%%s%%", fieldVal.String()))
		case opIn:
			buf.AppendString(fmt.Sprintf("%s in (?)", columnName))
		case opLT:
			buf.AppendString(fmt.Sprintf("%s < ?", columnName))
		case opGT:
			buf.AppendString(fmt.Sprintf("%s > ?", columnName))
		case opBetween:
			bwt, ok := fieldVal.Interface().([]string)
			if !ok || len(bwt) != 2 {
				continue
			}
			buf.AppendString(fmt.Sprintf("%s between ? and ? %s", columnName, joinStr))
			params = append(params, bwt[0], bwt[1])
			continue
		case opUnEqual:
			buf.AppendString(fmt.Sprintf("%s <> ?", columnName))
		default:
			buf.AppendString(fmt.Sprintf("%s = ?", columnName))
		}
		params = append(params, fieldVal.Interface())
		buf.AppendString(joinStr)
	}

	if buf.Buf == nil {
		return "", nil
	}
	whereStr := string(buf.Buf)
	return whereStr[:len(whereStr)-len(joinStr)], params
}
