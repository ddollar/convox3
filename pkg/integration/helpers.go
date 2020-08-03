package integration

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/pkg/errors"
)

func execute(w io.Writer, command string, args ...string) error {
	fmt.Fprintf(w, "%s %s\n", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)

	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}

func gitClone(w io.Writer, url, ref string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", errors.WithStack(err)
	}

	defer os.Chdir(wd)

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	if err := os.Chdir(dir); err != nil {
		return "", errors.WithStack(err)
	}

	if err := execute(w, "git", "init"); err != nil {
		return "", errors.WithStack(err)
	}

	if err := execute(w, "git", "remote", "add", "origin", url); err != nil {
		return "", errors.WithStack(err)
	}

	if err := execute(w, "git", "remote", "update", "origin"); err != nil {
		return "", errors.WithStack(err)
	}

	if err := execute(w, "git", "reset", "--hard", ref); err != nil {
		return "", errors.WithStack(err)
	}

	if err := execute(w, "git", "submodule", "update", "--init"); err != nil {
		return "", errors.WithStack(err)
	}

	return dir, nil
}

func terraformInputs(location string) ([]string, error) {
	res, err := http.Get(location)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var config struct {
		Variables []struct {
			Name       string         `hcl:"name,label"`
			Attributes hcl.Attributes `hcl:",remain"`
		} `hcl:"variable,block"`
	}

	if err := hclsimple.Decode("input.hcl", data, nil, &config); err != nil {
		return nil, errors.WithStack(err)
	}

	ins := []string{}

	for _, v := range config.Variables {
		ins = append(ins, v.Name)
	}

	sort.Strings(ins)

	return ins, nil
}
