package hotkey

import (
	"golang.org/x/sys/windows"
	"fmt"
	"unsafe"
	"time"
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
	m.user32.Release()
}

func (m *Manager) Register(hotkey Hotkey) {
	r1, _, err := m.reghotkey.Call(
		0, uintptr(hotkey.Id), uintptr(hotkey.Modifiers), uintptr(hotkey.KeyCode))
	if r1 == 1 {
		m.hotkeys[int16(hotkey.Id)] = &hotkey
		fmt.Println("Registered", hotkey)
	} else {
		fmt.Println("Failed to register", hotkey, ", error:", err)
	}
}

func (m *Manager) Listen() <-chan *Hotkey {
	channel := make(chan *Hotkey)
	go func() {
		peekmsg := m.user32.MustFindProc("PeekMessageW")
		for {
			var msg = &MSG{}
			peekmsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

			// Registered id is in the WPARAM field:
			if id := msg.WPARAM; id != 0 {
				fmt.Println("Hotkey pressed:", m.hotkeys[id])
				channel <- m.hotkeys[id]
			}

			time.Sleep(time.Millisecond * 50)
		}
	}()
	return channel
}
