# Overview

MQTT Request Response (MRR) is a library that facilitates easy request-response handling over a MQTT connection. 

# Installing

`go get -u github.com/joescharf/mrr`

# Usage

Import the library into your go code with:
`import "github.com/joescharf/mrr"`


This library currently uses the [Paho MQTT Client](https://github.com/eclipse/paho.mqtt.golang) to connect to MQTT servers


```
package main
import (
    MQTT "github.com/eclipse/paho.MQTT.golang"
    "github.com/joescharf/mrr"
)

func main () {

    // example device structure for custom service injection:
    dev := &devInfo{}

    // Set MQTT connection options (paho mqtt client):
    options := MQTT.NewClientOptions()
    options.AddBroker(viper.GetString("MQTT.url"))
    options.ClientID = viper.GetString("MQTT.client_id")
    options.Username = viper.GetString("MQTT.username")
    options.Password = viper.GetString("MQTT.password")

    // Create MQTT client based on options:
    mqtt := MQTT.NewClient(options)

    // Establish MRR with mqtt client:
    mrr := mrr.New(mqtt)

    // Inject the devInfo{} service for handlers:
    mrr.Map(dev)

    // Connect to MQTT Server
    mrr.Connect()
    api.SetupRoutes(mrr)

    // Add some routes & subscribe to topics
    m.Add("device/cmd/helloworld", 0, HandleHelloWorld)
    m.Add("device/cmd/open", 0, HandleDeviceOpen)

}

// Simple handler that just takes a mrr.Conversation
func HandleHelloWorld(c mrr.Conversation) error {
    c.String(200, "Hello World!!")
    return nil
}

// Another handler that also depends on the devInfo struct,
// And returns some data on a JSON payload
func HandleDeviceOpen(c *mrr.Conversation, d *devInfo) error {
    settings := d.Open()

    payload := make(map[string]interface{})
    payload["id"] = 42
    payload["settings"] = settings

    c.JSON(200, payload)
}

```

# Request

MRR expects a JSON encoded payload and will use the following fields to determine the response topic and QOS: 

* `_rt` - Response Topic - (string)
* `_rq` - Response QOS - (int)

## Example:

```
{"payload":{"option1":123, "option2":"test"}, "_rt":"response/topic", "_rq":1}
```


## Default Response Topic

If the response topic fields are omitted, and a response is issued using the `Conversation` response methods (i.e. `Conversation.String()`) MRR will by default use the `request topic + "_response"` and will match the request's QOS:

Request topic: `/device/cmd/helloworld`
Response topic: `/device/cmd/helloworld/_response`

# Handlers

Mrr uses Dependency Injection to resolve dependencies in a handler's argument list modeled after [Macaron Custom Services](https://go-macaron.com/docs/advanced/custom_services) 

map a service as follows:

```
dev := &devInfo{}
mrr := mrr.New(mqttClient)
mrr.Map(dev)

qos := 0
mrr.Add("topic/goes/here", qos, SomeHandler))

func SomeHandler(ctx *macaron.Context, d *devInfo) {
    // Do stuff
}
```

`Mrr` injects the `*mrr.Conversation` service and is most commonly used service in your handlers:

```
func SomeHandler(ctx *macaron.Context) {
```


