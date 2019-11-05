package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/event"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"gopkg.in/olahol/melody.v1"
	"net/http"
	"sync"
)

func getWebsocketHandler(subscriber event.Subscriber) http.Handler {
	r := chi.NewRouter()
	m := melody.New()
	ws := NewWebSocketService()

	r.Get("/", func(w http.ResponseWriter, req *http.Request) {

		if err := m.HandleRequest(w, req); err != nil {
			if _, werr := w.Write([]byte(err.Error())); werr != nil {
				logrus.Error(werr)
			}
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		ws.AddClient(NewClient(s))
	})

	m.HandleDisconnect(func(s *melody.Session) {
		ws.RemoveClient(s)
	})

	go publishEvents(m, subscriber)
	return r
}

func publishEvents(m *melody.Melody, subscriber event.Subscriber) {
	events := subscriber.Subscribe(event.Created, event.Deleted, event.Updated)
	for e := range events {
		data, _ := json.Marshal(e)
		err := m.Broadcast(data)
		if err != nil {
			logrus.Debug(err)
		}
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
