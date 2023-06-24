package event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	log "github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type MQTTEventSender struct {
	Cfg autopaho.ClientConfig
	cm  *autopaho.ConnectionManager
}

func NewMQTTEventSender(broker *url.URL, clientId string) *MQTTEventSender {
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{broker},
		KeepAlive:         30,
		ConnectRetryDelay: 10 * time.Second,
		OnConnectionUp:    func(*autopaho.ConnectionManager, *paho.Connack) { log.Println("MQTT Event Sender: connection up") },
		OnConnectError:    func(err error) { log.Errorf("error whilst attempting connection: %s\n", err) },
		Debug:             paho.NOOPLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID:      clientId,
			OnClientError: func(err error) { log.Errorf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
	return &MQTTEventSender{Cfg: cliCfg}
}

func (m *MQTTEventSender) Connect(ctx context.Context) error {
	cm, err := autopaho.NewConnection(ctx, m.Cfg)
	if err != nil {
		return err
	}
	m.cm = cm
	return nil
}

func (m *MQTTEventSender) Disconnect(ctx context.Context) error {
	return m.cm.Disconnect(ctx)
}

func (m *MQTTEventSender) Send(ctx context.Context, event api.Event) error {
	err := m.cm.AwaitConnection(ctx)
	if err != nil { // Should only happen when context is cancelled
		return fmt.Errorf("publisher done (AwaitConnection: %w)", err)
	}

	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("error marshaling json: %w", err)
	}

	_, err = m.cm.Publish(ctx, &paho.Publish{
		QoS:     byte(1),
		Topic:   fmt.Sprintf("%s/%s", event.Topic, event.Id),
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("error publishing: %s", err)
	}
	return nil
}
