package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net/mail"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var salt string = "tbu89L&t3cyMez52p3BuS$"

func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)

	return err == nil
}

func GetHash(s string) string {
	hash := sha256.Sum256([]byte(s + salt))
	return hex.EncodeToString(hash[:])
}

func CreateJWT(claims jwt.Claims) (string, error) {
	key := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return t, err
	}
	return t, nil
}
