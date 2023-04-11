package middleware

import (
	"log"
	"net/http"

	"vuegolang/graph/model"
	"vuegolang/pkg/ctxpayload"
)

type HandleAuthHTTPFunc func(w http.ResponseWriter, r *http.Request) (*http.Request, error)

func AuthMiddleware(HandleAuthHTTP HandleAuthHTTPFunc) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				var e error
				r, e = HandleAuthHTTP(w, r)
				if e != nil {
					http.Error(w, e.Error(), http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
			},
		)
	}
}

// Обрабатывает сессию клиента
func SessionHandleClient(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	var (
		payload = ctxpayload.Auth{ResponseWriter: w, Request: *r}
		ctx     = r.Context()
		auth    = model.AuthInfo{}
	)
	// Получим контекст
	cookie, errCookie := r.Cookie(ctxpayload.CookieName)
	if errCookie != nil {
		log.Printf("Cookie retrieval error: ", errCookie.Error())
		return r.WithContext(payload.WithContext(ctx)), nil
	}

	auth.Token = &cookie.Value

	payload.AuthInfo = auth
	ctx = payload.WithContext(ctx)

	return r.WithContext(ctx), nil
}

func HandleAuthHTTP(w http.ResponseWriter, r *http.Request) (*http.Request, error) {

	// Обработаем сессию клиента
	var e error
	r, e = SessionHandleClient(w, r)

	return r, e
}
