package coderunner

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
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
	file, err := os.CreateTemp("", "code-*.cpp")
	if err != nil {
		return "", err
	}
	defer os.Remove(file.Name())

	// TODO - идея постоянной перезаписи файлов, чтобы не создавать каждый файл. Или лучше - прокинуть эти файлы сразу в докер
	// и скриптами туда прокидывать код. Надо подумать ...

	if _, err := file.WriteString(code); err != nil {
		return "", err
	}

	containerName := "code-runner" // TODO просто как пример
	cmd := exec.Command("docker", "run", "--rm", "-v",
		fmt.Sprintf("%s:/code.cpp", file.Name()), // TODO научиться собирать мета инфу о файле, чтобы понимать в каком окружении собирать
		containerName, "sh", "-c", "g++ /code.cpp -o /code && /code")

	// Запускаем Docker
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return "hello", nil
}
