package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"
	"net/mail"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

// salt it are the salt for the passwords of the users
// it mus be change later for enviroment variable
var salt string = "tbu89L&t3cyMez52p3BuS$"

// IsEmail validate if the parameter string are a email
// It returns true if it's a email and false if its no
func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)

	return err == nil
}

// GetHash generate a hash for the string paramater
// and return a hash string created with encrypt sha256
func GetHash(s string) string {
	hash := sha256.Sum256([]byte(s + salt))
	return hex.EncodeToString(hash[:])
}

// CreateJWT create a jwt with a Claims what resive, and
// with the key generate a JWT and return it
func CreateJWT(claims jwt.Claims) (string, error) {
	key := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(key))
	if err != nil {
		return t, err
	}
	return t, nil
}

// loadENV get to the so variable, all the enviroments variables
// of the .env file
func LoadENV() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}
