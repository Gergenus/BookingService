package jwtpkg

import "github.com/golang-jwt/jwt/v5"

type UserClaims struct {
	jwt.RegisteredClaims
	UUID     string
	Username string
	Role     string
	Email    string
}
