package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type SPort struct {
	DeviceInfo
	Value              types.CurrentState
	CurrentColourValue int
	IsActive           bool
	Current            float64
}

func (sport *SPort) SetValue(state types.CurrentState) {
	sport.IsActive = state == types.CurrentStateOn
	sport.Value = state
}

func NewSPort() *SPort {
	var probe SPort
	probe.Type = "SPort"
	probe.Units = "State"
	return &probe
}
