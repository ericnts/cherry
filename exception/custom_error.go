package exception

import (
	"context"
	"github.com/ericnts/cherry/local"
	"google.golang.org/grpc/status"
)

type IError interface {
	Int() int
	Error() string
	Attachments() []error
	LocalKey() string
	GRPCStatus() *status.Status
}

type CustomError struct {
	ErrCode int
	ErrStr  string
	Errs    []error
}

func Custom(err Error, errStr string, errs ...error) CustomError {
	return CustomError{
		ErrCode: err.Int(),
		ErrStr:  errStr,
		Errs:    errs,
	}
}

func (e CustomError) Int() int {
	return e.ErrCode
}

func (e CustomError) Error() string {
	return local.Translate(context.TODO(), e.LocalKey())
}

func (e CustomError) Attachments() []error {
	return e.Errs
}

func (e CustomError) LocalKey() string {
	if len(e.ErrStr) > 0 {
		return e.ErrStr
	} else {
		return Error(e.ErrCode).LocalKey()
	}
}

func (e CustomError) GRPCStatus() *status.Status {
	return status.New(Error(e.ErrCode).Convert(), e.ErrStr)
}

func FromStatusErr(err error) CustomError {
	s, ok := status.FromError(err)
	if ok {
		return Custom(ConvertCode(s.Code()), s.Message(), err)
	} else {
		return Custom(Private, err.Error(), err)
	}
}
