package crypto

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// createToken creates a token for the user and sends it
func CreateToken(identity string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims = jwt.MapClaims{
		"identity": identity,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// sign the token
	tokenString, err := token.SignedString(PrivateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// verifyToken verifies a token
func VerifyToken(token *jwt.Token) (interface{}, error) {
	// check if the token is valid
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, jwt.ErrSignatureInvalid
	}

	return PublicKey, nil
}
