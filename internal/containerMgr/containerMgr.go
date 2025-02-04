package containermgr

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
	"volnerability-game/internal/domain"
	"volnerability-game/internal/lib/logger/utils"
)

const (
	maxContainers    = 10
	maxTasks         = 100
	maxExecutionTime = 3 // TODO какое максимальное время на выполнение кода? Возможно оно должно меняться в зависимости от задания?
)

type Orchestrator struct {
	Dir          string
	Queue        chan domain.Task
	results      map[string]string
	containers   []string
	containersMx sync.Mutex
	available    chan string
	l            *slog.Logger
}

// TODO остановка контейнеров
// TODO запускать остановленные контейнеры, если они уже созданы

func New(l *slog.Logger) Orchestrator {
	workingDir := os.TempDir() + "/codes"
	if err := os.Mkdir(workingDir, os.FileMode(os.O_CREATE|os.O_RDWR|os.O_APPEND)); err != nil {
		log.Fatal(err)
	}

	return Orchestrator{
		Dir:          workingDir,
		Queue:        make(chan domain.Task, maxTasks),
		available:    make(chan string, maxContainers),
		containers:   make([]string, 0, maxContainers),
		results:      map[string]string{},
		containersMx: sync.Mutex{},
		l:            l,
	}
}

func (o *Orchestrator) Stop() error {
	for _, containerName := range o.containers {
		if err := exec.Command("docker", "stop", "--name", containerName).Run(); err != nil {
			o.l.Error(fmt.Sprintf("failed stop container: %s", containerName), utils.Err(err))
			return err
		}
	}
	o.l.Info(fmt.Sprintf("containers: [%s] were stopped", strings.Join(o.containers, ", ")))
	return nil
}

func (o *Orchestrator) RunContainers() error {
	//  Параметры запуска: docker run --name code-runner -v /home/spacikl/codes/:/home/ code-runner
	ctx := context.Background()
	for i := range maxContainers {
		// TODO запускать по имени + i
		containerName := "code-runner-" + strconv.Itoa(i)
		if err := exec.Command("docker", "run", "--name", "code-runner", "-v", o.Dir, ":/home", "code-runner").Run(); err != nil {
			// TODO остановить запущенные контейнеры, если будет ошибка
			o.l.Error(fmt.Sprintf("failed start container: %s", containerName), utils.Err(err))
			return err
		}
		o.l.Info(fmt.Sprintf("start container: %s", containerName))
		o.containers = append(o.containers, containerName)
		o.available <- containerName
	}

	go o.taskProcessor(ctx)

	return nil
}

func (o *Orchestrator) taskProcessor(ctx context.Context) {
	for task := range o.Queue {
		containerId := <-o.available
		go o.executeTask(ctx, containerId, task)
	}
}

func (o *Orchestrator) executeTask(ctx context.Context, containerId string, t domain.Task) {
	o.l.Info(fmt.Sprintf("start task: [%s] on container: [%s]", t.ReqId, containerId))

	resp, err := o.runCode(containerId, t)
	if err != nil {
		o.l.Error("failed run code", utils.Err(err))
	}

	_ = resp
	// TODO закинуть респонс в канал ответа
	o.available <- containerId
}

func (o *Orchestrator) runCode(containerId string, t domain.Task) (domain.ExecuteResponse, error) {
	empty := domain.ExecuteResponse{}

	// TODO Хочется придумать пул свободных файлов, чтобы просто их перезаписывать
	file, err := os.CreateTemp(o.Dir, "code-*."+t.Lang)
	if err != nil {
		return empty, err
	}
	defer os.Remove(file.Name())

	if _, err := file.WriteString(t.Code); err != nil {
		return empty, err
	}

	if err := file.Close(); err != nil {
		return empty, err
	}

	cmd := o.cmd(file.Name(), t.Lang, containerId)

	slog.Info(cmd.String())

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = cmd.Run(); err != nil {
		return empty, err
	}

	execResp := domain.ExecuteResponse{
		Resp: stdout.String(),
	}

	return execResp, nil
}

func (o *Orchestrator) cmd(fileName, lang, containerId string) *exec.Cmd {
	fileName, _ = strings.CutPrefix(fileName, o.Dir)
	pathToFile := "/home/" + fileName

	slog.Info(fileName)
	runner := ""
	switch lang {
	case "c":
		runner = "" // TODO сюда нужно вставить команду на запуск си кода, если он будет. Пока на будущее
	case "py":
		runner = "python3"
	}

	return exec.Command("docker", "exec", containerId, runner, pathToFile)
}
