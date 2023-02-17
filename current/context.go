package current

import (
	"context"
	"github.com/ericnts/cherry/jwt"
	"github.com/gin-gonic/gin"
)

const (
	TokenKey    = "Authorization"
	UserIDKey   = "currentUserID"
	UserNameKey = "currentUserName"
	UserTypeKey = "currentUserType"
	UserKey     = "currentUser"
	IDKey       = "id"
	GroupIDKey  = "groupID"
	OfficeIDKey = "officeID"
)

func SetToken(ctx context.Context, token *jwt.CustomClaims) context.Context {
	if ctx == nil {
		return ctx
	}
	if c, ok := ctx.(*gin.Context); ok {
		c.Set(TokenKey, token)
	} else {
		ctx = context.WithValue(ctx, TokenKey, token)
	}
	SetUserID(ctx, token.Id)
	SetUserName(ctx, token.Audience)
	SetUserType(ctx, token.Subject)
	return ctx
}

func Token(ctx context.Context) *jwt.CustomClaims {
	if ctx == nil {
		return nil
	}
	token := ctx.Value(TokenKey)
	if token == nil {
		return nil
	}
	return token.(*jwt.CustomClaims)
}

func SetUserID(c context.Context, userID string) context.Context {
	return SetValue(c, UserIDKey, userID)
}

func UserID(c context.Context) string {
	return GetString(c, UserIDKey)
}

func SetUserName(c context.Context, userName string) context.Context {
	return SetValue(c, UserNameKey, userName)
}

func UserName(c context.Context) string {
	return GetString(c, UserNameKey)
}

func SetUserType(c context.Context, userType string) context.Context {
	return SetValue(c, UserTypeKey, userType)
}

func UserType(c context.Context) string {
	return GetString(c, UserTypeKey)
}

func SetUser(ctx context.Context, user interface{}) context.Context {
	return SetValue(ctx, UserKey, user)
}

func User(ctx context.Context) interface{} {
	if ctx == nil {
		return nil
	}
	return ctx.Value(UserKey)
}

func SetID(ctx context.Context, dataID string) context.Context {
	return SetValue(ctx, IDKey, dataID)
}

func ID(ctx context.Context) string {
	return GetString(ctx, IDKey)
}

func SetGroupID(ctx context.Context, dataID string) context.Context {
	return SetValue(ctx, GroupIDKey, dataID)
}

func GroupID(ctx context.Context) string {
	return GetString(ctx, GroupIDKey)
}

func SetOfficeID(ctx context.Context, dataID string) context.Context {
	return SetValue(ctx, OfficeIDKey, dataID)
}

func OfficeID(ctx context.Context) string {
	return GetString(ctx, OfficeIDKey)
}

func SetValue(ctx context.Context, key string, value interface{}) context.Context {
	if ctx == nil {
		return ctx
	}
	if c, ok := ctx.(*gin.Context); ok {
		c.Set(key, value)
	} else {
		ctx = context.WithValue(ctx, key, value)
	}
	return ctx
}

func GetInterface(ctx context.Context, key string) interface{} {
	if ctx == nil {
		return nil
	}
	return ctx.Value(key)
}

func GetString(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	id := ctx.Value(key)
	if id == nil {
		return ""
	}
	return id.(string)
}

func GetStringArray(ctx context.Context, key string) []string {
	if ctx == nil {
		return nil
	}
	value := ctx.Value(key)
	if value == nil {
		return nil
	}
	return value.([]string)
}
