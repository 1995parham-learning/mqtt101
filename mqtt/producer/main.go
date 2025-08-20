// produce messages to emqx using mqtt protocol.
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
	PublishDelay         = 3 * time.Second
	ConnectTimeout       = 5 * time.Second
	ConnectRetryInterval = 3 * time.Second
	MaxReconnectInterval = 9 * time.Second
)

type Config struct {
	ClientID string `koanf:"client_id"`
	URL      string `koanf:"url"`
}

// Connect into mqtt broker.
// nolint: exhaustruct
func Connect(cfg Config, logger *slog.Logger) *autopaho.ConnectionManager {
	mqttURL, _ := url.Parse(cfg.URL)

	conn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		ServerUrls: []*url.URL{mqttURL},
		ClientConfig: paho.ClientConfig{
			ClientID: cfg.ClientID,
		},
	})
	if err != nil {
		logger.Error("failed to create a new connection", "error", err)

		return nil
	}

	logger.Info("wait for the new connection 🚧")

	if err := conn.AwaitConnection(context.Background()); err != nil {
		logger.Error("failed to wait for a new connection", "error", err)

		return nil
	}

	logger.Info("connection successful ✅")

	return conn
}

// nolint: exhaustruct
func main() {
	logger := slog.Default().With("role", "producer")

	cfg := Config{
		ClientID: "mqtt101-producer-client",
		URL:      "mqtt://127.0.0.1:1883",
	}

	conn := Connect(cfg, logger.With("component", "client"))

	go func() {
		for {
			_, err := conn.Publish(context.Background(), &paho.Publish{
				QoS:     0,
				Topic:   "test",
				Payload: []byte("hello world"),
			})
			if err != nil {
				slog.Error("failed to publish", "error", err)
			}

			time.Sleep(PublishDelay)
		}
	}()

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
