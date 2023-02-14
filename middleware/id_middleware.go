package middleware

import (
	"github.com/ericnts/cherry/current"
	"github.com/ericnts/cherry/exception"
	"github.com/ericnts/cherry/results"
	"github.com/gin-gonic/gin"
	"strings"
)

func IDMiddleware(c *gin.Context) {
	groupID, ok := c.Params.Get(current.GroupIDKey)
	if ok {
		if len(strings.TrimSpace(groupID)) == 0 {
			results.Err(c, exception.Custom(exception.ParamInvalid, "集团ID不能为空"))
			return
		}
		current.SetGroupID(c, groupID)
	}
	officeID, ok := c.Params.Get(current.OfficeIDKey)
	if ok {
		if len(strings.TrimSpace(officeID)) == 0 {
			results.Err(c, exception.Custom(exception.ParamInvalid, "机构ID不能为空"))
			return
		}
		current.SetOfficeID(c, officeID)
	}
	dataID, ok := c.Params.Get(current.IDKey)
	if ok {
		if len(strings.TrimSpace(dataID)) == 0 {
			results.Err(c, exception.Custom(exception.ParamInvalid, "ID不能为空"))
			return
		}
		current.SetID(c, dataID)
	}
	c.Next()
}
