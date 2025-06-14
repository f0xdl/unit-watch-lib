package domain

import (
	"strings"
)

type DeviceStatus int

const (
	/* Default statuses*/

	UnknownStatus      DeviceStatus = 0
	NormalStatus       DeviceStatus = 1
	DefaultErrorStatus DeviceStatus = 2

	/* BV (Bill Validator) statuses */

	IdlingStatus           DeviceStatus = 0x14
	DropCassetteFullStatus DeviceStatus = 0x41
	CassetteRemovedStatus  DeviceStatus = 0x42
	ValidatorJammedStatus  DeviceStatus = 0x43
	CassetteJammedStatus   DeviceStatus = 0x44
	CheatedStatus          DeviceStatus = 0x45
	PauseStatus            DeviceStatus = 0x46
	FailureStatus          DeviceStatus = 0x47
)

func (s DeviceStatus) String() string {
	switch s {
	case NormalStatus:
		return "NORMAL"
	case DefaultErrorStatus:
		return "BASE-ERROR"
	case IdlingStatus:
		return "IDLING"
	case DropCassetteFullStatus:
		return "DROP-CASSETTE-FULL"
	case CassetteRemovedStatus:
		return "CASSETTE-REMOVED"
	case ValidatorJammedStatus:
		return "VALIDATOR-JAMMED"
	case CassetteJammedStatus:
		return "CASSETTE-JAMMED"
	case CheatedStatus:
		return "CHEATED"
	case PauseStatus:
		return "PAUSE"
	case FailureStatus:
		return "FAILURE"
	default:
		return "UNKNOWN"
	}
}

func DeviceStatusFromInt(status int) DeviceStatus {
	return DeviceStatus(status)
}

func ParseString(value string) DeviceStatus {
	switch strings.ToUpper(value) {
	case "NORMAL":
		return NormalStatus
	case "BASE-ERROR":
		return DefaultErrorStatus
	case "IDLING":
		return IdlingStatus
	case "DROP-CASSETTE-FULL":
		return DropCassetteFullStatus
	case "CASSETTE-REMOVED":
		return CassetteRemovedStatus
	case "VALIDATOR-JAMMED":
		return ValidatorJammedStatus
	case "CASSETTE-JAMMED":
		return CassetteJammedStatus
	case "CHEATED":
		return CheatedStatus
	case "PAUSE":
		return PauseStatus
	case "FAILURE":
		return FailureStatus
	default:
		return UnknownStatus
	}
}
