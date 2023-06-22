package main

import (
	"context"
	"encoding/json"
	"github.com/eclipse/paho.golang/packets"
	mqtt "github.com/eclipse/paho.golang/paho"
	"net"
)

type MQTTClient struct {
	client   *mqtt.Client
	packetId uint16
}

func NewMQTTSubscriber(conn net.Conn, clientId string) *MQTTClient {
	client := mqtt.NewClient(mqtt.ClientConfig{
		Conn:     packets.NewThreadSafeConn(conn),
		ClientID: clientId,
	})
	return &MQTTClient{client: client}
}

func (e *MQTTClient) Connect(ctx context.Context) (*mqtt.Connack, error) {
	return e.client.Connect(ctx, &mqtt.Connect{
		ClientID:  e.client.ClientID,
		KeepAlive: 16,
	})
}

func (e *MQTTClient) Disconnect() error {
	return e.client.Disconnect(&mqtt.Disconnect{})
}

func (e *MQTTClient) Subscribe(ctx context.Context, topic string) error {
	_, err := e.client.Subscribe(ctx, &mqtt.Subscribe{
		Subscriptions: map[string]mqtt.SubscribeOptions{
			topic: {
				QoS: byte(1),
			},
		},
	})
	return err
}

func (e *MQTTClient) Router() mqtt.Router {
	return e.client.Router
}

func (e *MQTTClient) Publish(ctx context.Context, topic string, data any) error {
	payload, _ := json.Marshal(data)
	_, err := e.client.Publish(ctx, &mqtt.Publish{
		PacketID: e.packetId,
		Topic:    topic,
		QoS:      1,
		Payload:  payload,
		Retain:   false,
	})
	e.packetId += 1
	return err
}
