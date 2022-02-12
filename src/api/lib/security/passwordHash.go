package security

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) string {
	// Generate a salted hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}

	return string(hash)
}

func CheckPass(password string, hash string) bool {
	// Check the password against the hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // nil means no error
}
