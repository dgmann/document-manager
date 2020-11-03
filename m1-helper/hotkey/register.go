package hotkey

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/dgmann/document-manager/m1-helper/bdt"
)

func Register(ctx context.Context, fileName, serverUrl string) {
	manager := NewManager()
	defer manager.Close()
	manager.Register(Hotkey{Id: 1, Modifiers: ModAlt + ModCtrl, KeyCode: 'P'})
	keyPresses := manager.Listen()
	println(len(keyPresses))
	for {
		select {
		case x, ok := <-keyPresses:
			if !ok {
				keyPresses = nil
				continue
			}

			println(x)
			f, err := os.Open(fileName)
			if err != nil {
				println("error reading patient file")
				continue
			}

			patient, err := bdt.Parse(f)
			_ = f.Close()

			cmd := exec.Command("explorer", fmt.Sprintf("%s/patient/%s", serverUrl, patient.Id))
			_ = cmd.Run()
		case <-ctx.Done():
			return
		}
	}
}
