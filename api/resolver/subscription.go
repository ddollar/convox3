package resolver

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/convox/console/api/model"
	"github.com/convox/convox/pkg/structs"
	"github.com/graph-gophers/graphql-go"
)

type Log struct {
	line string
}

func (l *Log) Line() string {
	return l.line
}

type InstallLogsArgs struct {
	Oid graphql.ID
	Iid graphql.ID
}

func (r *Root) InstallLogs(ctx context.Context, args InstallLogsArgs) (chan *Log, error) {
	ch := make(chan *Log)

	if _, err := authenticatedInstall(ctx, r.model, string(args.Oid), string(args.Iid)); err != nil {
		return nil, err
	}

	id := string(args.Iid)

	go progressiveLogs(ctx, installLogs(r.model, id), installDone(r.model, id), ch)

	return ch, nil
}

type RackLogsArgs struct {
	Oid   graphql.ID
	Rid   graphql.ID
	Since *int32
}

func (r *Root) RackLogs(ctx context.Context, args RackLogsArgs) (chan *Log, error) {
	o, err := authenticatedOrganization(ctx, r.model, string(args.Oid))
	if err != nil {
		return nil, err
	}

	rr, err := r.model.RackGet(string(args.Rid))
	if err != nil {
		return nil, err
	}

	if rr.Organization != o.ID {
		return nil, fmt.Errorf("invalid organization")
	}

	ch := make(chan *Log)

	go rackLogs(ctx, &Rack{Rack: *rr, model: r.model}, ch)

	return ch, nil
}

type UninstallLogsArgs struct {
	Oid graphql.ID
	Uid graphql.ID
}

func (r *Root) UninstallLogs(ctx context.Context, args UninstallLogsArgs) (chan *Log, error) {
	ch := make(chan *Log)

	if _, err := authenticatedUninstall(ctx, r.model, string(args.Oid), string(args.Uid)); err != nil {
		return nil, err
	}

	id := string(args.Uid)

	go progressiveLogs(ctx, uninstallLogs(r.model, id), uninstallDone(r.model, id), ch)

	return ch, nil
}

type DoneFunc func() (bool, error)
type LogsFunc func() (io.ReadCloser, error)

func installDone(m model.Interface, id string) DoneFunc {
	return func() (bool, error) {
		i, err := m.InstallGet(id)
		if err != nil {
			return true, err
		}

		switch i.Status {
		case "complete", "failed":
			return true, nil
		default:
			return false, nil
		}
	}
}

func installLogs(m model.Interface, id string) LogsFunc {
	return func() (io.ReadCloser, error) {
		return m.InstallLogs(id)
	}
}

func progressiveLogs(ctx context.Context, logs LogsFunc, done DoneFunc, ch chan *Log) error {
	pos := 0

	for {
		time.Sleep(1 * time.Second)

		select {
		case <-ctx.Done():
			return nil
		default:
			r, err := logs()
			if err != nil {
				fmt.Printf("err: %+v\n", err)
				continue
			}
			defer r.Close()

			data, err := ioutil.ReadAll(r)
			if err != nil {
				fmt.Printf("err: %+v\n", err)
				continue
			}

			s := bufio.NewScanner(bytes.NewReader(data[pos:]))

			for s.Scan() {
				ch <- &Log{line: s.Text()}
			}

			if err := s.Err(); err != nil {
				fmt.Printf("err: %+v\n", err)
			}

			pos = len(data)

			d, err := done()
			if err != nil {
				fmt.Printf("err: %+v\n", err)
			}

			fmt.Printf("d: %+v\n", d)

			if d {
				close(ch)
				return nil
			}
		}
	}

}

func rackLogs(ctx context.Context, r *Rack, ch chan *Log) error {
	c, err := r.client(ctx)
	if err != nil {
		return err
	}

	rc, err := c.SystemLogs(structs.LogsOptions{})
	if err != nil {
		return err
	}
	defer rc.Close()

	s := bufio.NewScanner(rc)

	for s.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
			ch <- &Log{line: s.Text()}
		}
	}

	if err := s.Err(); err != nil {
		fmt.Printf("err: %+v\n", err)
	}

	return nil
}

func uninstallDone(m model.Interface, id string) DoneFunc {
	return func() (bool, error) {
		i, err := m.UninstallGet(id)
		if err != nil {
			return true, err
		}

		switch i.Status {
		case "complete", "failed":
			return true, nil
		default:
			return false, nil
		}
	}
}

func uninstallLogs(m model.Interface, id string) LogsFunc {
	return func() (io.ReadCloser, error) {
		return m.UninstallLogs(id)
	}
}
