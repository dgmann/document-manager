package hotkey

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/dgmann/document-manager/m1-helper/bdt"
)

func Register(ctx context.Context, fileName, serverUrl string) {
	manager := NewManager()
	defer manager.Close()
	openPatientHotkey := Hotkey{Id: 1, Modifiers: ModAlt + ModCtrl, KeyCode: 'P'}
	manager.Register(openPatientHotkey)
	keyPresses := manager.Listen()
	for {
		select {
		case hotkey, ok := <-keyPresses:
			if !ok {
				keyPresses = nil
				continue
			}

			if hotkey.Id != openPatientHotkey.Id {
				log.Printf("unexpected hotkey: %s", hotkey.String())
				continue
			}

			f, err := os.Open(fileName)
			if err != nil {
				log.Printf("error reading patient file %s: %s\n", fileName, err.Error())
				continue
			}

			patient, err := bdt.Parse(f)

			if err = f.Close(); err != nil {
				log.Printf("error closing %s: %s\n", fileName, err.Error())
			}

			cmd := exec.Command("explorer", fmt.Sprintf("%s/patient/%s", serverUrl, patient.Id))
			if err := cmd.Start(); err != nil {
				log.Printf("error running %s\n", cmd.String())
			}
			go func() {
				if err := cmd.Wait(); err != nil {
					log.Printf("error waiting for %s\n", cmd.String())
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}
