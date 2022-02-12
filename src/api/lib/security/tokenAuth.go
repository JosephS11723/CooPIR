package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// VerifyToken verifies the token in the request.
func VerifyToken(c *gin.Context) bool {
	// get token from request header
	token, err := c.Request.Cookie("token")
	if err != nil {
		return false
	}

	// verify token
	_, err = jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	// return result
	if err == nil {
		return true
	} else {
		return false
	}
}

// CreateToken creates a new JWT token.
func CreateToken(c *gin.Context) (string, error) {
}


// RevokeToken revokes the current JWT token.
func RevokeToken(c *gin.Context) {


}