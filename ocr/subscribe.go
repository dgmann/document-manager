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

type Subscriber struct {
	MQTTClient
}

type Subscription struct {
	Topic string
	paho.SubscribeOptions
	Handler paho.MessageHandler
}

func NewSubscriber(ctx context.Context, broker *url.URL, clientId string, subscriptions []Subscription) *Subscriber {
	cliCfg := autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{broker},
		KeepAlive:         30,
		ConnectRetryDelay: 10 * time.Second,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Println("Subscriber: mqtt connection up")
			if _, err := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: func() map[string]paho.SubscribeOptions {
					options := make(map[string]paho.SubscribeOptions)
					for _, sub := range subscriptions {
						options[sub.Topic] = sub.SubscribeOptions
					}
					return options
				}(),
			}); err != nil {
				log.Printf("Subscriber: failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			log.Println("Subscriber: mqtt subscription made")
		},
		OnConnectError: func(err error) { log.Printf("Subscriber: error whilst attempting connection: %s\n", err) },
		ClientConfig: paho.ClientConfig{
			ClientID: clientId,
			Router: func() paho.Router {
				router := paho.NewStandardRouter()
				for _, sub := range subscriptions {
					router.RegisterHandler(sub.Topic, sub.Handler)
				}
				return router
			}(),
			OnClientError: func(err error) { log.Printf("Subscriber: server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("Subscriber: server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("Subscriber: server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
	return &Subscriber{MQTTClient{Cfg: cliCfg}}
}

func (p *Subscriber) Run(ctx context.Context, topic string, ocrRequestChan <-chan OCRRequest) error {
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
