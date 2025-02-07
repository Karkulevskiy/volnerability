package containermgr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"volnerability-game/internal/cfg"
	"volnerability-game/internal/common"
	"volnerability-game/internal/domain"
	"volnerability-game/internal/lib/logger/utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

const (
	maxContainers    = 10
	maxTasks         = 100
	maxExecutionTime = 3 // TODO какое максимальное время на выполнение кода? Возможно оно должно меняться в зависимости от задания?
)

type Orchestrator struct {
	TempDir      string
	TargetDir    string
	ImageName    string
	Queue        chan domain.Task
	containers   []string
	available    chan string
	l            *slog.Logger
	dockerClient *client.Client
}

// TODO удаление временных файлов после остановки приложения
func New(l *slog.Logger, cfg cfg.OrchestratorConfig) (Orchestrator, error) {
	tempDir := os.TempDir() + "/codes"
	if err := os.Mkdir(tempDir, os.FileMode(os.O_CREATE|os.O_RDWR|os.O_APPEND)); err != nil {
		if !errors.Is(err, os.ErrExist) {
			l.Error("failed to create temp folder", utils.Err(err))
			return Orchestrator{}, err
		}
		l.Info("temp folder for codes already created")
	}

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation()) // map client api version
	if err != nil {
		l.Error("failed to create docker client", utils.Err(err))
		return Orchestrator{}, err
	}

	return Orchestrator{
		TempDir:      tempDir,
		TargetDir:    cfg.TargetDir,
		ImageName:    cfg.ImageName,
		Queue:        make(chan domain.Task, maxTasks),
		available:    make(chan string, maxContainers),
		containers:   make([]string, 0, maxContainers),
		dockerClient: cli,
		l:            l,
	}, nil
}

func (o *Orchestrator) Stop() error {
	for _, id := range o.containers {
		if err := o.dockerClient.ContainerPause(context.Background(), id); err != nil {
			o.l.Error(fmt.Sprintf("failed stop container: %s", id), utils.Err(err))
			return err
		}
	}

	o.l.Info(fmt.Sprintf("containers: [%s] were stopped", strings.Join(o.containers, ", ")))

	close(o.available)
	close(o.Queue)

	return nil
}

func (o *Orchestrator) RunContainers() error {
	//  Параметры запуска: docker run --name code-runner -v /home/spacikl/codes/:/home/ code-runner
	ctx := context.Background()

	if err := o.buildImage(ctx); err != nil {
		return err
	}

	containers, err := o.dockerClient.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		o.l.Error("failed to retrive existing containers", utils.Err(err))
		return err
	}

	if err := o.createContainers(ctx, containers); err != nil {
		return err
	}

	for _, id := range o.containers {
		if err := o.dockerClient.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
			o.l.Info("container already started, unpause")
			if err := o.dockerClient.ContainerUnpause(ctx, id); err != nil {
				o.l.Error(fmt.Sprintf("failed start container: %s", id), utils.Err(err))
				return err
			}
		}
	}

	go o.taskProcessor()

	return nil
}

func (o *Orchestrator) createContainers(ctx context.Context, containers []types.Container) error {
	containersSet := common.ToSetBy(containers, func(c types.Container) string {
		if len(c.Names) == 0 || len(c.Names[0]) <= 1 {
			return ""
		}
		return c.Names[0][1:]
	}) // docker create containers /code-runner-1, exclude "/"

	for i := range maxContainers {
		containerName := o.ImageName + "-" + strconv.Itoa(i)
		if _, ok := containersSet[containerName]; !ok {
			o.l.Info(fmt.Sprintf("create container: %s", containerName))
			if _, err := o.createContainer(ctx, containerName); err != nil {
				return err
			}
		}

		o.containers = append(o.containers, containerName)
		o.available <- containerName
	}

	return nil
}

func imageExist(ctx context.Context, targetImage string, c *client.Client) (bool, error) {
	images, err := c.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return false, err
	}

	if slices.ContainsFunc(images, func(i image.Summary) bool {
		return slices.ContainsFunc(i.RepoTags, func(name string) bool { return name == targetImage })
	}) {
		return true, nil
	}

	return false, nil
}

func createTar() (io.ReadCloser, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return archive.TarWithOptions(wd, &archive.TarOptions{IncludeFiles: []string{"Dockerfile"}})
}

func (o *Orchestrator) buildImage(ctx context.Context) error {
	ok, err := imageExist(ctx, o.ImageName, o.dockerClient)
	if err != nil {
		return err
	}

	if ok {
		o.l.Info("docker image already build")
		return nil
	}

	tar, err := createTar()
	if err != nil {
		o.l.Error("failed to crate tar file", utils.Err(err))
		return err
	}

	resp, err := o.dockerClient.ImageBuild(ctx, tar, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{o.ImageName},
	})

	if err != nil {
		o.l.Error("failed to build image", utils.Err(err))
		return err
	}
	defer resp.Body.Close()

	if _, err = io.Copy(os.Stdout, resp.Body); err != nil {
		o.l.Error("failed to get build response", utils.Err(err))
		return err
	}

	return nil
}

func (o *Orchestrator) createContainer(ctx context.Context, containerName string) (container.CreateResponse, error) {
	resp, err := o.dockerClient.ContainerCreate(ctx, &container.Config{
		Image: o.ImageName,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: o.TempDir,
				Target: o.TargetDir,
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
	file, err := os.CreateTemp(o.TempDir, "code-*."+t.Lang)
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
