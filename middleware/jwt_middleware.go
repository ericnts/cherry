package middleware

import (
	"errors"
	"github.com/ericnts/cherry/current"
	"github.com/ericnts/cherry/exception"
	"github.com/ericnts/cherry/jwt"
	"github.com/ericnts/cherry/results"
	"github.com/ericnts/log"
	"github.com/gin-gonic/gin"
)

const (
	TokenKeyHead = "Authorization"
)

func JWTMiddleware(subjects ...jwt.Subject) func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(TokenKeyHead)
		if token == "" {
			results.Err(c, exception.TokenInvalid, errors.New("缺少Token信息"))
			return
		}
		jwtToken, err := jwt.ParseToken(token)
		if err != nil {
			log.With("token", token).Errorf("Token验证失败,%v", err)
			results.Err(c, err)
			return
		}
		current.SetToken(c, jwtToken)

		// 验证请求路径
		if jwtToken.IsRefresh {
			results.Err(c, exception.TokenInvalid)
			return
		}

		// 验证token类型 当token类型为机器人时，不进行访问路径校验
		if len(subjects) > 0 {
			var hit bool
			for i := range subjects {
				if subjects[i] == jwt.Subject(jwtToken.Subject) {
					hit = true
					break
				}
			}
			if !hit {
				results.Err(c, exception.Custom(exception.TokenInvalid, "没有权限进行此操作"))
				return
			}
		}
		c.Next()
	}
}
