package token

import (
	"context"

	"github.com/dgrijalva/jwt-go"
)

type jwtKey struct{}

var jwtValKey = jwtKey{}

func SetContext(c context.Context, value interface{}) context.Context {
	return context.WithValue(c, jwtValKey, value)
}

func FromContext(c context.Context) jwt.Claims {
	return c.Value(jwtValKey).(jwt.Claims)
}
