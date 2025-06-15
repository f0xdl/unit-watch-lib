package mqtt_manager

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog/log"
	"time"
)

type MqttManager struct {
	client      mqtt.Client
	topics      map[string]MessageHandler
	isConnected bool
}

func (m *MqttManager) SetTopicHandle(topic string, h MessageHandler) {
	if h != nil {
		m.topics[topic] = h
	} else {
		delete(m.topics, topic)
	}
}

func NewMqttManager(brokerURL, clientID, username, password string) *MqttManager {
	manager := &MqttManager{topics: make(map[string]MessageHandler)}
	log.Info().Msg("building mqtt client with auto-reconnect")
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetPassword(password)
	opts.SetUsername(username)

	// reconnect
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(30 * time.Second)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetConnectRetry(true)

	// keep alive
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(10 * time.Second)
	opts.SetConnectTimeout(10 * time.Second)
	opts.SetCleanSession(false)

	// Last Will and Testament
	opts.SetWill("system/"+clientID+"/status", "offline", 1, true)

	// Callback
	opts.SetOnConnectHandler(manager.onConnect)
	opts.SetConnectionLostHandler(manager.onConnectionLost)
	opts.SetReconnectingHandler(manager.onReconnecting)
	manager.client = mqtt.NewClient(opts)
	return manager
}

func (m *MqttManager) Connect() error {
	log.Info().Msg("connecting to MQTT broker")
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	log.Info().Msg("successfully connected to MQTT broker")
	m.isConnected = true
	return nil
}

func (m *MqttManager) Disconnect(timeout uint) {
	if m.client.IsConnected() {
		m.client.Publish("system/status", 1, true, "offline")
		m.client.Disconnect(timeout)
	}
	m.isConnected = false
}

func (m *MqttManager) IsConnected() bool {
	return m.isConnected && m.client.IsConnected()
}

func (m *MqttManager) Publish(topic string, payload interface{}) error {
	if !m.IsConnected() {
		return mqtt.ErrNotConnected
	}

	token := m.client.Publish(topic, 1, false, payload)
	token.Wait()
	return token.Error()
}

func subscribeTopics(client mqtt.Client, topics map[string]MessageHandler) {
	for topic, adapter := range topics {
		token := client.Subscribe(topic, 1, adapter.MessageHandle)
		if token.Wait() && token.Error() != nil {
			log.Error().
				Err(token.Error()).
				Str("topic", topic).
				Msg("error subscribing to topic")
		} else {
			log.Info().
				Str("topic", topic).
				Msg("subscribed to topic")
		}
	}
}

func (m *MqttManager) onReconnecting(_ mqtt.Client, _ *mqtt.ClientOptions) {
	log.Warn().Msg("attempting to reconnect to MQTT broker")
}

func (m *MqttManager) onConnectionLost(_ mqtt.Client, err error) {
	log.Error().
		Err(err).
		Msg("connection to MQTT broker lost")
}

func (m *MqttManager) onConnect(client mqtt.Client) {
	log.Info().Msg("connected to MQTT broker")
	subscribeTopics(client, m.topics)
}
