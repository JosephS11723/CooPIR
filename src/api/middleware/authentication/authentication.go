package authentication

import (
	"errors"
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/config"
	"github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthenticationMiddleware is a middleware that checks if the user is authenticated
func AuthenticationMiddleware(c *gin.Context) {
	// debug config boolean (no login if in auth debug mode)
	if !config.RequireLogin {
		c.Set("identity", "00000000-0000-0000-0000-000000000000")
		c.Next()
		return
	}

	// get the token from the request
	token, err := c.Cookie("token")

	// if there is no token, return unauthorized
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// parse the token
	parsedToken, err := jwt.Parse(token, crypto.VerifyToken)

	// if there is an error in parsing, return unauthorized
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// if the parsed token is not valid, return unauthorized
	if !parsedToken.Valid {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// get the identity from the token
	identity := parsedToken.Claims.(jwt.MapClaims)["identity"]

	// if the identity is not a string, return unauthorized
	if identity.(string) == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid identity"))
		return
	}

	// if the identity is an empty string, return unauthorized
	if identity.(string) == "" {
		c.AbortWithError(http.StatusUnauthorized, errors.New("invalid identity"))
		return
	}

	// set the identity to the context
	c.Set("identity", identity)

	// call the next function
	c.Next()
}
