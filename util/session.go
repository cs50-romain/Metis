package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	UserID		int
	UserName	string
	ExpiresAt	time.Time
	UserData	[]any
}

var sessions = make(map[uuid.UUID]Session)

func CreateSession(userid int, username string) string {
	sessionId := generateSessionId()

	newSession := Session{
		UserID: userid,
		UserName: username,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 10),
		UserData: []any{},
	}

	sessions[sessionId] = newSession
	return sessionId.String()
}

func GetSession(sessionid uuid.UUID) (Session, bool){
	session, ok := sessions[sessionid]
	return session, ok
}

func DeleteSession(sessionid uuid.UUID) {
	delete(sessions, sessionid)
}

func generateSessionId() uuid.UUID {
	return uuid.New()
}
