package service

import (
	"context"
	"log"

	"github.com/MakeNowJust/hotkey"
	"github.com/dgmann/document-manager/internal/m1-helper/client"
	"github.com/dgmann/document-manager/internal/m1-helper/server"
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct {
	cancel        context.CancelFunc
	ctx           context.Context
	hotkeyManager *hotkey.Manager
	openCmd       string
	fileName      string
	serverUrl     string
	port          string
}

func newProgram(openCmd, fileName, serverUrl string, port string) *program {
	return &program{
		ctx:           context.Background(),
		openCmd:       openCmd,
		fileName:      fileName,
		serverUrl:     serverUrl,
		port:          port,
		hotkeyManager: hotkey.New(),
	}
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	ctx, cancel := context.WithCancel(p.ctx)
	p.cancel = cancel
	p.hotkeyManager.Register(hotkey.Alt+hotkey.Ctrl, 'P', func() {
		go func() {
			if err := client.OpenPatient(p.openCmd, p.fileName, p.serverUrl); err != nil {
				log.Println(err)
			}
		}()
	})
	go server.Run(ctx, p.fileName, p.port)
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	p.cancel()
	p.hotkeyManager.Stop()
	return nil
}
