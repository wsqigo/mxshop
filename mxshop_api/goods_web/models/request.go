package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	ID          uint
	NickName    string
	AuthorityId int64
	jwt.StandardClaims
}
