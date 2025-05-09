package common

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsEmailValid(t *testing.T) {
	tests := []struct {
		email string
		valid bool
	}{
		{"fds", false},
		{"lae@gmail.com", true},
		{"laegmail.com", false},
		{"lae@.gmail.com", false},
		{"1@l.cc", true},
		{"1@l.c", true},
	}
	for _, test := range tests {
		isValid := IsEmailValid(test.email)
		assert.Equal(t, test.valid, isValid)
	}
}

func TestFromErrors(t *testing.T) {
	tests := []struct {
		baseErr    error
		targetErrs []error
		found      bool
	}{
		{fmt.Errorf("unknown error"), []error{ErrBadSubmit, ErrEmptyInput}, false},
		{fmt.Errorf("known error: %w", ErrEmptyInput), []error{ErrBadSubmit, ErrEmptyInput}, true},
	}
	for _, test := range tests {
		found := FromErrors(test.baseErr, test.targetErrs...)
		assert.Equal(t, test.found, found)
	}
}

func TestParseToken(t *testing.T) {
	tests := []struct {
		tokenString  string
		appSecret    []byte
		isTokenValid bool
		exptectedErr error
	}{
		{
			tokenString:  signToken(t, []byte("appSecret")),
			appSecret:    []byte("appSecret"),
			isTokenValid: true,
			exptectedErr: nil,
		},
		{
			tokenString:  signToken(t, []byte("notAppSecret")),
			appSecret:    []byte("appSecret"),
			isTokenValid: false,
			exptectedErr: ErrInvalidToken,
		},
		{
			tokenString:  signTokenOtherMethod(t, []byte("appSecret")),
			appSecret:    []byte("appSecret"),
			isTokenValid: false,
			exptectedErr: ErrInvalidSigningMethod,
		},
	}

	for _, test := range tests {
		_, err := ParseToken(test.tokenString, test.appSecret)
		assert.Equal(t, test.exptectedErr, err)
	}
}

func signToken(t *testing.T, appSecret []byte) string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString(appSecret)
	require.NoError(t, err)
	return tokenString
}

func signTokenOtherMethod(t *testing.T, appSecret []byte) string {
	token := jwt.New(jwt.SigningMethodHS384)
	tokenString, err := token.SignedString(appSecret)
	require.NoError(t, err)
	return tokenString
}
