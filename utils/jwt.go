// utils/jwt.go
package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"log"
	"time"
)

// Secret key to sign JWTs (should be moved to environment variables in production)
var SecretKey = []byte("your-secret-key")

// GenerateJWT generates a JWT token for the user
func GenerateJWT(userID uint, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),  // Expiry time of 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return "", err
	}

	return signedToken, nil
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	// Parsing the token and validating its claims
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Invalid signing method")
			return nil, errors.New("Invalid signing method")  // err should be a custom error or use `errors.New("invalid signing method")`
		}
		return SecretKey, nil
	})

	if err != nil {
		log.Printf("Error validating token: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
