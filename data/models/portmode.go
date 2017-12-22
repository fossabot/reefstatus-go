package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type PortMode struct {
	DeviceMode types.DeviceMode
	Port       int
	BlackOut   int
	Invert     bool
	ID         string
}

func (mode PortMode) IsProbe() bool {
	return mode.DeviceMode == types.DeviceModeDecrease ||
		mode.DeviceMode == types.DeviceModeIncrease ||
		mode.DeviceMode == types.DeviceModeSubstrate ||
		mode.DeviceMode == types.DeviceModeProbeAlarm
}
