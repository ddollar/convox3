package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/convox/console/pkg/worker"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
}

func run() error {
	go handleTermination()

	fmt.Println("worker starting")

	return worker.Work()
}

func handleTermination() {
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, syscall.SIGINT, syscall.SIGTERM)
	<-sigch
	fmt.Println("terminating")
	os.Exit(0)
}
