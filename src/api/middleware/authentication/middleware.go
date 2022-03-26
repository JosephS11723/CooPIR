package authentication

import (
	"net/http"

	"github.com/JosephS11723/CooPIR/src/api/lib/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// AuthenticationMiddleware is a middleware that checks if the user is authenticated
func AuthenticationMiddleware(c *gin.Context) {
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

	// set the identity to the context
	c.Set("identity", identity)

	// call the next function
	c.Next()
}
