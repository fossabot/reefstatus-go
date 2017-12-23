package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type SensorInfo struct {
	BaseInfo
	Format     int
	SensorType types.SensorType
	Index      int
	AlarmState types.CurrentState
}
