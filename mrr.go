package mrr

import (
	MQTT "github.com/eclipse/paho.MQTT.golang"

	"reflect"
	"runtime"

	"github.com/codegangsta/inject"
	"github.com/golang/glog"

	_ "github.com/davecgh/go-spew/spew"
)

type (
	// Mrr is the top-level framework instance
	Mrr struct {
		inject.Injector
		Client     MQTT.Client
		routes     []*Route
		debugTopic *Topic
	}

	Topic struct {
		Name string
		Qos  byte
	}

	Route struct {
		Topic   *Topic
		Handler Handler
		Err     error
	}

	// Handler can be any callable function.
	// Mrr will inject services into handler's argument list
	// and panics if an argument could not be fullfilled via dependency injection.
	Handler interface{}
)

func validateHandler(h Handler) Handler {

	// Test handler is a function:
	if reflect.TypeOf(h).Kind() != reflect.Func {
		panic("MRR handler must be a callable function")
	}

	return h
}

func New(c MQTT.Client) *Mrr {
	m := &Mrr{
		Injector: inject.New(),
		Client:   c,
	}
	return m
}

func NewTopic(name string, qos byte) *Topic {
	return &Topic{Name: name, Qos: qos}
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
func (m *Mrr) Add(topicName string, qos byte, h Handler) {

	// Subscribe to the topic and specify ServeMQTT as common handler:
	if token := m.Client.Subscribe(topicName, qos, m.routeMQTT); token.Wait() && token.Error() != nil {
		glog.Errorf("Error subscribing to topic: %s, Error: %s", topicName, token.Error())
	}
	glog.Infoln("Subscribed to Topic:", topicName)

	// Validate handler:
	h = validateHandler(h)

	//  Define the new route
	route := &Route{Topic: NewTopic(topicName, qos), Handler: h}

	// Add topic and handler to routing table:
	m.routes = append(m.routes, route)
}

func (m *Mrr) SetDebugTopic(n string) {
	m.debugTopic = NewTopic(n, 0)
}

// routeMQTT is common router for all subscriptions
func (m *Mrr) routeMQTT(c MQTT.Client, msg MQTT.Message) {

	// Create a request based on the msg
	request := &Request{
		Topic:   NewTopic(msg.Topic(), msg.Qos()),
		Payload: msg.Payload(),
	}

	// Get a new context (conversation) and inject it as a service
	conversation := m.newConversation(request)
	m.Map(conversation)

	// call the handler for the topic
	if route := m.findRoute(request.Topic); route != nil {
		handler := route.Handler
		_, err := m.Invoke(handler) // Dep injection
		if err != nil {
			glog.Errorln("Error Invoking() handler: ", err)
		}
	}
}

func (m *Mrr) newConversation(r *Request) ConversationInterface {
	c := &Conversation{
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
	name := topic.Name
	for _, route := range m.routes {
		if route.Topic.Name == name {
			return route
		}
	}
	return nil
}

func handlerName(h Handler) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}
