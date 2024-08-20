package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/dgmann/document-manager/internal/directorywatcher"
	"github.com/dgmann/document-manager/pkg/log"
	"github.com/namsral/flag"
)

var logger = log.Logger

func main() {
	configPath := flag.String("c", "./watcher.yaml", "specifies path to watcher config file")
	flag.Parse()
	logger.Info("loading config file", "file", *configPath)
	config, err := func() (config directorywatcher.Config, err error) {
		configFile, err := os.Open(*configPath)
		if err != nil {
			err = fmt.Errorf("error opening config file: %w", err)
			return
		}
		c, err := io.ReadAll(configFile)
		if err != nil {
			err = fmt.Errorf("error reading config file: %w", err)
			return
		}
		config, err = directorywatcher.LoadConfig(c)
		if err != nil {
			err = fmt.Errorf("error parsing config file: %w", err)
			return
		}
		return
	}()
	if err != nil {
		logger.Error("error loading config file", log.ErrAttr(err))
		os.Exit(1)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		cancel()
	}()
	directorywatcher.Watch(ctx, config)
}
