package main

import (
	"context"
	"github.com/eclipse/paho.golang/autopaho"
)

type MQTTClient struct {
	Cfg autopaho.ClientConfig
	cm  *autopaho.ConnectionManager
}

func (m *MQTTClient) Connect(ctx context.Context) error {
	cm, err := autopaho.NewConnection(ctx, m.Cfg)
	if err != nil {
		return err
	}
	m.cm = cm
	return nil
}

func (m *MQTTClient) Disconnect(ctx context.Context) error {
	return m.Disconnect(ctx)
}
