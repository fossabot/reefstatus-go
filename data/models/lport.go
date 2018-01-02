package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type LPort struct {
	DeviceInfo
	Value float64
}

func NewLPort(index int) *LPort {
	var lPort LPort
	lPort.Type = "LPort"
	lPort.Units = "%"
	lPort.PortNumber = index
	lPort.Id = fmt.Sprintf("L%d", 1+index)
	return &lPort
}

func (port *LPort) Update(controller *profilux.Controller) {
	port.Mode = controller.GetLPortFunction(port.PortNumber)
	port.Value = controller.GetLPortValue(port.PortNumber)
}

func (port *LPort) UpdateState(controller *profilux.Controller) {
	port.Value = controller.GetLPortValue(port.PortNumber)
}
