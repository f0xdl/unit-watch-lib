package domain

import "time"

type Device struct {
	UID          string       `json:"uid"`
	Label        string       `json:"label"`
	PasswordHash string       `json:"password_hash"`
	Active       bool         `json:"active"`
	Status       DeviceStatus `json:"status"`
	OwnerId      int64        `json:"owner_id"`
	Point        string       `json:"point"`
	ExpiresAt    time.Time    `json:"expires_at"`
	Online       bool         `json:"online"`
	LastSeen     time.Time    `json:"last_seen"`
}

type DeviceWithPassword struct {
	Device
	Password string
}

func GetDeviceWithPassword(d *Device, password string) *DeviceWithPassword {
	return &DeviceWithPassword{
		Device:   *d,
		Password: password,
	}
}
