package resolver

import (
	"bufio"
	"context"
	"fmt"

	"github.com/convox/convox/pkg/structs"
)

func rackLogs(ctx context.Context, r *Rack, ch chan *RackLog) error {
	c, err := r.client(ctx)
	if err != nil {
		return err
	}

	fmt.Printf("c: %+v\n", c)

	rc, err := c.SystemLogs(structs.LogsOptions{})
	if err != nil {
		return err
	}
	defer rc.Close()

	fmt.Printf("rc: %+v\n", rc)

	s := bufio.NewScanner(rc)

	for s.Scan() {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Printf("s.Text(): %+v\n", s.Text())

		ch <- &RackLog{line: s.Text()}
	}

	fmt.Println("done")

	return nil
}

func (rl *RackLog) Line() string {
	return rl.line
}
