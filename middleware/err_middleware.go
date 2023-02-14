package middleware

import (
	"github.com/ericnts/cherry/exception"
	"github.com/ericnts/cherry/results"
	"github.com/ericnts/log"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

func ErrMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("Panic info is: %v", err)
			log.Errorf("Panic info is: %s", debug.Stack())
			results.Err(c, exception.Private, err.(error))
		}
	}()
	c.Next()
	err := c.Errors.Last()
	if err == nil {
		return
	}
	switch err.Type {
	case gin.ErrorTypeBind:
		results.ParamErr(c, exception.ParamInvalid, err)
	default:
		results.Err(c, err)
	}
}
