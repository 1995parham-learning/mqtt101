package client

import (
	"errors"
	"fmt"
	"os"
	"sync"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	c         MQTT.Client
	opts      *MQTT.ClientOptions
	topics    []string
	mutex     sync.Mutex
	connected bool
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

	cli := Client{
		opts:      opts,
		connected: false,
	}

	return cli
}

func (cli *Client) Connect() {
	cli.connected = true

	cli.opts.SetOnConnectHandler(func(client MQTT.Client) {
		for _, topic := range cli.topics {
			if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
				fmt.Println(token.Error())
				os.Exit(1)
			}
		}
	})

	cli.c = MQTT.NewClient(cli.opts)

	if token := cli.c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (cli *Client) Publish(topic string, text string) {
	token := cli.c.Publish(topic, 0, false, text)
	token.Wait()
}

func (cli *Client) Register(topics []string) error {
	if cli.connected {
		return errors.New("you can not register a new topic after connecting the client to the broker")
	}

	cli.mutex.Lock()
	cli.topics = append(cli.topics, topics...)
	cli.mutex.Unlock()

	return nil
}
