package authentication

import (
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
)

// Basic user login fucntion
func UserLogin(email string, password string) bool {

	hash := dbInterface.RetrieveHashByEmail(email)

	// Check the password against the hash
	result := security.CheckPass(password, hash)

	return result
}
