package utils

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// SessionStore stores active session information
// In a production system, this would be stored in Redis or a database
var (
	activeSessions      = make(map[string]sessionData)
	activeSessionsMutex sync.RWMutex
)

type sessionData struct {
	UserID    uint
	ExpiresAt time.Time
}

// GenerateToken creates a new session ID for the given user ID
func GenerateToken(userID uint) (string, error) {
	// Generate a random session ID
	b := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	
	sessionID := base64.URLEncoding.EncodeToString(b)

	// Parse expiration time from environment variable or use default (24h)
	expiresInStr := os.Getenv("JWT_EXPIRES_IN")
	if expiresInStr == "" {
		expiresInStr = "24h"
	}

	// Parse expiration duration
	expiresIn, err := time.ParseDuration(expiresInStr)
	if err != nil {
		return "", fmt.Errorf("invalid JWT_EXPIRES_IN value: %v", err)
	}

	// Store session data
	activeSessionsMutex.Lock()
	activeSessions[sessionID] = sessionData{
		UserID:    userID,
		ExpiresAt: time.Now().Add(expiresIn),
	}
	activeSessionsMutex.Unlock()

	return sessionID, nil
}

// ValidateToken validates a session ID and returns the user ID if valid
func ValidateToken(sessionID string) (uint, error) {
	if sessionID == "" {
		return 0, errors.New("empty session ID")
	}

	activeSessionsMutex.RLock()
	session, exists := activeSessions[sessionID]
	activeSessionsMutex.RUnlock()

	if !exists {
		return 0, errors.New("invalid session")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Clean up expired session
		activeSessionsMutex.Lock()
		delete(activeSessions, sessionID)
		activeSessionsMutex.Unlock()
		
		return 0, errors.New("session expired")
	}

	return session.UserID, nil
}

// CleanupSessions removes expired sessions (can be called periodically)
func CleanupSessions() {
	now := time.Now()
	
	activeSessionsMutex.Lock()
	defer activeSessionsMutex.Unlock()
	
	for id, session := range activeSessions {
		if now.After(session.ExpiresAt) {
			delete(activeSessions, id)
		}
	}
}