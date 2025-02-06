package containermgr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"volnerability-game/internal/common"
	"volnerability-game/internal/domain"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
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
	available    chan string
	l            *slog.Logger
	dockerClient *client.Client
}

// TODO удаление временных файлов после остановки приложения
func New(l *slog.Logger) (Orchestrator, error) {
	workingDir := os.TempDir() + "/codes"
	if err := os.Mkdir(workingDir, os.FileMode(os.O_CREATE|os.O_RDWR|os.O_APPEND)); err != nil {
		l.Error("failed to create temp folder", utils.Err(err))
		return Orchestrator{}, err
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		l.Error("failed to create docker client", utils.Err(err))
		return Orchestrator{}, err
	}

	return Orchestrator{
		Dir:          workingDir,
		Queue:        make(chan domain.Task, maxTasks),
		available:    make(chan string, maxContainers),
		containers:   make([]string, 0, maxContainers),
		results:      map[string]string{},
		dockerClient: cli,
		l:            l,
	}, nil
}

func (o *Orchestrator) Stop() error {
	for _, id := range o.containers {
		if err := o.dockerClient.ContainerStop(context.Background(), id, container.StopOptions{}); err != nil { // deafult timeout 10s
			o.l.Error(fmt.Sprintf("failed stop container: %s", id), utils.Err(err))
			return err
		}
	}

	o.l.Info(fmt.Sprintf("containers: [%s] were stopped", strings.Join(o.containers, ", ")))

	close(o.available)
	close(o.Queue)

	return nil
}

// TODO получение имен контейнеров из енва
func (o *Orchestrator) RunContainers() error {
	//  Параметры запуска: docker run --name code-runner -v /home/spacikl/codes/:/home/ code-runner
	ctx := context.Background()
	containers, err := o.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		o.l.Error("failed to retrive existing containers", utils.Err(err))
		return err
	}

	if err := o.createContainers(ctx, containers); err != nil {
		return err
	}

	for _, id := range o.containers {
		if err := o.dockerClient.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
			o.l.Error(fmt.Sprintf("failed start container: %s", id), utils.Err(err))
			return err
		}
	}

	go o.taskProcessor()

	return nil
}

func (o *Orchestrator) createContainers(ctx context.Context, containers []types.Container) error {
	containersSet := common.ToSetBy(containers, func(c types.Container) string { return c.ID })

	for i := range maxContainers {
		containerName := "code-runner-" + strconv.Itoa(i)

		if _, ok := containersSet[containerName]; ok {
			continue
		}

		resp, err := o.createContainer(ctx, containerName)
		if err != nil {
			return err
		}
		o.l.Info(fmt.Sprintf("start container: %s", resp.ID))

		o.containers = append(o.containers, resp.ID)
		o.available <- containerName
	}

	return nil
}

func (o *Orchestrator) createContainer(ctx context.Context, containerName string) (container.CreateResponse, error) {
	resp, err := o.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: "python:3.9-slim", // TODO найти легковесные образы
		Cmd:   []string{"sleep", "infinity"},
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume,
				Source: o.Dir,
				Target: ":/home", // TODO некрасиво, нужно это все брать из конфига (файла)
			},
		},
	}, nil, nil, containerName)

	if err != nil {
		o.l.Error(fmt.Sprintf("failed create container: %s", containerName), utils.Err(err))
		return container.CreateResponse{}, err
	}

	return resp, nil
}

func (o *Orchestrator) taskProcessor() {
	for task := range o.Queue {
		containerId := <-o.available
		go o.executeTask(containerId, task)
	}
}

func (o *Orchestrator) executeTask(containerId string, t domain.Task) {
	o.l.Info(fmt.Sprintf("start task: [%s] on container: [%s]", t.ReqId, containerId))

	resp, err := o.runCode(containerId, t)
	if err != nil {
		o.l.Error("failed run code", utils.Err(err))
	}

	t.Resp <- resp
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

	ctx := context.Background()

	execResp, err := o.dockerClient.ContainerExecCreate(ctx, containerId, container.ExecOptions{
		Cmd:          cmd(file.Name(), t.Lang),
		AttachStdout: true,
		AttachStderr: true,
		Tty:          false,
	})

	if err != nil {
		return empty, err
	}

	attachResp, err := o.dockerClient.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{}) // TODO нужно для каждого контейнера держать свой процесс?? И просто через какой то метод его получать?
	if err != nil {
		return empty, err
	}
	defer attachResp.Close()

	output, err := io.ReadAll(attachResp.Reader)
	if err != nil {
		return empty, err
	}

	return domain.ExecuteResponse{Resp: string(output)}, nil
}

func cmd(fileName, lang string) []string {
	runner := ""
	switch lang {
	case "c":
		runner = "" // TODO сюда нужно вставить команду на запуск си кода, если он будет. Пока на будущее
	case "py":
		runner = "python3"
	}
	return []string{runner, fileName}
}
