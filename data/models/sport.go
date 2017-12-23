package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
)

type SPort struct {
	DeviceInfo
	Value              types.CurrentState
	CurrentColourValue int
	IsActive           bool
	Current            float64
}

func NewSPort(index int) *SPort {
	var probe SPort
	probe.Type = "SPort"
	probe.Units = "State"
	probe.PortNumber = index
	probe.Id = fmt.Sprintf("S%d", 1+index)
	return &probe
}

func (port *SPort) Update(controller *profilux.Controller) {
	port.Mode = controller.GetSPortFunction(port.PortNumber)
	port.Value = controller.GetSPortValue(port.PortNumber)
	port.IsActive = port.Value == types.CurrentStateOn

	port.DisplayName = controller.GetSPortName(port.PortNumber)
}
