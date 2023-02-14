package util

import (
	"crypto/md5"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func MakePWD(pwd string) ([]byte, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil, WrapErr(err, fmt.Sprintf("密码生成失败, %s", pwd))
	}
	return password, nil
}

func CheckPWD(pwd []byte, new []byte) error {
	return bcrypt.CompareHashAndPassword(pwd, new)
}

func MakeMD5(pwd string) string {
	data := []byte(pwd)
	has := md5.Sum(data)
	return fmt.Sprintf("%x", has)
}
