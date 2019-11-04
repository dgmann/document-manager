// +build !windows

package hotkey

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) Close() {
}

func (m *Manager) Register(hotkey Hotkey) {
}

func (m *Manager) Listen() <-chan *Hotkey {
	channel := make(chan *Hotkey)
	close(channel)
	return channel
}
