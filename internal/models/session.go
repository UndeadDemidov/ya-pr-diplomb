package models

import "github.com/google/uuid"

// Session model.
type Session struct {
	SessionID string    `json:"sessionID"`
	UserID    uuid.UUID `json:"userID"`
}
