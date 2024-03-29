package handler

import "github.com/golang-jwt/jwt"

func NewWithClaims(claims Claims) *jwt.Token {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func From(v interface{}) Claims {
	return *v.(*jwt.Token).Claims.(*Claims)
}

type Claims struct {
	jwt.StandardClaims
	Username string `json:"usr,omitempty"`
	ID       int64  `json:"id,omitempty"`
}
