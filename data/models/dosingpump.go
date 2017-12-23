package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type DosingPump struct {
	BaseInfo
	Channel  int
	Rate     int
	PerDay   int
	Settings profilux.TimerSettings
}

func NewDosingPump(index int) *DosingPump {
	var pump DosingPump
	pump.Channel = index
	pump.Type = "Dosing"
	pump.Units = "ml/day"
	pump.Id = fmt.Sprintf("Pump%d", 1+index)
	return &pump
}

func (pump *DosingPump) Update(controller *profilux.Controller) {
	pump.Settings = controller.GetTimerSettings(pump.Channel)
	pump.Rate = controller.GetDosingRate(pump.Channel)
	pump.PerDay = pump.Settings.SwitchingCount
}
