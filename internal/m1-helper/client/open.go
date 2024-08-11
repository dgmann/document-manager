package client

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dgmann/document-manager/pkg/bdt"
)

func OpenPatient(program, fileName, serverUrl string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("error reading patient file %s: %w", fileName, err)
	}

	patient, err := bdt.Parse(f)

	if err = f.Close(); err != nil {
		return fmt.Errorf("error closing %s: %w", fileName, err)
	}

	cmd := exec.Command(program, fmt.Sprintf("%s/patient/%s", serverUrl, patient.Id))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running %s. %w", cmd.String(), err)
	}
	return nil
}
