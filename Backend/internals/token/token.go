package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims holds JWT access token claims
type AccessClaims struct {
	UserID uint
	Email  string
	Role   string
	jwt.RegisteredClaims
}

// RefreshClaims holds JWT refresh token claims
type RefreshClaims struct {
	UserID uint
	Role string
	JTI    string
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a JWT access token (15 mins)
func GenerateAccessToken(userID uint, email, role string) (string, error) {
	claims := AccessClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			Issuer:    "insight-api",
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(os.Getenv("JWT_ACCESS_SECRET")))
}

// GenerateRefreshToken generates a JWT refresh token (7 days) with a unique JTI
func GenerateRefreshToken(userID uint,role, jti string) (string, error) {
	claims := RefreshClaims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(os.Getenv("JWT_REFRESH_SECRET")))
}
