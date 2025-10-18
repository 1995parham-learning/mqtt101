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

// nolint: exhaustruct, funlen
func main() {
	logger := slog.Default().With("role", "consumer")

	onMessage := func(pr paho.PublishReceived) (bool, error) {
		msg := pr.Packet
		logger.Info("New message", "topic", msg.Topic, "payload", string(msg.Payload))

		return true, nil
	}

	mqttURL, err := url.Parse("mqtt://127.0.0.1:1883")
	if err != nil {
		logger.Error("failed to parse MQTT URL", "error", err)

		return
	}

	conn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		ServerUrls: []*url.URL{mqttURL},
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

	logger.Info("wait for the new connection 🚧")

	if err := conn.AwaitConnection(context.Background()); err != nil {
		logger.Error("failed to wait for a new connection", "error", err)

		return
	}

	logger.Info("connection successful ✅")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	logger.Info("signal caught - exiting")

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := conn.Disconnect(ctx); err != nil {
		logger.Error("error during disconnect", "error", err)
	}

	logger.Info("shutdown complete")
}
