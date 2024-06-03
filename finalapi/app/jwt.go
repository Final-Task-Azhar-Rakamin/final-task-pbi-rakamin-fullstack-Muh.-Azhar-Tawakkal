package app

import "github.com/golang-jwt/jwt/v4"

var JWT_KEY = []byte("uY2rBGxXVGZqlC8o1gJwz14kBhZ9u/VQypB/oL0+q3s=")

type JWTClaim struct {
	Id    int
	Email string
	jwt.RegisteredClaims
}
