package sqlrunner

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFirstLevelRegexp(t *testing.T) {
	tests := []struct {
		input       string
		isInjection bool
	}{
		{
			`SELECT * FROM users WHERE username = ' OR 'a' = 'a' AND password = ' OR 'a' = 'a'`,
			true,
		},
		{
			`SELECT * FROM users WHERE username = ' OR '1' = '1' AND password = ' OR '1' = '1'`,
			true,
		},
		{
			`SELECT * FROM users WHERE username = ' OR '1' = '1' AND password = ' OR 'x1' = 'x1'`,
			true,
		},
		{
			`SELECT * FROM users WHERE username = ' OR '2' = '1' AND password = ' OR 'x1' = 'x1'`,
			false,
		},
		{
			`SELECT * FROM users WHERE username =  OR '1' = '1' AND password = ' OR 'x1' = 'x1'`,
			false,
		},
		{
			`SELECT * FROM users WHERE username =  OR '1' = '1' AND password = ' OR '1' = 'x1'`,
			false,
		},
		{
			`SELECT * FROM users WHERE username =  OR '1' = '1' AND password = ' OR '1 = 'x1'`,
			false,
		},
		{
			`SELECT * FROM users WHERE username =  OR '1' = 1' AND password = ' OR '1 = 'x1'`,
			false,
		},
	}

	for _, test := range tests {
		fmt.Println("Testing input:", test.input)
		isMatch := isFirstSqlInjection(test.input)
		require.Equal(t, test.isInjection, isMatch)
	}
}
