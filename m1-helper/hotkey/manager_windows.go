// +build windows
package hotkey

import (
	"fmt"
	"log"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

type Manager struct {
	user32    *windows.DLL
	reghotkey *windows.Proc
	hotkeys   map[int16]*Hotkey
}

func NewManager() *Manager {
	user32 := windows.MustLoadDLL("user32")
	return &Manager{
		user32:    user32,
		reghotkey: user32.MustFindProc("RegisterHotKey"),
		hotkeys:   make(map[int16]*Hotkey),
	}
}

func (m *Manager) Close() {
	if err := m.user32.Release(); err != nil {
		log.Fatal(err)
	}
}

func (m *Manager) Register(hotkey Hotkey) {
	r1, _, err := m.reghotkey.Call(
		0, uintptr(hotkey.Id), uintptr(hotkey.Modifiers), uintptr(hotkey.KeyCode))
	if r1 == 1 {
		m.hotkeys[int16(hotkey.Id)] = &hotkey
		fmt.Println("Registered", hotkey.String())
	} else {
		fmt.Println("Failed to register", hotkey.String(), ", error:", err)
	}
}

func (m *Manager) Listen() <-chan *Hotkey {
	channel := make(chan *Hotkey)
	go func() {
		peekmsg := m.user32.MustFindProc("PeekMessageW")
		for {
			var msg = &MSG{}
			_, _, _ = peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

			// Registered id is in the WPARAM field:
			if id := msg.WPARAM; id != 0 {
				pressed, ok := m.hotkeys[id]
				if !ok {
					log.Printf("hotkey with id %d not found", id)
				}
				select {
				case channel <- pressed:
					log.Println("Hotkey pressed:", pressed)
				default:
					log.Println("Hotkey pressed but channel buffer is full:", pressed)
				}

			}

			time.Sleep(time.Millisecond * 50)
		}
		fmt.Println("Stop listening for hotkeys")
		close(channel)
	}()
	return channel
}
