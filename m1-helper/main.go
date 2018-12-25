package main

import (
	"context"
	"flag"
	"github.com/dgmann/document-manager/m1-helper/hotkey"
	"github.com/dgmann/document-manager/m1-helper/server"
	"github.com/dgmann/document-manager/m1-helper/service"
)

func main() {
	fileName := flag.String("f", "", "BDT file containing current patient")
	serverUrl := flag.String("s", "http://localhost", "Document-Manager URL")
	flag.Parse()
	if *fileName == "" {
		panic("no file provided")
	}

	if service.Interactive() {
		ctx := context.Background()
		go hotkey.Register(ctx, *fileName, *serverUrl)
		go server.Run(ctx, *fileName)
	} else {
		s, logger, err := service.New(*fileName, *serverUrl)
		if err != nil {
			panic(err)
		}

		if err := s.Run(); err != nil {
			logger.Error(err)
		}
	}
}
