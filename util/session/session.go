package session

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/cs50-romain/Metis/util/task"
)

type Session struct {
	UserID		int
	UserName	string
	ExpiresAt	time.Time
	UserData	*task.Tasks
}

var sessions = make(map[string]Session)

func Print() {
	fmt.Println("Sessions:", sessions)
}

func CreateSession(username string) string {
	userid := len(sessions)
	sessionId := generateSessionId()

	newSession := Session{
		UserID: userid, 
		UserName: username,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 10),
		UserData: &task.Tasks{},
	}

	sessions[sessionId] = newSession
	return sessionId
}

func GetSession(sessionid string) (*Session, bool){
	session, ok := sessions[sessionid]
	if !ok {
		return nil, ok
	}

	return &session, true 
}

func DeleteSession(sessionid string) {
	delete(sessions, sessionid)
}

func generateSessionId() string {
	return uuid.New().String()
}
