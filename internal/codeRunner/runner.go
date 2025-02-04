package coderunner

import (
	"bytes"
	"context"
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"time"
	"volnerability-game/internal/domain"
)

type Runner interface {
	Run(code string) (string, error)
}

type CodeRunner struct {
	dir   string
	queue chan domain.Task
}

func New(dir string, queue chan domain.Task) *CodeRunner {
	return &CodeRunner{
		dir:   dir,
		queue: queue,
	}
}

func (r *CodeRunner) Run(code string, lang string) (string, error) {
	// TODO нужно быстро уметь валидировать, что код вообще билдиться
	// Механизм кеширования, используя очередь

	task := domain.Task{
		Code: code,
		Lang: lang,
	}

	r.queue <- task

}
