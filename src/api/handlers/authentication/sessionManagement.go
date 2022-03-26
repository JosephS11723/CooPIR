package authentication

import (
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// Login verifies user credentials and returns a token
func Login(c *gin.Context) {
	// get email from formdata
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No email provided"})
		return
	}

	// get password from formdata
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
	token, err := crypto.CreateToken(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// send token
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Logout deletes the token from the cookie
func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "", "", false, true)
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

	// create a new token
	token, err = crypto.CreateToken(identity.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	// send token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
