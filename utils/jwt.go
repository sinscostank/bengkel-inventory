// utils/jwt.go
package utils

import (
	"os"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// Secret key to sign JWTs (should be moved to environment variables in production)
var SecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

// UserClaims is the custom claims structure for the JWT.
type UserClaims struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, email, role string) (string, error) {
	claims := UserClaims{
		ID:    userID,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "bengkel-inventory",
		},
	}

	// Create the token with claims and sign it with the secret key
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(SecretKey)
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*UserClaims, error) {
	// Parse the JWT string and validate it
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the signing method is correct
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract the claims from the token
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
