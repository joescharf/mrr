package mrr

import (
	MQTT "github.com/eclipse/paho.MQTT.golang"

	"github.com/golang/glog"

	"github.com/davecgh/go-spew/spew"
)

type DeviceCmd struct {
	Event       string                 `json:"event"`
	DeviceIndex int                    `json:"device_index"`
	Payload     map[string]interface{} `json:"payload"`
}

func HandleConnect(c MQTT.Client) {
	glog.Infoln("MQTT Client Connected")
}

func HandleConnectionLost(c MQTT.Client, err error) {
	glog.Infoln("MQTT Connection Lost:", err)
}

func HandleMessage(c MQTT.Client, m MQTT.Message) {
	glog.Infoln("Received message on topic:", m.Topic())
	spew.Dump(m.Payload())
}
