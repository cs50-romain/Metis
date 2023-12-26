package session

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Username	string
	ExpiresAt	time.Time
}

type Store struct {
	Sessions map[string]Session
}

func Init() *Store{
	sessions := make(map[string]Session)
	store := &Store{sessions}
	return store
}

func (st *Store) Get(session_id string) Session {
	return st.Sessions[session_id]
}

// This function saves to map and then call save() to save to file
func (st *Store) Save(username string) string {
	unique_id := generateSessionId()	
	session := Session{
		Username: username,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 10),
	}
	st.Sessions[unique_id] = session
	return unique_id
}

func generateSessionId() string {
	return uuid.New().String()
}
