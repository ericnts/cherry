package results

import (
	"context"
	"github.com/ericnts/cherry/exception"
	"github.com/ericnts/cherry/local"
	"github.com/ericnts/cherry/util"
)

type Result struct {
	Code    int         `json:"code"`            //错误编码
	Data    interface{} `json:"data"`            //数据
	Page    *Page       `json:"page,omitempty"`  //分页信息
	Message string      `json:"message"`         //描述信息
	Error   string      `json:"error,omitempty"` //错误信息
}

type Page struct {
	No    int   `json:"no"`
	Size  int   `json:"size"`
	Count int64 `json:"count"`
}

func NewPage(no, size int, count int64) *Page {
	return &Page{
		No:    no,
		Size:  size,
		Count: count,
	}
}

func NewResult(c context.Context, data interface{}, page *Page, errs ...error) *Result {
	var customErr exception.IError
	customErr = exception.Ok
	var finalErr error
	if len(errs) > 0 {
		customErr = exception.Private
		for i := range errs {
			if errs[i] == nil {
				continue
			}
			e, ok := errs[i].(exception.IError)
			if ok {
				customErr = e
				for _, ce := range e.Attachments() {
					if ce == nil {
						continue
					}
					finalErr = util.WrapErr(finalErr, ce.Error())
				}
			} else {
				finalErr = util.WrapErr(finalErr, errs[i].Error())
			}
		}
	}
	result := &Result{
		Code:    customErr.Int(),
		Data:    data,
		Page:    page,
		Message: local.Translate(c, customErr.LocalKey()),
	}
	if finalErr != nil {
		result.Error = finalErr.Error()
	}
	return result
}

func (r Result) Ok() bool {
	return r.Code == exception.Ok.Int()
}
