package job

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"github.com/convox/console/api/model"
	"github.com/convox/console/pkg/storage"
	"github.com/pkg/errors"
)

var (
	m = model.New(storage.New("dynamo"))
)

func Execute(id string) error {
	var j *model.Job

	count := 0

	for {
		jj, err := m.JobGet(id)
		if err != nil {
			return errors.WithStack(err)
		}

		if jj.Pid != "" {
			j = jj
			break
		}

		count++

		if count > 60 {
			return errors.WithStack(fmt.Errorf("could not establish pid"))
		}

		time.Sleep(1 * time.Second)
	}

	j.Started = time.Now().UTC()
	j.Status = "running"

	if err := m.JobSave(j); err != nil {
		return errors.WithStack(err)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// go jobLogs(ctx, id)

	if err := execute(j); err != nil {
		fmt.Printf("ns=job at=Execute status=fail error=%q\n", err)
		m.JobFail(j.ID, fmt.Errorf("Workflow failed"))
		return errors.WithStack(err)
	}

	m.JobSucceed(j.ID)

	// displayJobLogs(j.ID)

	return nil
}

func execute(j *model.Job) error {
	switch j.Kind {
	case "legacy":
		return legacy(j)
	case "merge":
		return merge(j)
	case "review.close":
		return reviewClose(j)
	case "review.open":
		return reviewOpen(j)
	case "review.update":
		return reviewUpdate(j)
	}

	return nil
}

func legacy(j *model.Job) error {
	dir := ""

	desc := j.Params["description"]

	err := task(j, "clone", "Fetching source", func(w io.Writer) error {
		var err error
		dir, err = clone(w, j)
		return errors.WithStack(err)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	rack := ""
	app := ""
	release := ""

	for _, t := range j.Workflow.Tasks {
		switch t.Kind {
		case "build":
			app = t.Params["app_id"]
			manifest := t.Params["manifest"]

			rack, err = rackURL(t.Params["rack_id"])
			if err != nil {
				return errors.WithStack(err)
			}

			err = task(j, "build", fmt.Sprintf("Building %s", app), func(w io.Writer) error {
				var err error
				release, err = build(w, dir, rack, app, manifest, desc, false)
				return errors.WithStack(err)
			})
			if err != nil {
				return errors.WithStack(err)
			}
		case "copy":
			rackto, err := rackURL(t.Params["rack_id_to"])
			if err != nil {
				return errors.WithStack(err)
			}

			appto := t.Params["app_id_to"]

			err = task(j, "copy", fmt.Sprintf("Copying build from %s to %s", app, appto), func(w io.Writer) error {
				var err error
				release, err = copyBuild(w, rack, app, release, rackto, appto)
				rack = rackto
				app = appto
				return errors.WithStack(err)
			})
			if err != nil {
				return errors.WithStack(err)
			}
		case "promote":
			err = task(j, "promote", fmt.Sprintf("Promoting %s on %s", release, app), func(w io.Writer) error {
				return promote(w, rack, app, release)
			})
			if err != nil {
				return errors.WithStack(err)
			}
		case "run":
			err = task(j, "run", fmt.Sprintf("Running command on %s on %s", release, app), func(w io.Writer) error {
				return run(w, rack, app, release, t.Params["service"], t.Params["command"])
			})
			if err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func merge(j *model.Job) error {
	dir := ""

	err := task(j, "clone", "Fetching source", func(w io.Writer) error {
		var err error
		dir, err = clone(w, j)
		return errors.WithStack(err)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	appfrom := ""
	manifest := j.Workflow.Params["manifest"]
	rackfrom := ""
	releasefrom := ""

	desc := j.Params["description"]

	for _, t := range j.Workflow.Tasks {
		switch t.Kind {
		case "deploy":
			rack, err := rackURL(t.Params["rack"])
			if err != nil {
				return errors.WithStack(err)
			}

			app := t.Params["app"]
			release := ""

			if t.Params["test"] == "true" {
				testRelease := ""

				err = task(j, "test-build", fmt.Sprintf("Creating test build on %s", app), func(w io.Writer) error {
					var err error
					testRelease, err = build(w, dir, rack, app, manifest, "test build", true)
					return errors.WithStack(err)
				})
				if err != nil {
					return errors.WithStack(err)
				}

				err = task(j, "test-run", fmt.Sprintf("Running tests on %s", app), func(w io.Writer) error {
					return test(w, rack, app, testRelease)
				})
				if err != nil {
					return errors.WithStack(err)
				}
			}

			if appfrom != "" && releasefrom != "" && rackfrom != "" {
				err = task(j, "copy", fmt.Sprintf("Copying build from %s to %s", appfrom, app), func(w io.Writer) error {
					var err error
					release, err = copyBuild(w, rackfrom, appfrom, releasefrom, rack, app)
					return errors.WithStack(err)
				})
				if err != nil {
					return errors.WithStack(err)
				}
			} else {
				err = task(j, "build", fmt.Sprintf("Building %s", app), func(w io.Writer) error {
					var err error
					release, err = build(w, dir, rack, app, manifest, desc, false)
					rackfrom = rack
					appfrom = app
					releasefrom = release
					return errors.WithStack(err)
				})
				if err != nil {
					return errors.WithStack(err)
				}
			}

			if t.Params["promote"] == "auto" {
				if t.Params["before-service"] != "" && t.Params["before-command"] != "" {
					err = task(j, "run", fmt.Sprintf("Running before hook for %s on %s", release, app), func(w io.Writer) error {
						return run(w, rack, app, release, t.Params["before-service"], t.Params["before-command"])
					})
					if err != nil {
						return errors.WithStack(err)
					}
				}

				err = task(j, "promote", fmt.Sprintf("Promoting %s on %s", release, app), func(w io.Writer) error {
					return promote(w, rack, app, release)
				})
				if err != nil {
					return errors.WithStack(err)
				}

				if t.Params["after-service"] != "" && t.Params["after-command"] != "" {
					err = task(j, "run", fmt.Sprintf("Running after hook for %s on %s", release, app), func(w io.Writer) error {
						return run(w, rack, app, release, t.Params["after-service"], t.Params["after-command"])
					})
					if err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}

	return nil
}

func reviewClose(j *model.Job) error {
	rack, err := rackURL(j.Workflow.Params["rack"])
	if err != nil {
		return errors.WithStack(err)
	}

	app := j.Params["name"]

	err = task(j, "delete", fmt.Sprintf("Deleting %s", app), func(w io.Writer) error {
		return destroy(w, rack, app)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func reviewOpen(j *model.Job) error {
	return reviewUpdate(j)
}

func reviewUpdate(j *model.Job) error {
	app := j.Params["name"]
	dir := ""
	manifest := j.Workflow.Params["manifest"]

	rack, err := rackURL(j.Workflow.Params["rack"])
	if err != nil {
		return errors.WithStack(err)
	}

	if err := convox(ioutil.Discard, nil, rack, "apps", "info", app); err != nil {
		err = task(j, "create", fmt.Sprintf("Creating %s", app), func(w io.Writer) error {
			return create(w, rack, app)
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	err = task(j, "clone", "Fetching source", func(w io.Writer) error {
		var err error
		dir, err = clone(w, j)
		return errors.WithStack(err)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	desc := j.Params["description"]
	release := ""

	if j.Workflow.Params["env"] != "" {
		err = task(j, "env", fmt.Sprintf("Applying test environment to %s", app), func(w io.Writer) error {
			return env(w, rack, app, j.Workflow.Params["env"])
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	err = task(j, "build", fmt.Sprintf("Building %s", app), func(w io.Writer) error {
		var err error
		release, err = build(w, dir, rack, app, manifest, desc, true)
		return errors.WithStack(err)
	})
	if err != nil {
		return errors.WithStack(err)
	}

	if j.Workflow.Params["test"] == "true" {
		err = task(j, "test", fmt.Sprintf("Running tests for %s on %s", release, app), func(w io.Writer) error {
			return test(w, rack, app, release)
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}

	if j.Workflow.Params["demo"] == "true" {
		if j.Workflow.Params["before-service"] != "" && j.Workflow.Params["before-command"] != "" {
			err = task(j, "run", fmt.Sprintf("Running before hook for %s on %s", release, app), func(w io.Writer) error {
				return run(w, rack, app, release, j.Workflow.Params["before-service"], j.Workflow.Params["before-command"])
			})
			if err != nil {
				return errors.WithStack(err)
			}
		}

		err = task(j, "promote", fmt.Sprintf("Promoting %s on %s", release, app), func(w io.Writer) error {
			return promote(w, rack, app, release)
		})
		if err != nil {
			return errors.WithStack(err)
		}

		if j.Workflow.Params["after-service"] != "" && j.Workflow.Params["after-command"] != "" {
			err = task(j, "run", fmt.Sprintf("Running after hook for %s on %s", release, app), func(w io.Writer) error {
				return run(w, rack, app, release, j.Workflow.Params["after-service"], j.Workflow.Params["after-command"])
			})
			if err != nil {
				return errors.WithStack(err)
			}
		}

		ss, err := convoxID(ioutil.Discard, rack, "api", "get", fmt.Sprintf("/apps/%s/services", app))
		if err != nil {
			return errors.WithStack(err)
		}

		var services []struct {
			Domain string
			Name   string
		}

		if err := json.Unmarshal([]byte(ss), &services); err != nil {
			return errors.WithStack(err)
		}

		j.Services = map[string]string{}

		for _, s := range services {
			if s.Domain != "" {
				j.Services[s.Name] = s.Domain
			}
		}
	}

	return nil
}

func build(w io.Writer, dir, rack, app, manifest, description string, development bool) (string, error) {
	args := []string{"build", dir, "--app", app, "--description", description, "--id", "--manifest", manifest}

	if development {
		args = append(args, "--development")
	}

	return convoxID(w, rack, args...)
}

func clone(w io.Writer, j *model.Job) (string, error) {
	s, err := m.WorkflowSource(j.WorkflowID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	repo, err := m.WorkflowRepository(j.WorkflowID)
	if err != nil {
		return "", errors.WithStack(err)
	}

	dir, err := s.RepositoryClone(repo, j.Params["ref"], w)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return dir, nil
}

func copyBuild(w io.Writer, rack, app, release, rackto, appto string) (string, error) {
	data, err := convoxID(nil, rack, "api", "get", fmt.Sprintf("/apps/%s/releases/%s", app, release))
	if err != nil {
		return "", errors.WithStack(err)
	}

	var r struct {
		Build string `json:"build"`
	}

	if err := json.Unmarshal([]byte(data), &r); err != nil {
		return "", errors.WithStack(err)
	}

	tmp, err := ioutil.TempDir("", "")
	if err != nil {
		return "", errors.WithStack(err)
	}

	exp := filepath.Join(tmp, "build.tgz")

	if err := convox(w, nil, rack, "builds", "export", r.Build, "--app", app, "--file", exp); err != nil {
		return "", errors.WithStack(err)
	}

	fmt.Fprintf(w, "\n")

	nr, err := convoxID(w, rackto, "builds", "import", "--app", appto, "--file", exp, "--id")
	if err != nil {
		return "", errors.WithStack(err)
	}

	return nr, nil
}

func create(w io.Writer, rack, app string) error {
	return convox(w, nil, rack, "apps", "create", app)
}

func destroy(w io.Writer, rack, app string) error {
	return convox(w, nil, rack, "apps", "delete", app)
}

func env(w io.Writer, rack, app, env string) error {
	return convox(w, strings.NewReader(env), rack, "env", "set", "--app", app)
}

func promote(w io.Writer, rack, app, release string) error {
	return convox(w, nil, rack, "releases", "promote", release, "--app", app)
}

func run(w io.Writer, rack, app, release, service, command string) error {
	return convox(w, nil, rack, "run", service, command, "--app", app, "--release", release)
}

func test(w io.Writer, rack, app, release string) error {
	return convox(w, nil, rack, "test", "--app", app, "--release", release)
}
