package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var (
	ErrUnexpectedSignedMeth = fmt.Errorf("unexpected token signing method")
	ErrInvalidTokenClaims   = fmt.Errorf("invalid token claims")
)

// JWTManager is a JSON web token manager.
type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// UserClaims is a custom JWT claims that contains some user's information.
type UserClaims struct {
	jwt.RegisteredClaims
	UserUUID uuid.NullUUID `json:"userUUID"`
}

// NewJWTManager returns a new JWT manager.
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Generate generates and signs a new token for a user.
func (manager *JWTManager) Generate(usrUUID uuid.UUID) (token string, err error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(manager.tokenDuration)},
		},
		UserUUID: uuid.NullUUID{
			UUID:  usrUUID,
			Valid: true,
		},
	}

	claimedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return claimedToken.SignedString([]byte(manager.secretKey)) //nolint:wrapcheck
}

// Verify verifies the access token string and return a user claim if the token is valid.
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, ErrUnexpectedSignedMeth
			}

			return []byte(manager.secretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, ErrInvalidTokenClaims
	}

	return claims, nil
}
