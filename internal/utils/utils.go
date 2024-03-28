package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
)

func GenerateBase64Hash() string {
	// Generate 16 random bytes
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Printf("Error generating random bytes: %v", err)
		return ""
	}

	// Hash the random bytes using SHA-256
	hash := sha256.Sum256(randomBytes)
	// Encode the hash to a base64 string
	hashBase64 := base64.URLEncoding.EncodeToString(hash[:])

	return hashBase64
}

func GetEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetName(firstName, lastName string) string {
	result := ""
	if firstName != "" {
		result = firstName
	}
	if lastName != "" {
		if result != "" {
			result += " "
		}
		result += lastName
	}
	return result
}

func GetLanguage(lang string) string {
	if lang == "ru" {
		return "ru"
	}
	return "en"
}
