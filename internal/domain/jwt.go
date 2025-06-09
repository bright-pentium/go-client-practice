package domain

import "github.com/golang-jwt/jwt/v4"

type JwtClaims struct {
	Name  string  `json:"name,omitempty"`
	Scope string  `json:"scope,omitempty"`
	Type  JwtType `json:"type,omitempty" `
	jwt.RegisteredClaims
}

type JwtType string

const (
	UserType   JwtType = "user"
	ClientType JwtType = "client"
)
