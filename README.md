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
    "github.com/joescharf/mrr"
)
func main () {
    // Set MQTT connection optioins:
    options := mrr.NewClientOptions()
    options.Url = viper.GetString("MQTT.url")
    options.ClientId = viper.GetString("MQTT.client_id")
    options.Username = viper.GetString("MQTT.username")
    options.Password = viper.GetString("MQTT.password")

    // Create client with options
    m := mrr.NewClient(options)

    // Connect to server
    m.Connect()

    // Add some routes & subscribe to topics
    m.Add("device/cmd/helloworld", 0, HandleHelloWorld)

}

func HandleHelloWorld(c mrr.Conversation) error {
    c.String(200, "Hello World!!")
    return nil
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
