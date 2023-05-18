package authpayload

import (
	"context"
	"errors"
	"net/http"

	"vuegolang/graph/model"
)

type CtxKey string

const ctxAuthPayloadKey CtxKey = "ctxAuthPayloadKey"
const CookieName string = "_auth"

type Auth struct {
	http.ResponseWriter
	http.Request
	model.AuthInfo
}

func (a *Auth) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxAuthPayloadKey, a)
}

func FromContext(ctx context.Context) *Auth {
	return ctx.Value(ctxAuthPayloadKey).(*Auth)
}

func GetSessionKeyFromContext(ctx context.Context) (string, error) {
	auth := FromContext(ctx)

	if auth.AuthInfo.Token == "" {
		return "", errors.New("no token")
	}

	return auth.AuthInfo.Token, nil
}
