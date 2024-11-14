package main

import (
	"flag"
	"os"

	"github.com/nrmnqdds/gomaluum-api/cmd/application"
)

type CliFlag struct {
	worker bool
	api    bool
}

func main() {
	flags := getFlags()

	// If no flags are provided, print usage and exit
	if !flags.api && !flags.worker {
		flag.Usage()
		println("\nError: At least one flag (-a/--api or -w/--worker) must be provided")
		os.Exit(1)
	}

	switch {
	case flags.api:
		application.StartEchoServer()
	case flags.worker:
		application.StartAsynqServer()
	}
}

func getFlags() CliFlag {
	var apiFlag bool
	var worker bool

	// Override default usage message
	flag.Usage = func() {
		println("Usage of Gomaluum:")
		println("  -a, --api     Start the API server")
		println("  -w, --worker  Start the worker server")
		println("  -h, --help    Show this help message")
	}

	flag.BoolVar(&apiFlag, "api", false, "start API server")
	flag.BoolVar(&apiFlag, "a", false, "start API server (shorthand)")
	flag.BoolVar(&worker, "worker", false, "start worker server")
	flag.BoolVar(&worker, "w", false, "start worker server (shorthand)")

	flag.Parse()

	return CliFlag{
		api:    apiFlag,
		worker: worker,
	}
}
