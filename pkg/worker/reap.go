package worker

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/convox/console/pkg/settings"
	"github.com/pkg/errors"
)

var (
	reapInterval = 1 * time.Minute
)

func init() {
	if settings.Development {
		reapInterval = 10 * time.Second
	}
}

func processRunning(pid string) bool {
	if settings.Development {
		return processRunningLocal(pid)
	} else {
		return processRunningRack(pid)
	}
}

func processRunningLocal(pid string) bool {
	ipid, err := strconv.Atoi(pid)
	if err != nil {
		return false
	}

	ps, err := os.FindProcess(ipid)
	if err != nil {
		return false
	}

	if err := ps.Signal(syscall.Signal(0)); err != nil {
		return false
	}

	return true
}

func processRunningRack(pid string) bool {
	if _, err := Rack.ProcessGet(os.Getenv("APP"), pid); err != nil {
		fmt.Printf("ns=worker at=processRunningRack error=%q\n", err)
		return false
	}

	return true
}

func reap() {
	for range time.Tick(reapInterval) {
		if err := reapJobs(); err != nil {
			fmt.Printf("err: %+v\n", err)
			continue
		}
	}
}

func reapJobs() error {
	js, err := Models.JobListByStatus("running")
	if err != nil {
		return errors.WithStack(err)
	}

	for _, j := range js {
		if j.Created.After(time.Now().Add(-1 * time.Minute)) {
			continue
		}

		if !processRunning(j.Pid) {
			fmt.Printf("ns=worker at=reap status=fail id=%s pid=%s\n", j.ID, j.Pid)
			Models.JobFail(j.ID, fmt.Errorf("Workflow failed"))
		}
	}

	return nil
}
