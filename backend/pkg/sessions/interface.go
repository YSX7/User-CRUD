package sessions

type Interface interface {
	Add(tokenStr string, session *Session)
	Get(tokenStr string) (Session, bool)
	Remove(tokenStr string)
}
