package graph

import (
	"github.com/uptrace/bun"
	"vuegolang/pkg/sessions"
)

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Db       *bun.DB
	Sessions sessions.Interface
}
