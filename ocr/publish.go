package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"log"
	"net/url"
	"time"
)

type Publisher struct {
	MQTTClient
}

func NewPublisher(broker *url.URL, clientId string) *Publisher {
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{broker},
		KeepAlive:         30,
		ConnectRetryDelay: 10 * time.Second,
		OnConnectionUp:    func(*autopaho.ConnectionManager, *paho.Connack) { log.Println("Publisher: mqtt connection up") },
		OnConnectError:    func(err error) { log.Printf("Publisher: error whilst attempting connection: %s\n", err) },
		Debug:             paho.NOOPLogger{},
		ClientConfig: paho.ClientConfig{
			ClientID:      clientId,
			OnClientError: func(err error) { log.Printf("Publisher: server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("Publisher: server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("Publisher: server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
	return &Publisher{MQTTClient{Cfg: cliCfg}}
}

func (p *Publisher) Run(ctx context.Context, topic string, ocrRequestChan <-chan OCRRequest) error {
	for msg := range ocrRequestChan {
		// AwaitConnection will return immediately if connection is up; adding this call stops publication whilst
		// connection is unavailable.
		err := p.cm.AwaitConnection(ctx)
		if err != nil { // Should only happen when context is cancelled
			return fmt.Errorf("publisher done (AwaitConnection: %w)", err)
		}

		payload, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("error marshaling json: %w", err)
		}

		pr, err := p.cm.Publish(ctx, &paho.Publish{
			QoS:     byte(1),
			Topic:   topic,
			Payload: payload,
		})
		if err != nil {
			return fmt.Errorf("error publishing: %s", err)
		} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 { // 16 = Server received message but there are no subscribers
			fmt.Printf("reason code %d received\n", pr.ReasonCode)
		}
	}
	return nil
}
