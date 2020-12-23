package service

import (
	"fmt"

	"github.com/kardianos/service"
)

func New(fileName, serverUrl string, port string) (service.Service, service.Logger, error) {
	svcConfig := &service.Config{
		Name:        "M1Helper",
		DisplayName: "M1-Helper",
		Description: "Dienst für die Verbindung von M1 und DocumentManager",
		Arguments:   []string{"-f", fileName, "-s", serverUrl, "-p", port},
	}

	prg := newProgram(fileName, serverUrl, port)
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating service: %w", err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating logger: %w", err)
	}
	return s, logger, nil
}

func Interactive() bool {
	return service.Interactive()
}
