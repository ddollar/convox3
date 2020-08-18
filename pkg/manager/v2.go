package manager

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

type V2 struct {
	env map[string]string
	id  string
}

func NewV2(rid string) (*V2, error) {
	ri, err := m.RackRuntime(rid)
	if err != nil {
		return nil, err
	}

	creds, err := ri.Credentials()
	if err != nil {
		return nil, err
	}

	m := &V2{
		env: creds,
		id:  rid,
	}

	return m, nil
}

func (v *V2) Install(name, version, region string, params map[string]string, output io.Writer) error {
	home, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.WithStack(err)
	}

	args := []string{"rack", "install", "aws", "--name", name, "--raw"}

	cmd := exec.Command("convox2", args...)

	cmd.Env = []string{
		fmt.Sprintf("HOME=%s", home),
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("AWS_REGION=%s", region),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", region),
	}

	for k, v := range v.env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	or, err := cmd.StdoutPipe()
	if err != nil {
		return errors.WithStack(err)
	}

	cmd.Stderr = output

	if err := cmd.Start(); err != nil {
		return errors.WithStack(err)
	}

	if err := v.stream(output, or); err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	data, err := ioutil.ReadFile(filepath.Join(home, ".convox", "auth"))
	if err != nil {
		return errors.WithStack(err)
	}

	var auth map[string]string

	if err := json.Unmarshal(data, &auth); err != nil {
		return errors.WithStack(err)
	}

	r, err := m.RackGet(v.id)
	if err != nil {
		return errors.WithStack(err)
	}

	for k, v := range auth {
		r.Host = k
		r.Password = v
	}

	if err := m.RackSave(r); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (v *V2) Uninstall(output io.Writer) error {
	r, err := m.RackGet(v.id)
	if err != nil {
		return err
	}

	s, err := r.System()
	if err != nil {
		return errors.WithStack(err)
	}

	args := []string{"rack", "uninstall", "aws", s.Name, "--force"}

	cmd := exec.Command("convox2", args...)

	cmd.Env = []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	for k, v := range v.env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stdout = output
	cmd.Stderr = output

	if err := cmd.Start(); err != nil {
		return errors.WithStack(err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (v *V2) Update(name, version string, params map[string]string, output io.Writer) error {
	return fmt.Errorf("v2 racks are self updating")
}

func (v *V2) stream(w io.Writer, r io.Reader) error {
	var current int

	dec := json.NewDecoder(r)

	var step struct {
		Stack   string
		Current int
		Total   int
	}

	for dec.More() {
		if err := dec.Decode(&step); err != nil {
			return errors.WithStack(err)
		}

		if step.Current > current {
			fmt.Fprintf(w, "Progress: %d%%\n", (step.Current * 100 / step.Total))
			current = step.Current
		}
	}

	return nil
}
