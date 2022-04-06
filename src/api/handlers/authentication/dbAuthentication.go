package authentication

import (
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
)

// Basic user login fucntion
func userLogin(email string, password string) bool {
	// get password hash from database
	hash, err := dbInterface.RetrieveHashByEmail(email)

	// could not retrieve hash from database (email not in database)
	if err != nil {
		return false
	}

	// Check the password against the hash
	result := security.CheckPass(password, hash)

	return result
}

// Function that returns all cases that the user can view
func UserCases(userUUID string) []string {

	// Get the cases with the user's UUID
	cases, err := dbInterface.RetrieveViewCasesByUserUUID(userUUID)

	// Could not retrieve cases from database
	if err != nil {
		return nil
	}

	return cases
}
