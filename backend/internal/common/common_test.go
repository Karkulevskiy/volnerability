package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
		fmt.Println(test.email, isValid, test.valid)
		assert.Equal(t, test.valid, isValid)
	}
}
