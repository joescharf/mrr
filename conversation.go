package mrr

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"reflect"
)

type (
	Conversation interface {
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

	conversation struct {
		request *Request
		payload map[string]interface{}
		mrr     *Mrr
	}
)

func (c *conversation) Request() *Request {
	return c.request
}

func (c *conversation) SetRequest(r *Request) {
	c.request = r
}

func (c *conversation) Payload() map[string]interface{} {
	return c.payload
}

func (c *conversation) SetResponseTopic(name string, qos byte) {
	c.request.ResponseTopic = NewTopic(name, qos)
}

func (c *conversation) SetPayload(data []byte) {
	json.Unmarshal(data, &c.payload)
}

func (c *conversation) Param(name string) interface{} {
	return c.payload[name]
}

func (c *conversation) ParamBool(name string) bool {
	return c.payload[name].(bool)
}

func (c *conversation) ParamByte(name string) byte {
	val := c.ParamFloat64(name)
	bits := math.Float64bits(val)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes[0]
}

func (c *conversation) ParamInt64(name string) int64 {
	val := c.ParamFloat64(name)
	return int64(val)
}

func (c *conversation) ParamFloat64(name string) float64 {
	if val, ok := c.payload[name]; ok {

		if reflect.TypeOf(val).Name() == "float64" {
			return c.payload[name].(float64)
		}
	}
	return 0.0
}

func (c *conversation) ParamString(name string) string {
	if val, ok := c.payload[name]; ok {
		if reflect.TypeOf(val).Name() == "string" {
			return val.(string)
		}
	}
	return ""
}

func (c *conversation) String(code int, s string) (err error) {
	return c.Blob(code, []byte(s))
}

func (c *conversation) JSON(code int, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return
	}
	return c.JSONBlob(code, b)
}

func (c *conversation) JSONPretty(code int, i interface{}, indent string) (err error) {
	b, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		return
	}
	return c.JSONBlob(code, b)
}

func (c *conversation) JSONBlob(code int, b []byte) (err error) {
	return c.Blob(code, b)
}

// Blob writes the response to the designated responseTopic. If
// responseTopic is nil, we use default of:
// <incoming_topic>/_response with Qos matching the request Qos.
func (c *conversation) Blob(code int, b []byte) (err error) {
	// Determine what response topic is
	rt := c.Request().ResponseTopic
	if rt.Name() != "" {
		c.Mrr().Client.Publish(rt.Name(), rt.Qos(), false, b)
	} else {
		c.Mrr().Client.Publish(c.Request().Topic.Name()+"/_response", c.Request().Topic.Qos(), false, b)
	}
	return
}

func (c *conversation) Mrr() *Mrr {
	return c.mrr
}
