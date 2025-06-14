package mqtt_manager

import "strings"

type topicParser struct{}

func NewDefaultTopicParser() Parser {
	return &topicParser{}
}

func (p *topicParser) ParseUid(topic string) (string, error) {
	if topic == "" {
		return "", &ValidationError{Message: "topic cannot be empty"}
	}

	topicParts := strings.Split(topic, "/")
	if len(topicParts) != 3 {
		return "", &TopicFormatError{Topic: topic}
	}

	uid := topicParts[1]
	if uid == "" {
		return "", &ValidationError{Message: "uid cannot be empty"}
	}

	return uid, nil
}
