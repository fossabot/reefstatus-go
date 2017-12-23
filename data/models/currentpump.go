package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type CurrentPump struct {
	BaseInfo
	Value int
	Index int
}

func NewCurrentPump(index int) *CurrentPump {
	var pump CurrentPump
	pump.Type = "CurrentPump"
	pump.Units = "%"
	pump.Index = index
	pump.Id = fmt.Sprintf("CurrentPumpm%d", 1+index)

	return &pump
}

func (pump *CurrentPump) Update(controller *profilux.Controller) {
	pump.Value = controller.GetCurrentPumpValue(pump.Index)
}
