package security

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPass(password string) string {
	// Generate a salted hash of the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panicln(err)
	}

	return string(hash)
}
