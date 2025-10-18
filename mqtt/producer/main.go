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
)

type Config struct {
	ClientID string `koanf:"client_id"`
	URL      string `koanf:"url"`
}

// Connect into mqtt broker.
// nolint: exhaustruct
func Connect(cfg Config, logger *slog.Logger) *autopaho.ConnectionManager {
	mqttURL, err := url.Parse(cfg.URL)
	if err != nil {
		logger.Error("failed to parse MQTT URL", "error", err)

		return nil
	}

	conn, err := autopaho.NewConnection(context.Background(), autopaho.ClientConfig{
		ServerUrls:                 []*url.URL{mqttURL},
		KeepAlive:                  30,
		ConnectRetryDelay:          ConnectRetryInterval,
		OnConnectionUp:             func(cm *autopaho.ConnectionManager, _ *paho.Connack) { logger.Info("mqtt connection up") },
		OnConnectError:             func(err error) { logger.Error("error whilst attempting connection", "error", err) },
		ClientConfig: paho.ClientConfig{
			ClientID: cfg.ClientID,
		},
		ConnectTimeout: ConnectTimeout,
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
	if conn == nil {
		logger.Error("failed to establish connection")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(PublishDelay)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if _, err := conn.Publish(ctx, &paho.Publish{
					QoS:     0,
					Topic:   "test",
					Payload: []byte("hello world"),
				}); err != nil {
					logger.Error("failed to publish", "error", err)
				}
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	logger.Info("signal caught - exiting")
	cancel()

	// We could cancel the context at this point but will call Disconnect instead (this waits for autopaho to shutdown)
	disconnectCtx, disconnectCancel := context.WithTimeout(context.Background(), time.Second)
	defer disconnectCancel()

	if err := conn.Disconnect(disconnectCtx); err != nil {
		logger.Error("error during disconnect", "error", err)
	}

	logger.Info("shutdown complete")
}
