package service

import (
	"context"

	"github.com/dgmann/document-manager/m1-helper/hotkey"
	"github.com/dgmann/document-manager/m1-helper/server"
	"github.com/kardianos/service"
)

var logger service.Logger

type program struct {
	cancel    context.CancelFunc
	ctx       context.Context
	fileName  string
	serverUrl string
	port      string
}

func newProgram(fileName, serverUrl string, port string) *program {
	return &program{
		ctx:       context.Background(),
		fileName:  fileName,
		serverUrl: serverUrl,
		port:      port,
	}
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	ctx, cancel := context.WithCancel(p.ctx)
	p.cancel = cancel
	go hotkey.Register(ctx, p.fileName, p.serverUrl)
	go server.Run(ctx, p.fileName, p.port)
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	p.cancel()
	return nil
}
