package types

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type User struct {
	Name         string `json:"username"`
	PasswordHash string `json:"passwordhash"`
}

type Token struct {
	Name string
	*jwt.StandardClaims
}
