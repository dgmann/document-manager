package main

import (
	"log"
	"os"

	"github.com/dgmann/document-manager/internal/cli"
)

func main() {
	l := log.New(os.Stderr, "", 1)
	args := os.Args
	if err := cli.Init(args); err != nil {
		l.Fatalf("error initializing cli: %s", err)
	}
	subArgs := args[2:] // Drop program name and command.
	switch args[1] {
	case "record":
		if err := cli.Record().Execute(subArgs); err != nil {
			l.Fatal(err)
		}
	default:
		l.Fatalf("error: unknown command - %q\n", args[1])
	}
}
