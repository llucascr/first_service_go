package model

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	JWT_SECRET_KEY = "JWT_SECRET_KEY"
)

func (u *User) GenerateToken() (string, error) {
	secret := os.Getenv(JWT_SECRET_KEY)
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not configured")
	}

	claims := jwt.MapClaims{
		"user_id": u.UserID,
		"name":    u.Name,
		"email":   u.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token_string, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token_string, nil
}

func ValidateToken(token_string string) (string, error) {
	secret := os.Getenv(JWT_SECRET_KEY)
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET_KEY not configured")
	}

	token, err := jwt.Parse(token_string, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	name, ok := claims["name"].(string)
	if !ok {
		return "", fmt.Errorf("name claim missing")
	}

	return name, nil
}
