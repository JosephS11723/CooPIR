package dbInterface

import (
	"log"

	"github.com/JosephS11723/CooPIR/src/api/lib/security"
)

// Basic user login function
func UserLogin(email string, password string) bool {
	// get password hash from database
	hash, err := RetrieveHashByEmail(email)

	// could not retrieve hash from database (email not in database)
	if err != nil {
		log.Println("Failed to retrieve hash from database")
		return false
	}

	// Check the password against the hash
	result := security.CheckPass(password, hash)

	return result
}

// UserCases returns all cases that the user can view
func UserCases(userUUID string) []string {
	// Get the cases with the user's UUID
	cases, err := RetrieveViewCasesByUserUUID(userUUID)

	// Could not retrieve cases from database
	if err != nil {
		log.Println(err)
		return nil
	}

	return cases
}

// Returns true the user has responder rights
func UserResponderPermission(userUUID string) bool {
	// Find if user has supervisor rights
	result, err := FindResponderByUUID(userUUID)

	if err != nil {
		log.Println(err)
		return false
	}

	return result
}

// Returns true the user has supervisor rights
func UserSupervisorPermission(userUUID string) bool {
	// Find if user has supervisor rights
	result, err := FindSupervisorByUUID(userUUID)

	if err != nil {
		log.Println(err)
		return false
	}

	return result
}

// Returns true if the user has admin rights
func UserAdminPermission(userUUID string) bool {
	// Find if user has admin rights
	result, err := FindAdminByUUID(userUUID)

	if err != nil {
		log.Println(err)
		return false
	}

	return result
}
