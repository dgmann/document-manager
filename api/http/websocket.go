package http

import (
	"encoding/json"
	"github.com/dgmann/document-manager/api/services"
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func registerWebsocket(router *gin.Engine) {
	m := melody.New()
	ws := services.NewWebSocketService()

	router.GET("/notifications", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleConnect(func(s *melody.Session) {
		ws.AddClient(services.NewClient(s))
	})

	m.HandleDisconnect(func(s *melody.Session) {
		ws.RemoveClient(s)
	})

	go publishEvents(m)
}

func publishEvents(m *melody.Melody) {
	eventService := services.GetEventService()
	events := eventService.Subscribe(services.EventCreated, services.EventDeleted, services.EventUpdated)
	for event := range events {
		data, _ := json.Marshal(event)
		m.Broadcast(data)
	}
}
