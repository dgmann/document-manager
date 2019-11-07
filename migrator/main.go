package main

import (
	"fmt"
	"github.com/dgmann/document-manager/migrator/http"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("starting...")
	config := http.NewConfig()
	server, err := http.NewServer(config)
	if err != nil {
		logrus.Error(fmt.Errorf("error during startup: %w", err))
		return
	}
	logrus.Info("application started. Listening...")
	logrus.Fatal(server.Run())
}
