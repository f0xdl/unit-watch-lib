package domain

import "errors"

var (
	ErrObjectNotRegistered = errors.New("object does not registered")
	ErrValidatorDistracted = errors.New("validator distracted")
	ErrPingFailed          = errors.New("ping failed")
	ErrDeviceUnactivated   = errors.New("device unactivated")
	ErrUnknownDeviceStatus = errors.New("unknown device status")
	ErrDeviceRegistered    = errors.New("device registered")
)
