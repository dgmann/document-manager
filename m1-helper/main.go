package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgmann/document-manager/m1-helper/hotkey"
	"github.com/dgmann/document-manager/m1-helper/server"
	"github.com/dgmann/document-manager/m1-helper/service"
	service2 "github.com/kardianos/service"
)

func main() {
	fileName := flag.String("f", "./aow_pat.bdt", "BDT file containing current patient")
	serverUrl := flag.String("s", "http://localhost", "Document-Manager URL")
	port := flag.String("p", "3000", "port")
	flag.Parse()
	if *fileName == "" {
		log.Fatal("no BDT file provided")
	}

	s, logger, err := service.New(*fileName, *serverUrl, *port)
	if err != nil {
		log.Fatal(err)
	}

	if service.Interactive() {
		runInteractive(*fileName, *serverUrl, *port)
		if success := installUninstallService(s); !success {
			log.Printf("running interactive on port %s", *port)

		}
	} else {
		if err := s.Run(); err != nil {
			logger.Errorf("error starting service: %w", err)
		}
	}
}

func runInteractive(fileName, serverUrl, port string) {
	ctx := context.Background()
	go hotkey.Register(ctx, fileName, serverUrl)
	server.Run(ctx, fileName, port)
}

func installUninstallService(s service2.Service) bool {
	status, err := s.Status()
	if err != service2.ErrNotInstalled {
		log.Printf("error reading status: %s", err)
	}
	if status == service2.StatusRunning || status == service2.StatusStopped {
		if askForConfirmation("Do you want to uninstall the service?") {
			if err := s.Stop(); err != nil {
				log.Fatalf("error stopping service. Uninstall failed: %s", err)
			}
			log.Printf("successfully stopped service")
			if err := s.Uninstall(); err != nil {
				log.Fatalf("error uninstalling service: %s", err)
			}
			log.Printf("successfully uninstalled service")
			return true
		}
		return false
	} else if err == service2.ErrNotInstalled && askForConfirmation("Service not installed. Do you want to install the service?") {
		if err := s.Install(); err != nil {
			log.Fatalf("error installing service: %s", err)
		}
		if err := s.Start(); err != nil {
			log.Fatalf("error starting service: %s", err)
		}
		fmt.Println("service successfully installed and started")
		return true
	}

	return false
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
