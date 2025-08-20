// consume messages from emqx using mqtt protocol.
package main

import (
	"context"
	"log/slog"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

const (
	CheckDelay = time.Second
)

func onMessage(pr paho.PublishReceived) (bool, error) {
	msg := pr.Packet

	slog.Default().Info("New message", "topic", msg.Topic, "payload", string(msg.Payload))

	return true, nil
}

// nolint: exhaustruct
func main() {
	logger := slog.Default()

	mqttURL, _ := url.Parse("mqtt://127.0.0.1:1883")

	conn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		ServerUrls:   []*url.URL{mqttURL},
		ClientConfig: paho.ClientConfig{
			ClientID: "mqtt101-consumer-client",
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				onMessage,
			},
		},
		OnConnectionUp: func(conn *autopaho.ConnectionManager, _ *paho.Connack) {
			ctx := context.Background()

			if _, err := conn.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: "test", QoS: 0},
				},
			}); err != nil {
				logger.Error("subscribing failed", "error", err)
			}
		},
	})
	if err != nil {
		logger.Error("failed to create a new connection", "error", err)

		return
	}

	if err := conn.AwaitConnection(context.Background()); err != nil {
		logger.Error("failed to wait for a new connection", "error", err)

		return
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	logger.Info("signal caught - exiting")

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = conn.Disconnect(ctx)

	logger.Info("shutdown complete")
}
