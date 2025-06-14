package mqtt_manager

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/f0xdl/unit-watch-lib/domain"
	"strconv"
)

type message struct {
	mqtt.Message
	parser Parser
}

func NewMessage(msg mqtt.Message, parser Parser) (Message, error) {
	if msg == nil {
		return nil, &ValidationError{Message: "message cannot be nil"}
	}
	if parser == nil {
		return nil, &ValidationError{Message: "parser cannot be nil"}
	}

	return &message{
		Message: msg,
		parser:  parser,
	}, nil
}

func (msg *message) GetUid() (string, error) {
	return msg.parser.ParseUid(msg.Topic())
}

func (msg *message) GetDeviceStatus() (domain.DeviceStatus, error) {
	payload := msg.Payload()
	if len(payload) == 0 {
		return domain.UnknownStatus, &ValidationError{Message: "payload cannot be empty"}
	}

	payloadStr := string(payload)
	sInt, err := strconv.Atoi(payloadStr)
	if err != nil {
		return domain.UnknownStatus, &PayloadConversionError{Payload: payloadStr}
	}
	return domain.DeviceStatus(sInt), nil
}

func (msg *message) GetPayloadBool() (bool, error) {
	return strconv.ParseBool(string(msg.Payload()))
}
