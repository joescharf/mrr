package mrr

import (
// "errors"
// "github.com/davecgh/go-spew/spew"
// "strings"
// MQTT "github.com/eclipse/paho.MQTT.golang"
)

type (
	// Route contains a handler and information for matching against requests.
	Route struct {
		Topic   *Topic
		Handler HandlerFunc
		Err     error
	}
)

// func (r *Route) GetError() error {
// 	return r.err
// }

// func (r *Route) Handler(handler HandlerFunc) *Route {
// 	if r.err == nil {
// 		r.handler = handler
// 	}
// 	return r
// }

// func (r *Route) GetHandler() HandlerFunc {
// 	return r.handler
// }

// func (r *Route) Topic(name string, qos byte) *Route {
// 	if r.err == nil {
// 		r.topic = &Topic{
// 			name: name,
// 			qos:  qos,
// 		}
// 	}
// 	return r
// }

// func (r *Route) GetTopic() *Topic {
// 	return r.topic
// }
