package main

import (
	"crypto/tls"
	"log"
	"net/url"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


type Config struct {
	ClientID string `koanf:"client_id"`
	URL      string `koanf:"url"`
}

func Connect(cfg Config) mqtt.Client {

	opts := mqtt.NewClientOptions().AddBroker(cfg.URL).SetClientID(cfg.ClientID).SetAutoReconnect(true).
		SetConnectRetry(true).SetConnectTimeout(5 * time.Second).SetConnectRetryInterval(3 * time.Second).SetMaxReconnectInterval(9 * time.Second)

	//opts.SetKeepAlive(60 * time.Second)   //nolint:gomnd
	//opts.SetPingTimeout(10 * time.Second) //nolint:gomnd

	opts.SetConnectionAttemptHandler(func(broker *url.URL, tlsCfg *tls.Config) *tls.Config {
		log.Printf("ConnectionAttemptHandler: broker: %s\n", broker)
		return tlsCfg
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("ConnectionLostHandler: err: %s\n", err)
	})

	opts.SetConnectionLostHandler(func(client mqtt.Client, err error) {
		log.Printf("ConnectionLostHandler: err: %s\n", err)
	})

	opts.SetReconnectingHandler(func(client mqtt.Client, options *mqtt.ClientOptions) {
		log.Printf("ReconnectingHandler\n")
	})

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Printf("OnConnectHandler\n")
	})


	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return c
}

func main() {

	cfg := Config{
		ClientID: "test-client",
		URL:      "tcp://localhost:1883",
	}

	client := Connect(cfg)

	println(client.IsConnected())

	for {
		if token := client.Publish("test", 0, false, "test"); token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		} else {
			log.Println("successful publish")
		}
		time.Sleep(3 * time.Second)
	}

}
