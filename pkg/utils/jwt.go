package utils

import (
	"errors"
	"fmt"
	"gateway/config"
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

type jwtClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"sub"` // for map with claims inside
	Email  string
}

func ValidateToken(paramToken string, config *config.Config) (*jwtClaims, error) {
	if len(config.JwtConfig.SecretKey) > 0 {
		clams, err := validateBySecretKey(config.JwtConfig.SecretKey, paramToken)
		if err == nil {
			return clams, nil
		}
	}

	for _, path := range config.JwtConfig.PublicKey {
		clams, err := validateByKeyPem(path, paramToken)
		if err == nil {
			return clams, nil
		}
	}

	return nil, errors.New("Failed to parse token type")
}

func validateByKeyPem(pathKey string, paramToken string) (*jwtClaims, error) {
	var token *jwt.Token
	var err error
	dir, _ := os.Getwd()
	keyContent, err := os.ReadFile(fmt.Sprintf("%s%s%s", dir, "/", pathKey))
	if err != nil {
		return nil, fmt.Errorf("Error while read public key: %w", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyContent)
	token, err = jwt.ParseWithClaims(
		paramToken,
		&jwtClaims{},
		func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
			}
			return key, err
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("JWT is expired")
	}

	return claims, nil
}

func validateBySecretKey(secretKey string, paramToken string) (*jwtClaims, error) {
	var token *jwt.Token
	var err error
	key := []byte(secretKey)
	token, err = jwt.ParseWithClaims(
		paramToken,
		&jwtClaims{},
		func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
			}
			return key, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*jwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("Token is expired")
	}

	return claims, nil
}
