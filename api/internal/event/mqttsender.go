package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgmann/document-manager/api/pkg/api"
	mqtt "github.com/eclipse/paho.golang/paho"
	"net"
)

type MQTTEventSender struct {
	client   *mqtt.Client
	packetId uint16
}

func NewMQTTEventSender(conn net.Conn, clientId string) *MQTTEventSender {
	client := mqtt.NewClient(mqtt.ClientConfig{
		Conn:     conn,
		ClientID: clientId,
	})
	return &MQTTEventSender{client: client}
}

func (e *MQTTEventSender) Connect(ctx context.Context) (*mqtt.Connack, error) {
	return e.client.Connect(ctx, &mqtt.Connect{
		ClientID:  e.client.ClientID,
		KeepAlive: 16,
	})
}

func (e *MQTTEventSender) Disconnect() error {
	return e.client.Disconnect(&mqtt.Disconnect{})
}

func (e *MQTTEventSender) Send(ctx context.Context, event api.Event) error {
	payload, _ := json.Marshal(event)
	_, err := e.client.Publish(ctx, &mqtt.Publish{
		PacketID: e.packetId,
		Topic:    fmt.Sprintf("%s/%s", event.Topic, event.Id),
		QoS:      1,
		Payload:  payload,
		Retain:   false,
	})
	e.packetId += 1
	return err
}
