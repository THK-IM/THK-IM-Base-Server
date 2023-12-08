package mq

type (
	OnMessageReceived func(map[string]interface{}) error

	Subscriber interface {
		Sub(onReceived OnMessageReceived)
	}

	Publisher interface {
		Pub(id string, msg map[string]interface{}) error
	}
)
