package authentication

import (
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
)

// Basic user login fucntion
func userLogin(email string, password string) bool {
	// get password hash from database
	hash, err := dbInterface.RetrieveHashByEmail(email)

	if err != nil {
		return false
	}

	// Check the password against the hash
	result := security.CheckPass(password, hash)

	return result
}
