package mrr

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"
	"strconv"

	"github.com/golang/glog"
)

type (
	ConversationInterface interface {
		// Request returns `*Request`.
		Request() *Request

		// SetRequest sets `*Request`.
		SetRequest(r *Request)

		// Payload
		Payload() map[string]interface{}

		// Set payload, straight from mqtt msg
		SetPayload(data []byte)

		// Set response topic for reply
		SetResponseTopic(name string, qos byte)

		// Param returns payload parameter by name
		Param(name string) interface{}

		// ParamBool returns payload parameter as bool
		ParamBool(name string) bool

		// ParamByte returns parameter as byte via float64
		ParamByte(name string) byte

		// ParamInt64 returns parameter as int64, via float64
		ParamInt64(name string) int64

		// ParamFloat64 returns payload parameter as string
		ParamFloat64(name string) float64

		// ParamString returns payload parameter as string
		ParamString(name string) string

		// String sends a string response with status code.
		String(code int, s string) error

		// JSON sends a JSON response with status code.
		JSON(code int, i interface{}) error

		// JSONPretty sends a pretty-print JSON with status code.
		JSONPretty(code int, i interface{}, indent string) error

		// JSONBlob sends a JSON blob response with status code.
		JSONBlob(code int, b []byte) error

		// Blob sends a blob response with status code and content type.
		Blob(code int, b []byte) error

		// // NoContent sends a response with no body and a status code.
		// NoContent(code int) error

		// Mrr returns the `Mrr` instance.
		Mrr() *Mrr
	}

	// Conversation is the implementation of ConversationInterface:
	Conversation struct {
		request *Request
		payload map[string]interface{}
		mrr     *Mrr
	}
)

func (c *Conversation) Request() *Request {
	return c.request
}

func (c *Conversation) SetRequest(r *Request) {
	c.request = r
}

func (c *Conversation) Payload() map[string]interface{} {
	return c.payload
}

func (c *Conversation) SetResponseTopic(name string, qos byte) {
	c.request.ResponseTopic = NewTopic(name, qos)
}

func (c *Conversation) SetPayload(data []byte) {
	json.Unmarshal(data, &c.payload)
	glog.Infoln("Payload: ", c.payload)
}

func (c *Conversation) Param(name string) interface{} {
	return c.payload[name]
}

func (c *Conversation) ParamBool(name string) bool {
	return c.payload[name].(bool)
}

func (c *Conversation) ParamByte(name string) byte {
	val := c.ParamFloat64(name)
	bits := math.Float64bits(val)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes[0]
}

func (c *Conversation) ParamInt64(name string) int64 {
	val := c.ParamFloat64(name)
	return int64(val)
}

func (c *Conversation) ParamFloat64(name string) float64 {
	ret := 0.0
	if val, ok := c.payload[name]; ok {
		switch t := reflect.TypeOf(val).Name(); t {
		case "string":
			ret, _ = strconv.ParseFloat(val.(string), 64)
		case "float64":
			ret = val.(float64)
		default:
			ret = 0.0
		}
	}
	return ret
}

func (c *Conversation) ParamString(name string) string {
	ret := ""
	if val, ok := c.payload[name]; ok {

		switch t := reflect.TypeOf(val).Name(); t {
		case "string":
			ret = val.(string)
		case "float64":
			ret = strconv.FormatFloat(val.(float64), 'g', -1, 64)
		default:
			ret = ""
		}
	}
	return ret
}

func (c *Conversation) String(code int, s string) (err error) {
	return c.Blob(code, []byte(s))
}

func (c *Conversation) JSON(code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return
	}
	return c.JSONBlob(code, b)
}

func (c *Conversation) JSONPretty(code int, i interface{}, indent string) (err error) {
	b, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		return
	}
	return c.JSONBlob(code, b)
}

func (c *Conversation) JSONBlob(code int, b []byte) (err error) {
	return c.Blob(code, b)
}

// Blob writes the response to the designated responseTopic. If
// responseTopic is nil, we use default of:
// <incoming_topic>/_response with Qos matching the request Qos.
func (c *Conversation) Blob(code int, b []byte) (err error) {
	// Determine what response topic is
	rt := c.Request().ResponseTopic
	if rt.Name != "" {
		c.Mrr().Client.Publish(rt.Name, rt.Qos, false, b)
	} else {
		c.Mrr().Client.Publish(c.Request().Topic.Name+"/_response", c.Request().Topic.Qos, false, b)
	}

	// Check to see if we need to echo the response to the debug topic:
	dt := c.Mrr().debugTopic
	if dt.Name != "" {
		c.Mrr().Client.Publish(dt.Name, dt.Qos, false, b)
	}
	return
}

func (c *Conversation) Mrr() *Mrr {
	return c.mrr
}
