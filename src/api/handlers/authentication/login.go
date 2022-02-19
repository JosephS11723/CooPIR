package authentication

import (
	"log"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/dbInterface"
	"github.com/JosephS11723/CooPIR/src/api/lib/security"
	"github.com/gin-gonic/gin"
)

// Login is a function to handle the login request by verifying the user's email and password and returning a JWT token.
func Login(c *gin.Context) {
	// get email
	email, success := c.GetPostForm("email")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no username provided"})
		return
	}

	// get password
	password, success := c.GetPostForm("password")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no password provided"})
		return
	}

	// get user from database
	passHash := dbInterface.RetrieveHashByEmail(email)

	// verify password
	if !security.CheckPass(password, passHash) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid password"})
		return
	}

	// create new token
	token, err := security.CreateToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to create token"})
		return
	}

	// return new token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// RenewToken renews the current JWT token.
func RenewToken(c *gin.Context) {
	// verify token
	if !security.VerifyToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// create new token
	token, err := security.CreateToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to create token"})
		return
	}

	// return new token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout revokes the current JWT token.
func Logout(c *gin.Context) {
	// revoke token
	err := security.RevokeToken(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		log.Panicln("INVALID TOKEN")
	}

	// return success
	c.Status(http.StatusOK)
}

// Adds a user to the database by consuming a registration token.
func AddUser(c *gin.Context) {
	// verify token
	if !security.VerifyRegistrationToken(c) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		log.Panicln("INVALID REGISTRATION TOKEN")
	}

	// get email
	email, success := c.GetPostForm("email")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no email provided"})
		return
	}

	// get password
	password, success := c.GetPostForm("password")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no password provided"})
		return
	}

	// get username
	username, success := c.GetPostForm("username")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no username provided"})
		return
	}

	// TODO: fix role to have sanity checks. should we make this a read-only by default? should the registration token be attached to particular permissions or should we just let the user set it?
	// add role
	role, success := c.GetPostForm("role")
	if !success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "no role provided"})
		return
	}

	// hash password
	hashedPassword, err := security.HashPass(password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to hash password"})
		return
	}

	// TODO: figure out what cases a user will be added with a role
	var cases []string

	// add user to database
	_, err = dbInterface.MakeUser(username, email, role, cases, hashedPassword)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "failed to add user"})
		return
	}

	// return success
	c.Status(http.StatusOK)
}