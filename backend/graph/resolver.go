package graph

import (
	"crypto/rsa"

	"github.com/uptrace/bun"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Session struct {
	PrivateKey rsa.PrivateKey
}

type Sessions map[string]Session

type Resolver struct {
	Db *bun.DB
	Sessions
}
