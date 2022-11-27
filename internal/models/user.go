package models

import (
	"time"

	"github.com/google/uuid"
)

// ToDo generate mock and tests.

// Authenticator describes contract for different types of authentication.
type Authenticator interface {
	CleanCredentials()
	SanitizeCredentials()
}

// User base model.
type User struct {
	UserUUID  uuid.UUID `db:"uuid" validate:"omitempty"`
	Auth      Authenticator
	CreatedAt time.Time `db:"created_at" exhaustruct:"optional"`
	UpdatedAt time.Time `db:"updated_at" exhaustruct:"optional"`
}

// NewUser creates User with credentials.
func NewUser(auth Authenticator) (usr *User) {
	auth.CleanCredentials()

	return &User{
		UserUUID: uuid.New(),
		Auth:     auth,
	}
}

// Sanitize removes sensitive info.
func (u *User) Sanitize() {
	u.Auth.SanitizeCredentials()
}
