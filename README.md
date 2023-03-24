# MQTT 101

## Introduction

Technologies like CoAP or MQTT are known to be useful in the field of IoT, but in the today world
you can use them for any type of communication with your clients, so they are really important and useful.

MQTT protocol is supported on wide variety of devices from Android Smartphones to Linux servers that are placed in
datacenters. MQTT protocol has two sides, consumer and the producer and both sides are simple and easy to implement.
The hardest part is implementing the broker that handles the connection between producers and consumers.
I've worked with [EMQX](https://vernemq.com/) and [VerneMQ](https://vernemq.com/) as brokers and both of them are awesome.

In this repository we implement consumer and producer for MQTT protocol with Golang based on Paho library from Eclipse.

## Implementation

The pure Go implementation of mqtt protocol is available from [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang).
The following code grabbed from the library and shows that the given on-connect handler runs on another go routine.

```go
c.setConnected(connected)
DEBUG.Println(CLI, "client is connected/reconnected")
if c.options.OnConnect != nil {
 go c.options.OnConnect(c)
}
```

This implementation works great, and you must not forget to re-subscribe on connection lost
To this just do subscribe with on-connect.
