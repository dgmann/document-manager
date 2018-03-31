package services

import (
	"gopkg.in/olahol/melody.v1"
	"sync"
)

type Client struct {
	Name    string
	session *melody.Session
}

func NewClient(session *melody.Session) *Client {
	return &Client{Name: "", session: session}
}

type WebsocketService struct {
	clients map[*melody.Session]*Client
	mutex   *sync.Mutex
}

func NewWebSocketService() *WebsocketService {
	return &WebsocketService{clients: make(map[*melody.Session]*Client), mutex: new(sync.Mutex)}
}

func (ws *WebsocketService) AddClient(client *Client) {
	ws.mutex.Lock()
	ws.clients[client.session] = client
	ws.mutex.Unlock()
}

func (ws *WebsocketService) RemoveClient(session *melody.Session) {
	ws.mutex.Lock()
	delete(ws.clients, session)
	ws.mutex.Unlock()
}

func (ws *WebsocketService) GetClient(session *melody.Session) *Client {
	return ws.clients[session]
}
