package session

import (
	"fmt"
	"encoding/json"
	"os"
	"log"	
	"time"

	"github.com/google/uuid"
	"github.com/cs50-romain/Metis/util/task"
)

const file_saving_location = "./session-data/"

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

func PreFillSessions() {
	fmt.Println("Pre filling sessions map")
	b, err := os.ReadFile("./session-data/session-data.json")
	if err != nil {
		log.Println("[ERROR] PreFillSessions ReadFile error -> ", err)
	}

	err = json.Unmarshal([]byte(b), &sessions)
	if err != nil {
		log.Println("[ERROR] PreFillSessions() -> ", err)
	}
}

func Save() {
	b, err := json.Marshal(sessions)
	if err != nil {
		log.Print("[INFO] Couldn't save file, retrying in...")
	} else {
		log.Println("[INFO] Session data File saved")
	}

	err = os.WriteFile(file_saving_location + "/session-data.json", b, 0644)
	if err != nil {
		log.Print("[ERROR] Error saving to file -> ", err)
	}
}
