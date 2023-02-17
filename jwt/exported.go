package jwt

import (
	"fmt"
	"github.com/ericnts/cherry/exception"
	"github.com/ericnts/log"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Subject string

const (
	User      Subject = "user"
	Service   Subject = "service"
	ThirdUSer Subject = "thirdUser"
	Proxy     Subject = "proxy"
)

var (
	jwtOption   *Options
	SingleCheck = &memorySingleCheck{
		data: make(map[string]int64, 10),
	}
)

type CustomClaims struct {
	jwt.StandardClaims

	Raw       string
	IsRefresh bool     //刷新token
	Single    bool     //单节点登录
	Scopes    []string //服务范围
}

func GenRefreshToken(id, name string, subject Subject, scopes []string, single bool) (string, error) {
	return GenToken(id, name, subject, scopes, true, single)
}

func GenSingleToken(id, name string, subject Subject, scopes []string) (string, error) {
	return GenToken(id, name, subject, scopes, false, true)
}

func GenToken(id, name string, subject Subject, scopes []string, refresh bool, single bool) (string, error) {
	expires := jwtOption.Expires
	if refresh {
		expires = jwtOption.RefreshExpires
	}
	now := time.Now().Unix()
	data := CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        id,                             //id
			Issuer:    jwtOption.Issuer,               //签发者
			NotBefore: now,                            //生效时间
			IssuedAt:  now,                            //签发时间
			ExpiresAt: time.Now().Add(expires).Unix(), //过期时间
			Audience:  name,
			Subject:   string(subject),
		},
		IsRefresh: refresh, //刷新token
		Single:    single,  //单节点登录
		Scopes:    scopes,
	}
	SingleCheck.Logon(fmt.Sprintf("auth:%s:%s", subject, id), now)
	key, err := jwtOption.GetKey()
	if err != nil {
		return "", err
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, data).SignedString(key)
	return fmt.Sprintf("Bearer %s", token), err
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	tokenString = strings.Replace(tokenString, "%20", " ", 1)
	if index := strings.Index(tokenString, " "); index >= 0 {
		tokenString = tokenString[index+1:]
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtOption.GetPub()
	})
	if err != nil {
		log.WithError(err).Warn("token解析失败")
		return nil, exception.Custom(exception.TokenInvalid, err.Error())
	}
	if !token.Valid {
		return nil, exception.TokenInvalid
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, exception.TokenInvalid
	}
	claims.Raw = tokenString
	if claims.Single {
		if !SingleCheck.Check(fmt.Sprintf("auth:%s:%s", claims.Audience, claims.Id), claims.IssuedAt) {
			return nil, exception.TokenReplaced
		}
	}
	return claims, nil
}
