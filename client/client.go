package client

import (
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	C MQTT.Client
}

func New(server string, id string) Client {
	f := func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("CLIENT ID: %s\n", id)
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}

	return NewWithMessageHandler(server, id, f)
}

func NewWithMessageHandler(server string, id string, f MQTT.MessageHandler) Client {
	opts := MQTT.NewClientOptions().AddBroker(server)
	opts.SetClientID(id)
	opts.SetDefaultPublishHandler(f)

	cli := Client{C: MQTT.NewClient(opts)}

	return cli
}

func (cli *Client) Connect(topics ...string) {
	if token := cli.C.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for _, topic := range topics {
		cli.Subscribe(topic)
	}
}

func (cli *Client) Subscribe(topic string) {
	if token := cli.C.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func (cli *Client) Publish(topic string, text string) {
	token := cli.C.Publish(topic, 0, false, text)
	token.Wait()
}
