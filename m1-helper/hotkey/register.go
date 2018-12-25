package hotkey

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/m1-helper/m1"
	"io/ioutil"
	"log"
	"os/exec"
)

func Register(ctx context.Context, fileName, serverUrl string) {
	manager := NewManager()
	manager.Register(Hotkey{Id: 1, Modifiers: ModAlt + ModCtrl, KeyCode: 'P'})
	keyPresses := manager.Listen()
	for {
		select {
		case <-keyPresses:
			f, err := ioutil.ReadFile(fileName)
			if err != nil {
				println("error reading patient file")
			}

			patient, err := m1.Parse(f)

			cmd := exec.Command(fmt.Sprintf("%s/patient/%s", serverUrl, patient.Id))
			err = cmd.Run()

			if err != nil {
				fmt.Printf("an error occurred: %s\n", err)
				log.Fatal(err)
			}
		case <-ctx.Done():
			return
		}
	}
}
