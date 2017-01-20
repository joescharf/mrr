package mrr

import (
	MQTT "github.com/eclipse/paho.MQTT.golang"
)

type (
	ClientOptions struct {
		Url                   string
		ClientId              string
		Username              string
		Password              string
		ConnectionHandler     MQTT.OnConnectHandler
		ConnectionLostHandler MQTT.ConnectionLostHandler
	}
)

func NewClientOptions() *ClientOptions {
	// Return clientoptions with default handlers set
	return &ClientOptions{
		ConnectionHandler:     HandleConnect,
		ConnectionLostHandler: HandleConnectionLost,
	}
}
