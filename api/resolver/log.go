package resolver

import (
	"bufio"
	"context"

	"github.com/convox/convox/pkg/structs"
)

func rackLogs(ctx context.Context, r *Rack, ch chan *RackLog) error {
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
		}

		ch <- &RackLog{line: s.Text()}
	}

	return nil
}

func (rl *RackLog) Line() string {
	return rl.line
}
