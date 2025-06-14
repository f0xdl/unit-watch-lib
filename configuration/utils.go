package configuration

func ValidateLoggerLevel(level string) bool {
	switch level {
	case
		"debug",
		"info",
		"warn",
		"error":
		return true
	default:
		return false
	}
}
