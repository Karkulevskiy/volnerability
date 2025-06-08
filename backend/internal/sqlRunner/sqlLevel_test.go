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
			`SELECT   *   FROM   users   WHERE   username   =   '   OR   'a'   =   'a'   AND   password   =   ' OR 'a'   =   'a'`,
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

func TestSecondLevelRegexp(t *testing.T) {
	tests := []struct {
		input       string
		isInjection bool
	}{
		// SELECT name, email FROM clients WHERE name LIKE '%$search%'
		// $search = "' UNION SELECT username, password FROM users--"
		// Result: SELECT name, email FROM clients WHERE name LIKE '%' UNION SELECT username, password FROM users--%'
		{"SELECT name, email FROM clients WHERE name LIKE '%' UNION SELECT username, password FROM users--%'", true},
		{"SELECT   name  ,   email   FROM   clients   WHERE   name   LIKE   '%'   UNION   SELECT   username  ,   password   FROM   users--%'", true},
		{"SELECT name, email FROM clients WHERE name LIKE '%' UNION select username, password FROM users--%'", true},
		{"SELECT name, email FROM clients WHERE name LIKE '%' UNION sxlect username, password FROM users--%'", false},
		{"SELECT name, email FOM clients WHERE name LIKE '%' UNION SELECT username, password FROM users--%'", false},
	}
	for _, test := range tests {
		fmt.Println("Testing input:", test.input)
		isMatch := isSecondSqlInjection(test.input)
		require.Equal(t, test.isInjection, isMatch)
	}
}

func TestThirdLevelRegexp(t *testing.T) {
	tests := []struct {
		input       string
		isInjection bool
	}{
		// db.Exec("INSERT INTO feedback (text) VALUES ('" + userInput + "')")
		// userInput = '); DROP TABLE users;--
		// Result: INSERT INTO feedback (text) VALUES (''); DROP TABLE users;--')
		{"'; DROP TABLE users;--", true},
		{"'; DROP TABLE users; --", true},
		{"';DROP TABLE users;--", true},
		{"'; DROP TABLE users;-- extra text", true},
		{"'; DROP TABLE users", false},
		{"'; DROP TABLE other_table;--", false},
		{"'; DROP TABLE users; -- extra", true},
		{"'; INSERT INTO feedback (text) VALUES ('test');--", false},
		{"'; DROP TABLE users;--", true},
		{"'; DROP TABLE users;-- more text", true},
	}
	for _, test := range tests {
		fmt.Println("Testing input:", test.input)
		isMatch := isThirdSqlInjection(test.input)
		require.Equal(t, test.isInjection, isMatch)
	}
}
