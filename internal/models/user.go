package models

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User base model
type User struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty"`
	Email     string    `json:"email" db:"email" validate:"omitempty,lte=60,email"`
	Password  string    `json:"password,omitempty" db:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

// SanitizePassword removes password
func (u *User) SanitizePassword() {
	u.Password = ""
}

// HashPassword hashes password without salt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword(u.saltPassword([]byte(u.Password)), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// ValidatePassword compares user password and payload
func (u *User) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), u.saltPassword([]byte(password)))
}

// PrepareCreate cleans credentials
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := u.HashPassword(); err != nil {
		return err
	}

	return nil
}

// saltPassword mixes password with salt
// ToDo salt must be added
func (u *User) saltPassword(password []byte) []byte {
	return password
}
