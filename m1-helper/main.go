package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/MakeNowJust/hotkey"
	"github.com/dgmann/document-manager/m1-helper/client"
	"github.com/dgmann/document-manager/m1-helper/icon"
	"github.com/dgmann/document-manager/m1-helper/server"
	"github.com/dgmann/document-manager/m1-helper/service"
	"github.com/getlantern/systray"
	service2 "github.com/kardianos/service"
	"log"
	"os"
	"strings"
)

var NotInstalledGer = "Der angegebene Dienst ist kein installierter Dienst."

func main() {
	fileName := flag.String("f", lookupEnvOrString("M1_BDT_FILE", "./aow_pat.bdt"), "BDT file containing current patient. Env: M1_BDT_FILE")
	serverURL := flag.String("s", lookupEnvOrString("DOCUMENT_MANAGER_URL", "http://localhost"), "Document-Manager URL")
	port := flag.String("p", lookupEnvOrString("M1_HELPER_PORT", "3000"), "port")
	interactive := flag.Bool("i", false, "run in interactive mode")
	openCommand := flag.String("c", "explorer", "default command to open Document-Manager. Default: explorer [SERVERURL]/patient/[ID]")
	flag.Parse()
	if *fileName == "" {
		log.Fatal("no BDT file provided")
	}

	s, logger, err := service.New(*openCommand, *fileName, *serverURL, *port)
	if err != nil {
		log.Fatal(err)
	}

	if service.Interactive() {
		if *interactive || !installUninstallService(s) {
			log.Printf("running interactive on port %s", *port)
			runInteractive(*openCommand, *fileName, *serverURL, *port)
		}
	} else {
		if err := s.Run(); err != nil {
			logger.Errorf("error starting service: %w", err)
		}
	}
}

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func installUninstallService(s service2.Service) bool {
	status, err := s.Status()
	if err != nil {
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
	} else if askForConfirmation("Service not installed. Do you want to install the service?") {
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

func runInteractive(openCommand, fileName, serverUrl, port string) {
	ctx, cancel := context.WithCancel(context.Background())

	openPatientCmd := func() {
		if err := client.OpenPatient(openCommand, fileName, serverUrl); err != nil {
			log.Println(err)
		}
	}

	manager := hotkey.New()
	manager.Register(hotkey.Alt+hotkey.Ctrl, 'P', func() {
		go openPatientCmd()
	})
	go server.Run(ctx, fileName, port)

	onExit := func() {
		cancel()
	}
	systray.Run(onReady(openPatientCmd), onExit)
}

func onReady(openPatientCmd func()) func() {
	return func() {
		systray.SetIcon(icon.Data)
		systray.SetTitle("Document Manager Helper")
		systray.SetTooltip("Document Manager Helper")
		mOpen := systray.AddMenuItem("Patient", "Derzeitigen Patienten öffnen")
		go func() {
			<-mOpen.ClickedCh
			openPatientCmd()
		}()
		mQuit := systray.AddMenuItem("Beenden", "Helper beenden")
		go func() {
			<-mQuit.ClickedCh
			systray.Quit()
		}()
	}
}
