package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"twitter-api/pkg/types"
)

func TestComparePassword_CorrectPassword(t *testing.T) {
	password := "mySecret123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if !ComparePassword(string(hashed), password) {
		t.Errorf("expected password to match, but it did not")
	}
}

func TestComparePassword_WrongPassword(t *testing.T) {
	password := "mySecret123"
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if ComparePassword(string(hashed), "wrongPassword") {
		t.Errorf("expected password to not match, but it did")
	}
}

func TestComparePassword_InvalidHash(t *testing.T) {
	invalidHash := "not-a-valid-bcrypt-hash"
	password := "anyPassword"

	if ComparePassword(invalidHash, password) {
		t.Errorf("expected function to return false for invalid hash, but it returned true")
	}
}
func TestGetAccessExpireTime(t *testing.T) {
	expectedMin := time.Now().Add(24 * time.Hour)
	actual := GetAccessExpireTime()

	if actual.Before(expectedMin.Add(-1*time.Second)) || actual.After(expectedMin.Add(1*time.Second)) {
		t.Errorf("Access token expiry not within 1 second margin: got %v, expected around %v", actual, expectedMin)
	}
}

func TestGetRefreshExpireTime(t *testing.T) {
	expectedMin := time.Now().Add(7 * 24 * time.Hour)
	actual := GetRefreshExpireTime()

	if actual.Before(expectedMin.Add(-1*time.Second)) || actual.After(expectedMin.Add(1*time.Second)) {
		t.Errorf("Refresh token expiry not within 1 second margin: got %v, expected around %v", actual, expectedMin)
	}
}

var testSecretKey = []byte("test-secret")

func init() {
	secretKey = testSecretKey
}

func TestCreateAccessToken(t *testing.T) {
	tokenStr, err := CreateAccessToken(123, types.UserRole("admin"))
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return testSecretKey, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, float64(123), claims["user_id"])
	require.Equal(t, "admin", claims["role"])

	exp, ok := claims["exp"].(float64)
	require.True(t, ok)
	require.Greater(t, int64(exp), time.Now().Unix())
}

func TestCreateRefreshToken(t *testing.T) {
	tokenStr, err := CreateRefreshToken(123, types.UserRole("admin"))
	require.NoError(t, err)
	require.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return testSecretKey, nil
	})
	require.NoError(t, err)
	require.True(t, token.Valid)

	claims, ok := token.Claims.(jwt.MapClaims)
	require.True(t, ok)
	require.Equal(t, float64(123), claims["user_id"])
	require.Equal(t, "admin", claims["role"])

	exp, ok := claims["exp"].(float64)
	require.True(t, ok)
	require.Greater(t, int64(exp), time.Now().Unix())
}
