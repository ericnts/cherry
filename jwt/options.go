package jwt

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/pem"
	"github.com/ericnts/config"
	"github.com/ericnts/log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func init() {
	options, err := config.Load[*Options]("jwt")
	if err != nil {
		log.Panicf("jwt配置加载失败, %v", err)
	}
	jwtOption = options
}

type Options struct {
	_key *rsa.PrivateKey
	_pub *rsa.PublicKey

	Issuer         string        `yaml:"issuer"`
	Key            string        `yaml:"key"`
	Pub            string        `yaml:"pub"`
	Expires        time.Duration `yaml:"expires"`
	RefreshExpires time.Duration `yaml:"refreshExpires"`
}

func (o *Options) GetKey() (*rsa.PrivateKey, error) {
	if o._key == nil {
		bs, err := base64.StdEncoding.DecodeString(o.Key)
		if err != nil {
			log.WithError(err).Error("加载jwt.key失败")
			return nil, err
		}
		key, err := jwt.ParseRSAPrivateKeyFromPEM(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: bs,
		}))
		if err != nil {
			log.WithError(err).Error("解析jwt.key失败")
			return nil, err
		}
		o._key = key
	}
	return o._key, nil
}

func (o *Options) GetPub() (*rsa.PublicKey, error) {
	if o._pub == nil {
		bs, err := base64.StdEncoding.DecodeString(o.Pub)
		if err != nil {
			log.WithError(err).Error("加载jwt.pub失败")
			return nil, err
		}
		pub, err := jwt.ParseRSAPublicKeyFromPEM(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: bs,
		}))
		if err != nil {
			log.WithError(err).Error("解析jwt.pub失败")
			return nil, err
		}
		o._pub = pub
	}

	return o._pub, nil
}
