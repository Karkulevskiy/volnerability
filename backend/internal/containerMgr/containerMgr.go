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
	maxContainers = 10
	maxTasks      = 100
)

type Orchestrator struct {
	WorkingDir   string
	TargetDir    string
	ImageName    string
	Queue        chan domain.Task
	containers   []string
	available    chan string
	l            *slog.Logger
	dockerClient *client.Client
}

func New(l *slog.Logger, cfg cfg.OrchestratorConfig) (Orchestrator, error) {
	wdPath, err := wdPathForCodes()
	if err != nil {
		l.Error("failed get working dir")
		return Orchestrator{}, err
	}

	if err := os.Mkdir(wdPath, 0777); err != nil {
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
		WorkingDir:   wdPath,
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
		if err := o.dockerClient.ContainerStop(context.Background(), id, container.StopOptions{}); err != nil {
			o.l.Error(fmt.Sprintf("failed to stop container: %s", id), utils.Err(err))
			return err
		}
	}

	for _, id := range o.containers {
		if err := o.dockerClient.ContainerRemove(context.Background(), id, container.RemoveOptions{}); err != nil {
			o.l.Error(fmt.Sprintf("failed to delete container: %s", id), utils.Err(err))
			return err
		}
	}

	o.l.Info(fmt.Sprintf("containers: [%s] were deleted", strings.Join(o.containers, ", ")))

	close(o.available)
	close(o.Queue)

	if err := o.cleanTempFiles(); err != nil {
		return err
	}

	return nil
}

func (o *Orchestrator) cleanTempFiles() error {
	if err := os.RemoveAll(o.WorkingDir); err != nil {
		o.l.Error("failed to remove working directory with codes", utils.Err(err))
		return err
	}
	return nil
}

func (o *Orchestrator) RunContainers() error {
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
		o.l.Info(fmt.Sprintf("start container: %s", id))
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

func (o *Orchestrator) buildImage(ctx context.Context) error {
	ok, err := imageExist(ctx, o.ImageName, o.dockerClient)
	if err != nil {
		return err
	}

	if ok {
		o.l.Info("docker image already build")
		return nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	tar, err := archive.TarWithOptions(wd, &archive.TarOptions{IncludeFiles: []string{"Dockerfile"}})
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
				Source: o.WorkingDir,
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

	select {
	case t.Resp <- resp:
	default:
		o.l.Info(fmt.Sprintf("task: %s was closed, data lost", t.ReqId))
	}

	o.available <- containerId
}

func (o *Orchestrator) runCode(containerId string, t domain.Task) (domain.ExecuteResponse, error) {
	empty := domain.ExecuteResponse{}
	// TODO Хочется придумать пул свободных файлов, чтобы просто их перезаписывать
	fileName := createFileName()
	file, err := os.OpenFile(o.WorkingDir+"/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
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
		Cmd:          cmd(fileName),
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

	resp, err := parseExecResp(output)
	if err != nil {
		return empty, err
	}

	fmt.Printf("\n\nRESPONSE: %s\n\n", resp)
	return domain.ExecuteResponse{Resp: resp}, nil
}
