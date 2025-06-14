package containermgr

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// I know, tests are useless. Btw. They are required for project :)

// func TestCreateFileName(t *testing.T) {
// 	tests := []struct {
// 		expectedName string
// 		lang         string
// 	}{
// 		{expectedName: "code-", lang: "py"},
// 		{expectedName: "code-", lang: "c"},
// 	}
// 	for _, test := range tests {
// 		fileName := createFileName(test.lang)
// 		assert.Contains(t, fileName, test.expectedName)
// 		assert.Contains(t, fileName, test.lang)
// 	}
// }

func TestWdPathForCodes(t *testing.T) {
	tests := []struct {
		expectedWd string
	}{
		{expectedWd: "/codes"},
	}
	for _, test := range tests {
		currWd, err := wdPathForCodes()
		assert.NoError(t, err)
		lastFolderInd := strings.LastIndex(currWd, "/")
		assert.Equal(t, test.expectedWd, currWd[lastFolderInd:])
	}
}

func TestParseExecResp(t *testing.T) {
	tests := []struct {
		output      []byte
		expResponse string
		err         error
	}{
		{
			output:      []byte("12345678good payload "),
			expResponse: "good payload",
			err:         nil,
		},
		{
			output:      []byte("1234567"),
			expResponse: "",
			err:         fmt.Errorf("failed to parse execution response: %v", []byte("1234567")),
		},
	}
	for _, test := range tests {
		resp, err := parseExecResp(test.output)
		if test.err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
		assert.Equal(t, test.expResponse, resp)
	}
}

// func TestCmd(t *testing.T) {
// 	tests := []struct {
// 		fileName    string
// 		lang        string
// 		expectedCmd []string
// 	}{
// 		{fileName: "myFile", lang: "py", expectedCmd: []string{"python3", "/home/myFile"}},
// 	}
// 	for _, test := range tests {
// 		resp := cmd(test.fileName, test.lang)
// 		assert.Equal(t, test.expectedCmd, resp)
// 	}
// }
