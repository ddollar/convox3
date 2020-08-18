package worker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/queue"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/console/pkg/storage"
	"github.com/convox/convox/pkg/options"
	"github.com/convox/convox/pkg/structs"
	"github.com/convox/convox/sdk"
	"github.com/convox/logger"
	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

var (
	Logger = logger.New("ns=worker")
	Models model.Interface
	Rack   structs.Provider
)

func init() {
	Models = model.New(storage.New("dynamo"))

	r, err := sdk.NewFromEnv()
	if err != nil {
		panic(err)
	}
	Rack = r
}

func Cancel(jid string) error {
	j, err := Models.JobGet(jid)
	if err != nil {
		return errors.WithStack(err)
	}

	if j.Pid != "" {
		if err := stop(j.Pid); err != nil {
			return errors.WithStack(err)
		}
	}

	j.Pid = ""

	if err := Models.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func Work() error {
	Logger.Logf("at=start")

	go reap()

	q := queue.New(settings.WorkerQueue)

	for {
		work, err := q.Dequeue(wait)
		if err != nil {
			return errors.WithStack(err)
		}

		switch work["type"] {
		case "install":
			if err := workInstall(work); err != nil {
				return errors.WithStack(err)
			}
		case "job":
			if err := workJob(work); err != nil {
				return errors.WithStack(err)
			}
		case "uninstall":
			if err := workUninstall(work); err != nil {
				return errors.WithStack(err)
			}
		case "update":
			if err := workUpdate(work); err != nil {
				return errors.WithStack(err)
			}
		default:
			Logger.Logf("error=unknown type=%s", work["type"])
		}
	}

	return nil
}

func workInstall(work map[string]string) error {
	id := work["id"]

	fmt.Printf("ns=worker at=install id=%q\n", id)

	if id == "" {
		return errors.WithStack(fmt.Errorf("install has no id"))
	}

	i, err := Models.InstallGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if i.Status != "starting" {
		fmt.Printf("install in %s status, expected starting: %s\n", i.Status, id)
		return nil
	}

	pid, err := run("rack", "install", id)
	if err != nil {
		return errors.WithStack(err)
	}

	i.Pid = pid

	if err := Models.InstallSave(i); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func workJob(work map[string]string) error {
	id := work["id"]

	fmt.Printf("ns=worker at=job id=%q\n", id)

	if id == "" {
		return errors.WithStack(fmt.Errorf("job has no id"))
	}

	j, err := Models.JobGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if j.Status != "starting" {
		fmt.Printf("job in %s status, expected starting: %s\n", j.Status, id)
		return nil
	}

	pid, err := run("job", id)
	if err != nil {
		return errors.WithStack(err)
	}

	j.Pid = pid

	if err := Models.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func workUninstall(work map[string]string) error {
	id := work["id"]

	fmt.Printf("ns=worker at=uninstall id=%q\n", id)

	if id == "" {
		return errors.WithStack(fmt.Errorf("uninstall has no id"))
	}

	u, err := Models.UninstallGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if u.Status != "starting" {
		fmt.Printf("uninstall in %s status, expected starting: %s\n", u.Status, id)
		return nil
	}

	pid, err := run("rack", "uninstall", id)
	if err != nil {
		return errors.WithStack(err)
	}

	u.Pid = pid

	if err := Models.UninstallSave(u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func workUpdate(work map[string]string) error {
	id := work["id"]

	fmt.Printf("ns=worker at=update id=%q\n", id)

	if id == "" {
		return errors.WithStack(fmt.Errorf("update has no id"))
	}

	u, err := Models.UpdateGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	if u.Status != "starting" {
		fmt.Printf("update in %s status, expected starting: %s\n", u.Status, id)
		return nil
	}

	pid, err := run("rack", "update", id)
	if err != nil {
		return errors.WithStack(err)
	}

	u.Pid = pid

	if err := Models.UpdateSave(u); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func run(command string, args ...string) (string, error) {
	if settings.Development {
		return runLocal(command, args...)
	} else {
		return runRack(command, args...)
	}
}

func runLocal(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true, Pgid: 0}

	if err := cmd.Start(); err != nil {
		return "", errors.WithStack(err)
	}

	go cmd.Wait()

	return fmt.Sprintf("%d", cmd.Process.Pid), nil
}

func runRack(command string, args ...string) (string, error) {
	ps, err := Rack.ProcessRun(os.Getenv("APP"), "worker", structs.ProcessRunOptions{
		Command: options.String(fmt.Sprintf("%s %s", command, shellquote.Join(args...))),
	})
	if err != nil {
		return "", errors.WithStack(err)
	}

	return ps.Id, nil
}

func stop(pid string) error {
	if settings.Development {
		return stopLocal(pid)
	} else {
		return stopRack(pid)
	}
}

func stopLocal(pid string) error {
	// TODO: implement via queue message
	return nil
}

func stopRack(pid string) error {
	Rack.ProcessStop(os.Getenv("APP"), pid)

	return nil
}

var waitLock sync.Mutex

func wait(work map[string]string) (bool, error) {
	waitLock.Lock()
	defer waitLock.Unlock()

	switch work["type"] {
	case "install":
		return waitInstall(work)
	case "job":
		return waitJob(work)
	case "update":
		return waitUpdate(work)
	case "uninstall":
		return waitUninstall(work)
	default:
		return false, errors.WithStack(fmt.Errorf("unknown work type: %s", work["type"]))
	}
}

func waitInstall(work map[string]string) (bool, error) {
	i, err := Models.InstallGet(work["id"])
	if err != nil {
		return true, errors.WithStack(err)
	}

	i.Status = "starting"

	if err := Models.InstallSave(i); err != nil {
		return true, nil
	}

	return false, nil
}

func waitJob(work map[string]string) (bool, error) {
	w, err := Models.WorkflowGet(work["workflow"])
	if err != nil {
		return true, errors.WithStack(err)
	}

	o, err := Models.OrganizationGet(w.OrganizationId)
	if err != nil {
		return true, errors.WithStack(err)
	}

	jss, err := Models.JobListByStatus("starting")
	if err != nil {
		return true, errors.WithStack(err)
	}

	jsr, err := Models.JobListByStatus("running")
	if err != nil {
		return true, errors.WithStack(err)
	}

	js := append(jss, jsr...)

	oc := 0

	for _, j := range js {
		if j.OrganizationID == w.OrganizationId {
			oc++
		}
	}

	jc := o.JobConcurrency()

	Logger.Logf("at=wait org=%s job=%s running=%d concurrency=%d", w.OrganizationId, work["id"], oc, jc)

	if oc >= o.JobConcurrency() {
		Logger.Logf("at=wait org=%s job=%s status=wait reason=concurrency", w.OrganizationId, work["id"])
		return true, nil
	}

	for _, j := range js {
		switch strings.Split(j.Kind, ".")[0] {
		case "merge":
			if j.WorkflowID == work["workflow"] {
				Logger.Logf("at=wait org=%s job=%s status=wait reason=merge", w.OrganizationId, work["id"])
				return true, nil
			}
		case "review":
			if j.WorkflowID == work["workflow"] && j.Params["name"] == work["name"] {
				Logger.Logf("at=wait org=%s job=%s status=wait reason=review", w.OrganizationId, work["id"])
				return true, nil
			}
		}
	}

	Logger.Logf("at=wait org=%s job=%s status=proceed", w.OrganizationId, work["id"])

	j, err := Models.JobGet(work["id"])
	if err != nil {
		return true, errors.WithStack(err)
	}

	j.Status = "starting"

	if err := Models.JobSave(j); err != nil {
		return true, nil
	}

	return false, nil
}

func waitUninstall(work map[string]string) (bool, error) {
	u, err := Models.UninstallGet(work["id"])
	if err != nil {
		return true, errors.WithStack(err)
	}

	u.Status = "starting"

	if err := Models.UninstallSave(u); err != nil {
		return true, nil
	}

	return false, nil
}

func waitUpdate(work map[string]string) (bool, error) {
	u, err := Models.UpdateGet(work["id"])
	if err != nil {
		return true, errors.WithStack(err)
	}

	u.Status = "starting"

	if err := Models.UpdateSave(u); err != nil {
		return true, nil
	}

	return false, nil
}
