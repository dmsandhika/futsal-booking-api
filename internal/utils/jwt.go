package utils

import (
	"time"
	"os"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(adminID int)(string, error){
	claims := jwt.MapClaims{
		"admin_id": adminID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}