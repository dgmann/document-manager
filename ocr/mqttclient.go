package main

import (
	"context"
	mqtt "github.com/eclipse/paho.golang/paho"
	"net"
)

type MQTTSubscriber struct {
	client *mqtt.Client
}

func NewMQTTSubscriber(conn net.Conn, clientId string) *MQTTSubscriber {
	client := mqtt.NewClient(mqtt.ClientConfig{
		Conn:     conn,
		ClientID: clientId,
	})
	return &MQTTSubscriber{client: client}
}

func (e *MQTTSubscriber) Connect(ctx context.Context) (*mqtt.Connack, error) {
	return e.client.Connect(ctx, &mqtt.Connect{
		ClientID:  e.client.ClientID,
		KeepAlive: 16,
	})
}

func (e *MQTTSubscriber) Disconnect() error {
	return e.client.Disconnect(&mqtt.Disconnect{})
}

func (e *MQTTSubscriber) Subscribe(ctx context.Context, topic string) error {
	_, err := e.client.Subscribe(ctx, &mqtt.Subscribe{
		Subscriptions: map[string]mqtt.SubscribeOptions{
			topic: {
				QoS: byte(1),
			},
		},
	})
	return err
}

func (e *MQTTSubscriber) Router() mqtt.Router {
	return e.client.Router
}
