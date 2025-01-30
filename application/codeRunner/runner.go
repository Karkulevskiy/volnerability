package coderunner

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Runner interface {
	Run(code string) (string, error)
}

type CodeRunner struct {
}

func New() *CodeRunner {
	return &CodeRunner{}
}

func (r *CodeRunner) Run(code string) (string, error) {
	// TODO нужно быстро уметь валидировать, что код вообще билдиться
	// TODO сделать унифицированную папку
	file, err := os.CreateTemp("/home/spacikl/codes", "code-*.py")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(code); err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	fileName, ok := strings.CutPrefix(file.Name(), "/home/spacikl/codes/")
	if !ok {
		return "", fmt.Errorf("write here error")
	}
	fileName = "/home/" + fileName
	slog.Info(fileName)
	containerName := "code-runner"
	//  Параметры запуска: docker run --name code-runner -v /home/spacikl/codes/:/home/ code-runner
	cmd := exec.Command("docker", "exec", containerName, "python3", fileName)

	slog.Info(cmd.String())

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return stdout.String(), nil
}
