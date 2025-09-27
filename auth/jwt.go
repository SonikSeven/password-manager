package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(jwtSecret []byte, userID int64) (string, string, error) {
	// Access token
	accessTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	signedAccessToken, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedRefreshToken, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return signedAccessToken, signedRefreshToken, nil
}
