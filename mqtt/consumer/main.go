package main

import (
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	CheckDelay = time.Second
)

func onMessage(_ mqtt.Client, msg mqtt.Message) {
	log.Printf("New message on %s: %s", msg.Topic(), string(msg.Payload()))
}

func main() {
	opts := mqtt.NewClientOptions().
		AddBroker("127.0.0.1").
		SetClientID("mqtt101-consumer-client")

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		if token := client.Subscribe("test", 0, onMessage); token.Wait() && token.Error() != nil {
			log.Printf("subscribing failed %s\n", token.Error())
		}
	})

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for {
		if !client.IsConnectionOpen() {
			log.Println("client is not connected, wait for connection")
		} else {
			time.Sleep(CheckDelay)
		}
	}
}
