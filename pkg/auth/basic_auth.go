package auth

import (
	"fmt"
	"strings"

	"github.com/UndeadDemidov/ya-pr-diplomb/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

var _ models.Authenticator = (*BasicAuth)(nil)

// BasicAuth implements basic (login-password) authentication strategy.
type BasicAuth struct {
	Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Password string `json:"password,omitempty" db:"password"`
}

func NewBasicAuth(email, password string) *BasicAuth {
	ba := BasicAuth{email, password}
	ba.CleanCredentials()

	return &ba
}

// HashPassword hashes password without salt.
func (ba *BasicAuth) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(ba.saltPassword([]byte(ba.Password)), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to generate hash from password: %w", err)
	}

	ba.Password = string(hashedPassword)

	return nil
}

// ValidatePassword compares user password and payload.
func (ba *BasicAuth) ValidatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(ba.Password), ba.saltPassword([]byte(password))); err != nil {
		return fmt.Errorf("failed to compare hash with password: %w", err)
	}

	return nil
}

// CleanCredentials cleans credentials.
func (ba *BasicAuth) CleanCredentials() {
	ba.Email = strings.ToLower(strings.TrimSpace(ba.Email))
	ba.Password = strings.TrimSpace(ba.Password)
}

// saltPassword mixes password with salt.
// ToDo salt must be added.
func (ba *BasicAuth) saltPassword(password []byte) []byte {
	return password
}

func (ba *BasicAuth) MarshalZerologObject(e *zerolog.Event) {
	e.Str("email", ba.Email)
}
