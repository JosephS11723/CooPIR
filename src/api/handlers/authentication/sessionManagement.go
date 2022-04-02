package authentication

import (
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/dbtypes"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Login verifies user credentials and returns a token
func Login(c *gin.Context) {
	// get email from parameter
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No email provided"})
		return
	}

	// get password from parameter
	password := c.PostForm("password")
	if password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No password provided"})
		return
	}

	// check login credentials
	if !userLogin(email, password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// create token
	//uuid, err := dbInterface.FindUUIDByEmail(email)
	uuid := "00000000-0000-0000-0000-000000000000"
	/*if err != nil {
		log.Println("Failed to find user uuid by email")
		c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid credentials"))
	}*/

	token, err := crypto.CreateToken(uuid)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// set token in cookie
	c.SetCookie("token", token, 3600, "", "", false, config.HTTPOnly)

	// send token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout deletes the token from the cookie
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, config.HTTPOnly)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}

// RenewToken renews the token
func RenewToken(c *gin.Context) {
	// get the token from the request
	token, err := c.Cookie("token")

	// if there is no token, return unauthorized
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	// parse the token
	parsedToken, err := jwt.Parse(token, crypto.VerifyToken)

	// if there is an error in parsing, return unauthorized
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// if the parsed token is not valid, return unauthorized
	if !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// get the identity from the token
	identity := parsedToken.Claims.(jwt.MapClaims)["identity"]

	// debug
	if config.AuthenticationDebug {
		log.Println("Renewing token for user ", identity)
	}

	// create a new token
	token, err = crypto.CreateToken(identity.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// send token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Adds a user to the database. only admins can do this.
func AddUser(c *gin.Context) {
	// verify token
	/*if !security.VerifyRegistrationToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		log.Println("INVALID REGISTRATION TOKEN")

		// http 401
		c.AbortWithStatus(http.StatusUnauthorized)
	}*/

	// TODO: verify user sending request is admin

	// get email
	email, success := c.GetPostForm("email")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No email provided"})
		return
	}

	// get password
	password, success := c.GetPostForm("password")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No password provided"})
		return
	}

	// TODO: fix role to have sanity checks. should we make this a read-only by default? should the registration token be attached to particular permissions or should we just let the user set it?
	// add role
	role, success := c.GetPostForm("role")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "No role provided"})
		return
	}

	// hash password
	hashedPassword, err := security.HashPass(password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	// TODO: figure out what cases a user will be added with a role
	var cases []string

	// add user to database
	_, err = dbInterface.MakeUser(
		dbtypes.NewUser{
			Email:    email,
			Password: hashedPassword,
			Role:     role,
			Cases:    cases,
		},
	)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to add user"})
		return
	}

	// return success
	c.Status(http.StatusOK)
}
