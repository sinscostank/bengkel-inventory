// utils/auth.go
package utils

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

// GenerateHash generates a hash from a password
func GenerateHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error generating hash: %v", err)
		return "", err
	}
	return string(hash), nil
}

// CheckHash checks if the given password matches the hashed password
func CheckHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		log.Println("Password mismatch:", err)
		return false
	}
	return true
}
