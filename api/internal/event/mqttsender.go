package event

import (
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/api/pkg/api"
	mqtt "github.com/eclipse/paho.golang/paho"
	"net"
)

type MQTTEventSender struct {
	client *mqtt.Client
}

func NewMQTTEventSender(conn net.Conn) *MQTTEventSender {
	client := mqtt.NewClient(mqtt.ClientConfig{
		Conn: conn,
	})
	return &MQTTEventSender{client: client}
}

func (e *MQTTEventSender) Connect(ctx context.Context) (*mqtt.Connack, error) {
	return e.client.Connect(ctx, &mqtt.Connect{})
}

func (e *MQTTEventSender) Disconnect() error {
	return e.client.Disconnect(&mqtt.Disconnect{})
}

func (e *MQTTEventSender) Send(ctx context.Context, event api.Event) error {
	payload, _ := json.Marshal(event)
	_, err := e.client.Publish(ctx, &mqtt.Publish{
		Topic:   api.RecordTopic,
		QoS:     byte(1),
		Payload: payload,
	})
	return err
}
