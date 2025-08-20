// produce messages to emqx using mqtt protocol.
package main

import (
	"crypto/tls"
	"log/slog"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

func Connect(cfg Config, logger *slog.Logger) mqtt.Client {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.URL).
		SetClientID(cfg.ClientID).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectTimeout(ConnectTimeout).
		SetConnectRetryInterval(ConnectRetryInterval).
		SetMaxReconnectInterval(MaxReconnectInterval)

	// opts.SetKeepAlive(60 * time.Second)   //nolint:gomnd
	// opts.SetPingTimeout(10 * time.Second) //nolint:optionsgomnd

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		logger.Info("ConnectionAttemptHandler", "broker", broker)

		return tlsCfg
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logger.Info("ConnectionLostHandler", "error", err)
	})

	opts.SetReconnectingHandler(func(_ mqtt.Client, _ *mqtt.ClientOptions) {
		logger.Info("ReconnectingHandler")
	})

	opts.SetOnConnectHandler(func(_ mqtt.Client) {
		logger.Info("OnConnectHandler")
	})

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		logger.Error("mqtt connection failed", "error", token.Error())
	}

	return c
}

func main() {
	logger := slog.Default()

	cfg := Config{
		ClientID: "mqtt101-producer-client",
		URL:      "mqtt://127.0.0.1:1883",
	}

	client := Connect(cfg, logger.With("component", "client"))

	logger.Info("client connected", "is-connected", client.IsConnected())

	for {
		if !client.IsConnectionOpen() {
			logger.Info("client is not connected, wait for connection")
		} else {
			token := client.Publish("test", 0, false, "test")

			<-token.Done()

			if token.Error() != nil {
				logger.Error("", "error", token.Error())
			} else {
				logger.Info("successful publish")
			}
		}

		time.Sleep(PublishDelay)
	}
}
