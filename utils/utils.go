package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"net/mail"
)

var addition string = "tbu89L&t3cyMez52p3BuS$"

func IsEmail(s string) bool {
	_, err := mail.ParseAddress(s)

	return err == nil
}

func GetHash(s string) string {
	hash := sha256.Sum256([]byte(s + addition))
	return hex.EncodeToString(hash[:])
}
