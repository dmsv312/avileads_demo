package utils

import "sync"

type Session struct {
	History []map[string]string
}

type SessionManager struct {
	sessions map[string]*Session
	mu       sync.Mutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (sm *SessionManager) GetOrCreateSession(sessionID string) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[sessionID]; exists {
		return session
	}

	newSession := &Session{History: []map[string]string{}}
	sm.sessions[sessionID] = newSession
	return newSession
}

func (session *Session) AddMessage(role, content string) {
	session.History = append(session.History, map[string]string{
		"role":    role,
		"content": content,
	})
}
