package mrr

type (
	Topic struct {
		name string
		qos  byte
	}
)

func NewTopic(name string, qos byte) *Topic {
	return &Topic{name: name, qos: qos}
}
func (t *Topic) Name() string {
	return t.name
}

func (t *Topic) SetName(name string) {
	t.name = name
}

func (t *Topic) Qos() byte {
	return t.qos
}

func (t *Topic) SetQos(qos byte) {
	t.qos = qos
}
