package main

import (
	"fmt"
	"os"
	"time"

	"github.com/convox/console/pkg/job"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.WithStack(fmt.Errorf("must specify a job id"))
	}

	go timeout(4 * time.Hour)

	return job.Execute(os.Args[1])
}

func timeout(d time.Duration) {
	time.Sleep(d)
	fmt.Fprintf(os.Stderr, "ERROR: timeout\n")
	os.Exit(1)
}
