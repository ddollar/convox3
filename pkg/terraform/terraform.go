package terraform

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/convox/console/pkg/settings"
	"github.com/convox/logger"
	"github.com/pkg/errors"
)

const (
	terraformDownload = "https://releases.hashicorp.com/terraform/%s/terraform_%s_linux_amd64.zip"
)

type Terraform struct {
	Writer io.Writer

	backend  string
	logger   *logger.Logger
	provider string
}

// TODO update terraform

func New(provider, backend string) (*Terraform, error) {
	t := &Terraform{
		Writer:   os.Stdout,
		backend:  backend,
		logger:   logger.New("ns=terraform"),
		provider: provider,
	}

	t.logger.Logf("provider=%s", provider)

	return t, nil
}

func (t *Terraform) Install(name, version string, env map[string]string, params map[string]string) error {
	args := []string{"rack", "install", t.provider, name, "-v", version}

	for k, v := range params {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	t.logger.At("install").Logf("args=%v", args)

	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.WithStack(err)
	}

	if err := t.convox(tmp, env, args...); err != nil {
		fmt.Fprintf(t.Writer, "install failed, cleaning up\n")

		if err := t.convox(tmp, env, "rack", "uninstall", name); err != nil {
			return errors.WithStack(fmt.Errorf("install cleanup failed: %v", err))
		}

		return errors.WithStack(fmt.Errorf("install failed: %v", err))
	}

	return nil
}

func (t *Terraform) Uninstall(name, version string, env, params map[string]string) error {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.WithStack(err)
	}

	if err := t.downloadTerraform(tmp); err != nil {
		return errors.WithStack(err)
	}

	args := []string{"rack", "install", t.provider, name, "--prepare", "-v", version}

	for k, v := range params {
		switch k {
		case "release":
		default:
			args = append(args, fmt.Sprintf("%s=%s", k, v))
		}
	}

	t.logger.At("uninstall").Logf("name=%q version=%q args=%v", name, version, args)

	if err := t.convox(tmp, env, args...); err != nil {
		return errors.WithStack(err)
	}

	if err := t.convox(tmp, env, "rack", "uninstall", name); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *Terraform) Update(name, version string, env map[string]string, params map[string]string) error {
	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return errors.WithStack(err)
	}

	if err := t.downloadTerraform(tmp); err != nil {
		return errors.WithStack(err)
	}

	args := []string{"rack", "install", t.provider, name, "-v", version}

	for k, v := range params {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	t.logger.At("update").Logf("args=%v", args)

	if err := t.convox(tmp, env, args...); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *Terraform) convox(home string, env map[string]string, args ...string) error {
	cmd := exec.Command("convox", args...)

	cmd.Env = t.env(home)

	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	cmd.Stdout = t.Writer
	cmd.Stderr = t.Writer

	if err := cmd.Start(); err != nil {
		return errors.WithStack(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go handleTermination(ctx, cmd)
	defer cancel()

	if err := cmd.Wait(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *Terraform) downloadTerraform(dir string) error {
	data, err := t.state()
	if err != nil {
		return errors.WithStack(fmt.Errorf("could not fetch state: %s", err))
	}

	var state struct {
		TerraformVersion string `json:"terraform_version"`
	}

	if err := json.Unmarshal(data, &state); err != nil {
		return errors.WithStack(err)
	}

	if state.TerraformVersion == "" {
		return errors.WithStack(fmt.Errorf("unknown terraform version in state"))
	}

	if err := downloadTerraformVersion(state.TerraformVersion, dir); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (t *Terraform) env(home string) []string {
	return []string{
		fmt.Sprintf("CONVOX_TERRAFORM_BACKEND=%s", t.backend),
		fmt.Sprintf("CONVOX_TERRAFORM_BACKEND_INSECURE=%t", settings.Development),
		fmt.Sprintf("HOME=%s", home),
		fmt.Sprintf("PATH=%s:%s", home, os.Getenv("PATH")),
	}
}

func (t *Terraform) state() ([]byte, error) {
	ht := &http.Transport{}

	if settings.Development {
		ht.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	c := &http.Client{Transport: ht}

	res, err := c.Get(fmt.Sprintf("%s/state", t.backend))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return data, nil
}

func downloadTerraformVersion(version, dir string) error {
	fmt.Printf("version: %+v\n", version)

	u := fmt.Sprintf(terraformDownload, version, version)

	res, err := http.Get(u)
	if err != nil {
		return errors.WithStack(err)
	}
	defer res.Body.Close()

	fd, err := os.Create(filepath.Join(dir, "terraform.zip"))
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := io.Copy(fd, res.Body); err != nil {
		return errors.WithStack(err)
	}

	cmd := exec.Command("unzip", "terraform.zip")
	cmd.Dir = dir

	if err := cmd.Run(); err != nil {
		return errors.WithStack(fmt.Errorf("could not unzip terraform"))
	}

	return nil
}

func handleTermination(ctx context.Context, cmd *exec.Cmd) {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	for {
		select {
		case <-ctx.Done():
			return
		case s := <-ch:
			if p := cmd.Process; p != nil {
				if err := p.Signal(s); err != nil {
					fmt.Printf("ERROR: %s\n", err)
				}
			}
		}
	}
}
