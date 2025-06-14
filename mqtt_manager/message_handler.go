package mqtt_manager

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
)

type MessageHandlerAdapter struct {
	handler func(mqtt.Client, Message)
	parser  Parser
}

func NewMessageHandlerAdapter(h func(mqtt.Client, Message), parser Parser) MessageHandler {
	if parser == nil {
		parser = NewDefaultTopicParser()
	}
	return &MessageHandlerAdapter{
		handler: h,
		parser:  parser,
	}
}

func (adapter *MessageHandlerAdapter) MessageHandle(client mqtt.Client, msg mqtt.Message) {
	customMsg, err := NewMessage(msg, adapter.parser)
	if err != nil {
		log.Error().Err(err).Msg("error parsing custom message")
		return
	}

	adapter.handler(client, customMsg)
}
