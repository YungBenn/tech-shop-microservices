package token

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaim struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID string, email string) (tokenString string, err error) {
	JwtKey := []byte("secret")

	expirationTime := jwt.NewNumericDate(time.Now().Add(1 * time.Hour))
	claims := &JWTClaim{
		UserID:           userID,
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expirationTime},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(JwtKey)

	return
}
