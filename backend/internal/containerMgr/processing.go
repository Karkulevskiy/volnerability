package containermgr

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

func wdPathForCodes() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	ind := strings.LastIndex(wd, "/")
	if ind == -1 {
		return "", fmt.Errorf("failed cut dir suffix")
	}
	return wd[:ind] + "/" + "codes", nil
}

func createFileName() string {
	return fmt.Sprintf("code-%s.%s", uuid.NewString(), "py")
}

func parseExecResp(output []byte) (string, error) {
	str := string(output)
	if len(str) <= 8 {
		return "", fmt.Errorf("failed to parse execution response: %v", output)
	}
	return str[8 : len(str)-1], nil
}

func cmd(fileName string) []string {
	return []string{"python3", "/home/" + fileName}
}
