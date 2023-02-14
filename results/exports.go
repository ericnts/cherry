package results

import (
	"github.com/ericnts/cherry/base"
	"github.com/ericnts/cherry/exception"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(c *gin.Context, data interface{}, page ...*Page) {
	if c.IsAborted() {
		return
	}
	if len(page) > 0 {
		c.JSON(http.StatusOK, NewResult(c, data, page[0], exception.Ok))
	} else {
		c.JSON(http.StatusOK, NewResult(c, data, nil, exception.Ok))
	}
}

func Success2(c *gin.Context, data interface{}) {
	result, ok := data.(*base.CollectionResult)
	if ok {
		Success(c, result.List, NewPage(result.Page.PageNO, result.Page.PageSize, result.Page.Count))
	} else {
		Success(c, data)
	}
}

func Err(c *gin.Context, errs ...error) {
	if c.IsAborted() {
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, NewResult(c, nil, nil, errs...))
}

func ParamErr(c *gin.Context, errs ...error) {
	c.AbortWithStatusJSON(http.StatusOK, NewResult(c, nil, nil, errs...))
}

func Warn(c *gin.Context, data interface{}, errs ...error) {
	WarnWithPage(c, data, nil, errs...)
}

func WarnWithPage(c *gin.Context, data interface{}, page *Page, errs ...error) {
	if c.IsAborted() {
		return
	}
	c.AbortWithStatusJSON(http.StatusOK, NewResult(c, data, page, errs...))
}
