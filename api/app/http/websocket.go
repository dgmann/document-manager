package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
	"sync"
)

func registerWebsocket(router *gin.Engine, subscriber app.Subscriber) {
	m := melody.New()
	ws := NewWebSocketService()

	router.GET("/notifications", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		ws.AddClient(NewClient(s))
	})

	m.HandleDisconnect(func(s *melody.Session) {
		ws.RemoveClient(s)
	})

	go publishEvents(m, subscriber)
}

func publishEvents(m *melody.Melody, subscriber app.Subscriber) {
	events := subscriber.Subscribe(app.EventCreated, app.EventDeleted, app.EventUpdated)
	for event := range events {
		data, _ := json.Marshal(event)
		m.Broadcast(data)
	}
}

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
