package token

import (
	"time"

	"github.com/YungBenn/tech-shop-microservices/internal/auth/entity"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaim struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateJWT(user entity.User) (string, *JWTClaim, error) {
	JwtKey := []byte("secret")

	claims := &JWTClaim{
		UserID:           user.ID,
		Email:            user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)

	return tokenString, claims, err
}

func VerifyJWT(token string) (*JWTClaim, error) {
	JwtKey := []byte("secret")

	claims := &JWTClaim{}
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, jwt.ErrInvalidKey
		
		}
		return JwtKey, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, claims, keyFunc)
	if err != nil {
		return nil, jwt.ErrInvalidKey
	}

	payload, ok := jwtToken.Claims.(*JWTClaim)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return payload, nil
}