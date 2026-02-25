package utils

import (
	"time"

	"fiber.com/session-api/config"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, userName, email, role string) (string, error) {
	expiresAt := time.Now().Add(time.Duration(config.AppConfig.JWTExpiresHour) * time.Hour)

	claims := JWTClaims{
		UserID:   userID,
		UserName: userName,
		Email:    email,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func ValidateToken(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	return claims, nil
}
