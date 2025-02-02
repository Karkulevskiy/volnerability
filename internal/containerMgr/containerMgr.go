package containermgr

import (
	"log"
	"os"
	"os/exec"
)

type Orchestrator struct {
	Dir string
}

func New() Orchestrator {
	workingDir := os.TempDir() + "/codes"
	if err := os.Mkdir(workingDir, os.FileMode(os.O_CREATE|os.O_RDWR|os.O_APPEND)); err != nil {
		log.Fatal(err)
	}
	return Orchestrator{Dir: workingDir}
}

func (o *Orchestrator) RunContainers() error {
	//  Параметры запуска: docker run --name code-runner -v /home/spacikl/codes/:/home/ code-runner
	cmd := exec.Command("docker", "run", "--name", "code-runner", "-v", o.Dir, ":/home", "code-runner")
	return cmd.Run()
}
