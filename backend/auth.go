package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"vuegolang/graph/model"
	"vuegolang/pkg/ctxpayload"
	"vuegolang/pkg/middleware"
)

// Обрабатывает сессию клиента
func SessionHandleClient(w http.ResponseWriter, r *http.Request) (*http.Request, error) {
	var auth = model.AuthInfo{}

	// Получим контекст
	ctx := r.Context()

	cookie, errCookie := r.Cookie(ctxpayload.CookieName)
	if errCookie != nil {
		log.Printf(errCookie.Error())
	}

	if errJson := json.Unmarshal([]byte(cookie.Value), &auth); errJson != nil {
		log.Println(errJson)
		return r, errors.New("wrong cookie")
	}
	s := ctxpayload.ContextPayload{ResponseWriter: w, AuthInfo: auth}
	ctx = context.WithValue(ctx, middleware.CtxAuthPayloadKey, s)

	// Сюда запишем сессию, если сработает кейс
	// var sess *model.Session
	// var ClientID string

	// Проверим наличие cookie c ClientID
	// cookie, err := r.Cookie("_cid")
	// if err == nil {
	// 	// ClientID = cookie.Value
	//
	// 	// У клиента есть ClientID
	// 	// 1. Проверим наличие заголовка Session-ID
	// 	// 2. Получаем токен и валидируем его
	// 	// 2.1. Токен валидный: сохраним сессию из токена
	// 	// 2.2. Токен протух: сохраним сессию из токена
	// 	// 2.3. Токен Invalid: создадим новую сессии
	//
	// 	// Ищем заголовок Session-ID
	// 	if t := r.Header.Get("Session-ID"); t != "" {
	//
	// 		// Нашли сессию
	// 		log.Println("Implement validation")
	// 		// if ss, err2 := s.ValidateSessionToken(t); err2 == nil {
	// 		// 	sess = ss
	// 		// }
	// 	}
	//
	// 	// Этот метод теперь удален
	// 	// sess = model.NewSessionWithSid(cookie.Value)
	// }

	// Если сессии нет: создаем сессию
	// if sess == nil {
	//
	// 	sess = model.NewSession()
	//
	// 	if ClientID != "" {
	// 		sess.AddClientID(ClientID)
	// 	}
	// }

	// Если есть ошибка при чтении cookie
	// if err != nil {
	//
	// 	// Получим ID клиента
	// 	log.Println("Implement cookie set")
	// 	// cid, err2 := sess.GetSid()
	// 	// if err2 != nil {
	// 	// 	fmt.Printf(err.Error())
	// 	// 	return r
	// 	// }
	//
	// 	// Создадим cookie
	// 	cookie = &http.Cookie{
	// 		Name:     "_auth",
	// 		Value:    "test",
	// 		Expires:  time.Now().AddDate(1, 0, 0),
	// 		Path:     "/",
	// 		SameSite: http.SameSiteStrictMode,
	// 	}
	// 	// Установим cookie
	// 	http.SetCookie(w, cookie)
	// }

	// Сохраним сессию в контекст и вернем *http.Request
	return r.WithContext(ctx), nil // r.WithContext(sess.WithContext(ctx))
}

func HandleAuthHTTP(w http.ResponseWriter, r *http.Request) (*http.Request, error) {

	// Обработаем сессию клиента
	var e error
	r, e = SessionHandleClient(w, r)

	return r, e
}
