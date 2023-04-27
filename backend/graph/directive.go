package graph

import (
	"context"
	"log"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"vuegolang/graph/model"
	"vuegolang/pkg/authpayload"
	"vuegolang/pkg/token"
)

func (r *Resolver) Auth(ctx context.Context, _ interface{}, next graphql.Resolver) (interface{}, error) {
	var (
		errorDefault = &gqlerror.Error{Message: "Доступ запрещён"}
	)

	auth := authpayload.FromContext(ctx)

	if auth.AuthInfo.Token == "" {
		log.Println("No token")
		return nil, errorDefault
	}

	sessionKey := auth.AuthInfo.Token

	tokenStr, errToken := token.FromJsonString(auth.Token)
	if errToken != nil {
		log.Println("Token from Json error: ", errToken)
		return nil, errorDefault
	}

	// Расшифровка
	session, okSession := r.Sessions.Get(sessionKey)
	if !okSession {
		log.Println("No session for key=", sessionKey)
		return nil, errorDefault
	}

	if session.Role != model.RoleAdmin {
		return nil, errorDefault
	}

	if errPkey := session.PrivateKey.Validate(); errPkey != nil {
		log.Println("Private key error: ", errPkey)
		return nil, errorDefault
	}

	claims, errDecrypt := token.DecryptBytes(tokenStr, session.PrivateKey)
	if errDecrypt != nil {
		log.Println("Decrypt error: ", errDecrypt)
		return nil, errorDefault
	}

	// TODO: делать редирект на логин когда токен протух
	if claims.Expired < time.Now().Unix() {
		log.Printf("Token expired for %v [%v]\n", claims.Id, time.Unix(claims.Expired, 0))
		return nil, &gqlerror.Error{Message: "Ваша сессия истекла. Пожалуйста, перелогиньтесь"}
	}

	return next(ctx)
}
