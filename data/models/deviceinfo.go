package models

import (
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
)

type DeviceInfo struct {
	BaseInfo
	PortNumber int
	Mode       profilux.PortMode
}

func (info DeviceInfo) IsConstant() bool {
	return info.Mode.DeviceMode == types.DeviceModeAlwaysOn || info.Mode.DeviceMode == types.DeviceModeAlwaysOff
}
