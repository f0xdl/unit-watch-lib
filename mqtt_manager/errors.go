package mqtt_manager

type TopicFormatError struct {
	Topic string
}

func (e *TopicFormatError) Error() string {
	return "topic format error: " + e.Topic
}

type PayloadConversionError struct {
	Payload string
}

func (e *PayloadConversionError) Error() string {
	return "error converting status to int: " + e.Payload
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}
