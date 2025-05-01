package containermgr

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/pkg/archive"
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

func createTar() (io.ReadCloser, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return archive.TarWithOptions(wd, &archive.TarOptions{IncludeFiles: []string{"Dockerfile"}})
}

func createFileName(lang string) string {
	return fmt.Sprintf("code-%s.%s", uuid.NewString(), lang)
}

func parseExecResp(output []byte) (string, error) {
	str := string(output)
	if len(str) <= 8 {
		return "", fmt.Errorf("failed to parse execution response: %v", output)
	}
	return str[8 : len(str)-1], nil
}

func cmd(fileName, lang string) []string {
	runner := ""
	switch lang {
	case "c":
		runner = "" // TODO сюда нужно вставить команду на запуск си кода, если он будет. Пока на будущее
	case "py":
		runner = "python3"
	}
	return []string{runner, "/home/" + fileName}
}
