package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var defaultSecret = []byte("your_jwt_secret_key") // замени в проде

// Берём секрет из ENV, иначе дефолтный
func jwtSecret() []byte {
	sec := os.Getenv("JWT_SECRET")
	if sec == "" {
		return defaultSecret
	}
	return []byte(sec)
}

// GenerateJWT — выдаёт JWT с полем username и сроком 1 час
func GenerateJWT(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret())
}

// ParseJWT — возвращает username из токена
func ParseJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Проверим алгоритм
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret(), nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}
	u, _ := claims["username"].(string)
	if u == "" {
		return "", errors.New("username not found in token")
	}
	return u, nil
}
