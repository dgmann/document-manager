package service

import (
	"github.com/kardianos/service"
	"log"
)

func New(fileName, serverUrl string) (service.Service, service.Logger, error) {
	svcConfig := &service.Config{
		Name:        "M1Helper",
		DisplayName: "M1-Helper",
		Description: "Dienst für die Verbindung von M1 und DocumentManager",
	}

	prg := newProgram(fileName, serverUrl)
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
	return s, logger, nil
}

func Interactive() bool {
	return service.Interactive()
}
