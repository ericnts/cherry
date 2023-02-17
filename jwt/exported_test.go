package jwt

import (
	"encoding/base64"
	"encoding/pem"
	"github.com/ericnts/log"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"testing"
	"time"
)

func TestKey(t *testing.T) {
	keyData, _ := ioutil.ReadFile("../conf/key")
	b, _ := pem.Decode(keyData)
	s := base64.StdEncoding.EncodeToString(b.Bytes)
	log.Info(s)

	keyBs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	KeyBs := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBs,
	})
	ioutil.WriteFile("../conf/key2", KeyBs, fs.ModePerm)
}

func TestPub(t *testing.T) {
	keyData, _ := ioutil.ReadFile("../conf/key.pub")
	b, _ := pem.Decode(keyData)
	s := base64.StdEncoding.EncodeToString(b.Bytes)
	log.Info(s)

	keyBs, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	KeyBs := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: keyBs,
	})

	ioutil.WriteFile("../conf/key.pub2", KeyBs, fs.ModePerm)
}

func TestGenToken(t *testing.T) {
	id := "1234"
	s, err := GenToken(id, "", User, nil, false, false)
	assert.Nil(t, err)
	cc, err2 := ParseToken(s)
	assert.Nil(t, err2)
	assert.Equal(t, id, cc.Id)
}

func TestNormal(t *testing.T) {
	println(time.Now().UnixMilli())
	println(time.Second.Milliseconds())
}
