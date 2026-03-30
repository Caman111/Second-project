package auth

import (
	"fmt"
	"sync"
	"time"
)

type Session struct {
	Phone string
	Code  string
}

type AuthRepository struct {
	mu           sync.RWMutex
	sessions     map[string]Session
	lastRequests map[string]time.Time
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{
		sessions:     make(map[string]Session),
		lastRequests: make(map[string]time.Time),
	}
}

func (r *AuthRepository) CheckLimit(phone string) error {
	r.mu.RLock()
	lastTime, ok := r.lastRequests[phone]
	r.mu.RUnlock()

	if ok && time.Since(lastTime) < 1*time.Minute {
		secondsLeft := int(60 - time.Since(lastTime).Seconds())
		return fmt.Errorf("подождите %d секунд", secondsLeft)
	}
	return nil
}

func (r *AuthRepository) UpdateLimit(phone string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastRequests[phone] = time.Now()
}

func (r *AuthRepository) SaveSession(sessionId string, session Session) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sessions[sessionId] = session
}

func (r *AuthRepository) GetSession(sessionId string) (Session, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	session, ok := r.sessions[sessionId]
	if ok {
		delete(r.sessions, sessionId)
	}
	return session, ok
}

func (r *AuthRepository) DeleteSession(id string) {
    r.mu.Lock()
    defer r.mu.Lock()
    delete(r.sessions, id)
}