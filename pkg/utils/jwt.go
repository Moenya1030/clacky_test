package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"task-manager/config"
)

// CustomClaims defines the claims structure for JWT tokens
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken creates a JWT token for the given user ID
func GenerateToken(userID uint) (string, error) {
	// Get JWT configuration
	jwtConfig := config.GetConfig().JWT

	// Create token claims
	now := time.Now()
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(jwtConfig.ExpiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID if valid
func ValidateToken(tokenString string) (uint, error) {
	if tokenString == "" {
		return 0, errors.New("empty token")
	}

	// Get JWT configuration
	jwtConfig := config.GetConfig().JWT

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtConfig.Secret), nil
		},
	)

	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return 0, errors.New("invalid token claims")
	}

	return claims.UserID, nil
}

// GetUserIDFromToken extracts the user ID from a valid JWT token
func GetUserIDFromToken(tokenString string) (uint, error) {
	// Get JWT configuration
	jwtConfig := config.GetConfig().JWT

	// Parse the token without validation (just to extract claims)
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtConfig.Secret), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, errors.New("token expired")
		}
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract user ID from claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}