package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type DeviceInfo struct {
	BaseInfo
	PortNumber int
	Mode       PortMode
}

func (info DeviceInfo) IsConstant() bool {
	return info.Mode.DeviceMode == types.DeviceModeAlwaysOn || info.Mode.DeviceMode == types.DeviceModeAlwaysOff
}
