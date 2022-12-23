package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// ToDo generate mock and tests.

// Authenticator describes contract for different types of authentication.
type Authenticator interface {
	CleanCredentials()
}

// User base model.
type User struct {
	UserUUID  uuid.UUID `db:"uuid" validate:"omitempty"`
	Auth      Authenticator
	CreatedAt time.Time `db:"created_at" exhaustruct:"optional"`
	UpdatedAt time.Time `db:"updated_at" exhaustruct:"optional"`
}

// MarshalZerologObject gives context to zerlog logs.
func (u User) MarshalZerologObject(e *zerolog.Event) {
	e.Stringer("uuid", u.UserUUID)
}

// NewUser creates User with credentials.
func NewUser(auth Authenticator) (usr *User) {
	auth.CleanCredentials()

	return &User{
		UserUUID: uuid.New(),
		Auth:     auth,
	}
}
