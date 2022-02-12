package security

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password and returns the hexdump as a string
func HashPass(password string) string {
	// Generate a salted hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(hash))

	return string(hash)
}

func CheckPass(password string, hash string) bool {
	// Check the password against the hash
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil // nil means no error
}
