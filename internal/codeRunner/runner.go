package coderunner

import (
	"bytes"
	"context"
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
	dir string
}

func New(dir string) *CodeRunner {
	return &CodeRunner{
		dir: dir,
	}
}

func (r *CodeRunner) Run(code string, lang string) (string, error) {
	// TODO нужно быстро уметь валидировать, что код вообще билдиться
	// Механизм кеширования, используя очередь
	file, err := os.CreateTemp(r.dir, "code-*."+lang)
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

	cmd := r.cmd(file.Name(), lang)

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

func (r *CodeRunner) cmd(fileName, lang string) *exec.Cmd {
	fileName, _ = strings.CutPrefix(fileName, r.dir)
	pathToFile := "/home/" + fileName

	slog.Info(fileName)
	runner := ""
	switch lang {
	case "c":
		runner = "" // TODO
	case "py":
		runner = "python3"
	}

	return exec.Command("docker", "exec", containerName(), runner, pathToFile)
}

// TODO спрашивать из свободного пула контейнеров у оркестратора
func containerName() string {
	// TODO для оркестратора
	return "code-runner"
}
