package security

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password and returns the hexdump as a string
func HashPass(password string) (string, error) {
	// Generate a salted hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPass(password string, hash string) bool {
	// Check the password against the hash
	log.Println("PASS", password)
	log.Println("HASH", hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Println("[CheckPass ERROR}", err)
	}
	return err == nil // nil means no error
}
