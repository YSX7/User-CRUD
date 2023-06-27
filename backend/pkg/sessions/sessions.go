package sessions

import (
	"crypto/rsa"

	"vuegolang/graph/model"
)

type Session struct {
	Id    int
	Login string
	model.Role
	PrivateKey rsa.PrivateKey
}

type sessions map[string]Session

func (s sessions) Remove(tokenStr string) {
	delete(s, tokenStr)
}

func (s sessions) Get(tokenStr string) (Session, bool) {
	value, ok := s[tokenStr]
	return value, ok
}

func (s sessions) Add(tokenStr string, session *Session) {
	s[tokenStr] = *session
}

func New() Interface {
	return &sessions{}
}
