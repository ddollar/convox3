package job

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"time"

	shellquote "github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

func convox(w io.Writer, r io.Reader, rack string, args ...string) error {
	fmt.Fprintf(w, "$ convox %s\n", shellquote.Join(args...))

	cmd := exec.Command("convox", args...)

	cmd.Env = []string{fmt.Sprintf("RACK_URL=%s", rack)}

	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = w

	return cmd.Run()
}

func convoxID(w io.Writer, rack string, args ...string) (string, error) {
	if w != nil {
		fmt.Fprintf(w, "$ convox %s\n", shellquote.Join(args...))
	}

	cmd := exec.Command("convox", args...)

	var buf bytes.Buffer

	cmd.Env = []string{fmt.Sprintf("RACK_URL=%s", rack)}
	cmd.Stdout = &buf
	cmd.Stderr = w

	if err := cmd.Run(); err != nil {
		return "", errors.WithStack(err)
	}

	return buf.String(), nil
}

func displayJobLogs(id string) {
	j, err := m.JobGet(id)
	if err == nil {
		fmt.Printf("job: %s\n", id)
		for i, t := range j.Tasks {
			fmt.Printf("task %d: [%s] %s\n", i, t.Status, t.Title)
			if t.Status != "running" {
				continue
			}
			out, err := m.TaskOutput(t)
			if err == nil {
				fmt.Println(out)
			}
		}
		fmt.Printf("\n------------------------\n")
	}
}

func jobLogs(ctx context.Context, id string) {
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick:
			displayJobLogs(id)
		}
	}
}

func rackURL(id string) (string, error) {
	r, err := m.RackGet(id)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return r.URL()
}
