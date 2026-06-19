package jwt

import (
	"coi/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

const (
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 7 * 24 * time.Hour
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type MyClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func getAccessSecretKey() []byte {
	return []byte(config.GetEnv("JWT_ACCESS_SECRET_KEY", config.GetEnv("JWT_SECRET_KEY", "default_secret_key")))
}

func getRefreshSecretKey() []byte {
	return []byte(config.GetEnv("JWT_REFRESH_SECRET_KEY", config.GetEnv("JWT_SECRET_KEY", "default_secret_key")))
}

func GenerateTokenPair(userID int, username string) (*TokenPair, error) {
	accessToken, err := generateToken(userID, username, TokenTypeAccess, defaultAccessTokenTTL, getAccessSecretKey())
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(userID, username, TokenTypeRefresh, defaultRefreshTokenTTL, getRefreshSecretKey())
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func GenerateToken(userID int, username string) (string, error) {
	tokenPair, err := GenerateTokenPair(userID, username)
	if err != nil {
		return "", err
	}
	return tokenPair.AccessToken, nil
}

func RefreshAccessToken(refreshTokenString string) (string, error) {
	claims, err := ValidateRefreshToken(refreshTokenString)
	if err != nil {
		return "", err
	}

	return generateToken(claims.UserID, claims.Username, TokenTypeAccess, defaultAccessTokenTTL, getAccessSecretKey())
}

func ValidateAccessToken(tokenString string) (*MyClaims, error) {
	return validateToken(tokenString, getAccessSecretKey(), TokenTypeAccess)
}

func ValidateRefreshToken(tokenString string) (*MyClaims, error) {
	return validateToken(tokenString, getRefreshSecretKey(), TokenTypeRefresh)
}

func ValidateToken(tokenString string) (*MyClaims, error) {
	return ValidateAccessToken(tokenString)
}

func generateToken(userID int, username, tokenType string, ttl time.Duration, secret []byte) (string, error) {
	claims := MyClaims{
		UserID:   userID,
		Username: username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "coi",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateToken(tokenString string, secret []byte, expectedTokenType string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing algorithm: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		if claims.TokenType != expectedTokenType {
			return nil, fmt.Errorf("invalid token type: expected %s", expectedTokenType)
		}
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}