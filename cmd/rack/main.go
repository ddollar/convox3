package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/manager"
	"github.com/convox/console/pkg/settings"
	"github.com/convox/console/pkg/storage"
	"github.com/convox/console/pkg/tracker"
	"github.com/convox/convox/pkg/structs"
	"github.com/convox/convox/sdk"
	"github.com/pkg/errors"
)

var (
	m    = model.New(storage.New("dynamo"))
	rack *sdk.Client
)

func init() {
	r, err := sdk.NewFromEnv()
	if err != nil {
		panic(err)
	}
	rack = r
}

type WorkFunc func(id string) error
type SuccessFunc func(id string) error
type FailureFunc func(id string, err error) error

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}

	time.Sleep(5 * time.Second) // give time for write buffers to clear
}

func run() error {
	if len(os.Args) < 3 {
		return errors.WithStack(fmt.Errorf("usage: rack <job> <id>"))
	}

	go timeout(4 * time.Hour)

	switch os.Args[1] {
	case "install":
		return handler(os.Args[2], install, m.InstallSucceed, m.InstallFail)
	case "uninstall":
		return handler(os.Args[2], uninstall, m.UninstallSucceed, m.UninstallFail)
	case "update":
		return handler(os.Args[2], update, m.UpdateSucceed, m.UpdateFail)
	default:
		return errors.WithStack(fmt.Errorf("unknown job: %s", os.Args[1]))
	}
}

func handler(id string, work WorkFunc, success SuccessFunc, failure FailureFunc) error {
	if err := work(id); err != nil {
		return failure(id, err)
	}
	return success(id)
}

func install(id string) error {
	i, err := m.InstallGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	w := writer(i.Key())
	defer w.Close()

	fmt.Fprintf(w, "Installing Rack...\n")

	i.Started = time.Now().UTC()
	i.Status = "running"

	if err := m.InstallSave(i); err != nil {
		return writeError(w, errors.WithStack(err))
	}

	params := map[string]string{}

	for k, v := range i.Params {
		params[k] = v
	}

	mg, err := manager.New(i.Engine, i.RackID)
	if err != nil {
		return err
	}

	if err := mg.Install(i.Name, i.Version, i.Region, params, w); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func timeout(d time.Duration) {
	time.Sleep(d)
	fmt.Fprintf(os.Stderr, "ERROR: timeout\n")
	os.Exit(1)
}

func uninstall(id string) error {
	u, err := m.UninstallGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	w := writer(u.Key())
	defer w.Close()

	fmt.Fprintf(w, "Uninstalling Rack...\n")

	r, err := m.RackGet(u.RackID)
	if err != nil {
		return writeError(w, err)
	}

	mg, err := manager.New(u.Engine, r.ID)
	if err != nil {
		return writeError(w, fmt.Errorf("could not initialize manager"))
	}

	if err := mg.Uninstall(w); err != nil {
		writeError(w, errors.WithStack(err))
		return nil
	}

	if err := m.RackDelete(u.RackID); err != nil {
		return writeError(w, err)
	}

	return nil
}

func update(id string) error {
	u, err := m.UpdateGet(id)
	if err != nil {
		return errors.WithStack(err)
	}

	u.Started = time.Now().UTC()
	u.Status = "running"

	if err := m.UpdateSave(u); err != nil {
		return errors.WithStack(err)
	}

	w := writer(u.Key())
	defer w.Close()

	fmt.Fprintf(w, "Updating Rack...\n")

	r, err := m.RackGet(u.RackID)
	if err != nil {
		return writeError(w, errors.WithStack(err))
	}

	s, err := r.System()
	if err != nil {
		return writeError(w, errors.WithStack(err))
	}

	params := map[string]string{}

	for k, v := range r.Parameters {
		if v != "" {
			params[k] = v
		}
	}

	for k, v := range u.Params {
		params[k] = v
	}

	mg, err := manager.New("v3", r.ID)
	if err != nil {
		return err
	}

	if err := mg.Update(s.Name, u.Version, params, w); err != nil {
		writeError(w, errors.WithStack(err))
	}

	return nil
}

func writer(key string) io.WriteCloser {
	r, w := io.Pipe()

	go write(key, r)

	return w
}

func write(key string, r io.Reader) {
	buf := []byte{}

	data := make([]byte, 1024*1024)

	for {
		n, err := r.Read(data)
		if err == io.EOF {
			return
		}

		buf = append(buf, data[0:n]...)

		fmt.Printf("writing %d bytes\n", len(buf))

		if _, err := rack.ObjectStore(settings.App, key, bytes.NewReader(buf), structs.ObjectStoreOptions{}); err != nil {
			fmt.Printf("err: %+v\n", err)
		}
	}
}

func writeError(w io.Writer, err error) error {
	id := tracker.CaptureError(err)
	fmt.Fprintf(w, "ERROR: we have been notified about a system error: (%s)\n", id)
	return errors.WithStack(err)
}
