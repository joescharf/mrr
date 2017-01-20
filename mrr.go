package mrr

import (
	MQTT "github.com/eclipse/paho.MQTT.golang"

	_ "encoding/json"
	_ "github.com/davecgh/go-spew/spew"
	"reflect"
	"runtime"

	"github.com/golang/glog"
)

type (
	// Mrr is the top-level framework instance
	Mrr struct {
		Client  MQTT.Client
		options *ClientOptions
		routes  []*Route
	}

	// HandlerFunc defines a function to serve MQTT requests.
	HandlerFunc func(Conversation) error
)

func NewClient(o *ClientOptions) *Mrr {
	m := &Mrr{
		options: o,
	}

	// Set configuration
	co := MQTT.NewClientOptions()
	co.AddBroker(o.Url)
	co.SetClientID(o.ClientId)
	co.SetUsername(o.Username)
	co.SetPassword(o.Password)
	co.SetOnConnectHandler(o.ConnectionHandler)
	co.SetConnectionLostHandler(o.ConnectionLostHandler)

	// Create the MQTT client:
	m.Client = MQTT.NewClient(co)
	return m
}

// Connect connects to mqtt server
func (m *Mrr) Connect() {
	glog.Infoln("MQTT Connecting")
	if token := m.Client.Connect(); token.Wait() && token.Error() != nil {
		glog.Errorln("MQTT Connect Error", token.Error())
	}
	glog.Infoln("MQTT Connected")
}

// Add subscribes to a topic and adds the topic and handler to routing table
func (m *Mrr) Add(topicName string, qos byte, h HandlerFunc) {

	// Subscribe to the topic and specify ServeMQTT as common handler:
	if token := m.Client.Subscribe(topicName, qos, m.routeMQTT); token.Wait() && token.Error() != nil {
		glog.Errorf("Error subscribing to topic: %s, Error: %s", topicName, token.Error())
	}
	glog.Infoln("Subscribed to Topic:", topicName)

	//  Define the new route
	route := &Route{Topic: NewTopic(topicName, qos), Handler: h}

	// Add topic and handler to routing table:
	m.routes = append(m.routes, route)
}

// RouteMQTT is common router for all subscriptions
func (m *Mrr) routeMQTT(c MQTT.Client, msg MQTT.Message) {

	// Create a request based on the msg
	request := &Request{
		Topic:   NewTopic(msg.Topic(), msg.Qos()),
		Payload: msg.Payload(),
	}

	// Get a new context (conversation)
	conversation := m.newConversation(request)

	// call the handler for the topic
	if route := m.findRoute(request.Topic); route != nil {
		handler := route.Handler
		handler(conversation)
	}
}

func (m *Mrr) newConversation(r *Request) Conversation {
	c := &conversation{
		request: r,
		mrr:     m,
	}

	// Convert the message payload:
	c.SetPayload(r.Payload)

	// Set response topic (_rt) and qos (_rq):
	c.SetResponseTopic(c.ParamString("_rt"), c.ParamByte("_rq"))

	return c
}

func (m *Mrr) findRoute(topic *Topic) *Route {
	name := topic.Name()
	for _, route := range m.routes {
		if route.Topic.Name() == name {
			return route
		}
	}
	return nil
}

func handlerName(h HandlerFunc) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}
