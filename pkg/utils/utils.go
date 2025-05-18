package utils

import (
	"errors"
	"fmt"
	"time"
	"twitter-api/pkg/errs"
	"twitter-api/pkg/types"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash: %w", err)
	}

	return string(hash), nil
}

func ComparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false
		}
		return false
	}
	return true
}

func GetAccessExpireTime() time.Time {
	const accessExpireTime = 24 * time.Hour // 1 day

	return time.Now().Add(accessExpireTime)
}

func GetRefreshExpireTime() time.Time {
	const refreshExpireTime = 7 * 24 * time.Hour // 7 days

	return time.Now().Add(refreshExpireTime)
}

func CreateAccessToken(userID int, role types.UserRole) (string, error) {
	accessToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"role":    role,
			"exp":     GetAccessExpireTime().Unix(),
		},
	)

	accessTokenStr, err := accessToken.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return accessTokenStr, nil
}

func CreateRefreshToken(userID int, role types.UserRole) (string, error) {
	refreshToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"user_id": userID,
			"role":    role,
			"exp":     GetRefreshExpireTime().Unix(),
		},
	)

	refreshTokenStr, err := refreshToken.SignedString(secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return refreshTokenStr, nil
}

func GenerateJwtTokens(userID int, role types.UserRole) (accessToken, refreshToken string, err error) {
	accessToken, err = CreateAccessToken(userID, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err = CreateRefreshToken(userID, role)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func VerifyRefreshToken(tokenStr string, userID int, userRole types.UserRole) error {
	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return errs.ErrTokenExpired
		}
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errs.ErrFailedToParseClaims
	}

	id, ok := claims["user_id"].(float64)
	if !ok || int(id) != userID {
		return errs.ErrInvalidUserID
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return errs.ErrFailedToParseUserRole
	}

	if types.UserRole(roleStr) != userRole {
		return errs.ErrInvalidUserRole
	}

	return nil
}

func VerifyAccessToken(tokenStr string) (int, types.UserRole, error) {
	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		if _, ok := tkn.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tkn.Header["alg"])
		}
		return secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, "", errs.ErrTokenExpired
		}
		return 0, "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, "", errs.ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", errs.ErrFailedToParseClaims
	}

	id, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", errs.ErrFailedToParseUserID
	}

	roleStr, ok := claims["role"].(string)
	if !ok {
		return 0, "", errs.ErrFailedToParseUserRole
	}

	return int(id), types.UserRole(roleStr), nil
}
