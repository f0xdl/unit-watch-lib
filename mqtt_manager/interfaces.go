package mqtt_manager

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/f0xdl/unit-watch-lib/domain"
)

type MessageHandler interface {
	MessageHandle(client mqtt.Client, msg mqtt.Message)
}
type Message interface {
	mqtt.Message
	GetUid() (string, error)
	GetDeviceStatus() (domain.DeviceStatus, error)
	GetPayloadBool() (bool, error)
}

type Parser interface {
	ParseUid(topic string) (string, error)
}
