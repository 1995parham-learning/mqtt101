<h1 align="center">MQTT 101</h1>
<p align="center">
  <img src=".github/assets/banner.png"><br>
</p>

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

## Deployment

As you may guess, the hardest problem after writing the first MQTT application is deploying the Broker.
I've chosen EMQ here because of its performance.
At [AoLab](https://github.com/AoLab) Iman and I deployed it with default configuration on a small (4G, 4 cores) virtual machine,
and it handled 1000 concurrent connections without any issue.

Now, we want to deploy it on cloud, and we want to use its [official chart](https://github.com/emqx/emqx/tree/master/deploy/charts).
On the production, you must pay attention to security and performance.
For security, you can choose between different methods of authentication and authorization, which are provided by EMQ.
The most flexible option is using HTTP authentication because you can write whatever you want on your HTTP server.

```bash
helm repo add emqx https://repos.emqx.io/charts
helm dependency build
helm install emqx . -f values.yaml
```
