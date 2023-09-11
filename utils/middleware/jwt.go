package middleware

import (
	"errors"
	"fmt"
	"time"
	"wechat/utils"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(utils.JwtKey)

type JwtCustomClaims struct {
	ID               int
	Role             int
	RegisteredClaims jwt.RegisteredClaims
}

func (j JwtCustomClaims) Valid() error {
	return nil
}

// 生成token
func GenerateToken(uid int, role int) (string, error) {
	//自定义数据
	myJwtCustomClaims := JwtCustomClaims{
		ID:   uid,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(utils.TokenTime) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myJwtCustomClaims)
	return token.SignedString(jwtKey)
}

// 解析token
func ParseToken(tokenStr string) (JwtCustomClaims, error) {
	myJwtCustomClaims := JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, &myJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err == nil && !token.Valid {
		err = errors.New("invalid Token error! [47]")
	}

	return myJwtCustomClaims, err
}

// token校验
func IsTokenValid(token string) bool {
	// myJwtCustomClaims, err:= ParseToken(token)
	_, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// 判断用户角色
func RoleJuge(token string, isRole int) bool {
	myJwtCustomClaims, err := ParseToken(token)
	// _, err := ParseToken(token)
	if err != nil {
		fmt.Println(err)
		return false
	} else if myJwtCustomClaims.Role != isRole {
		return false
	}
	return true
}
