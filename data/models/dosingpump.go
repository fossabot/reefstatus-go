package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type TimerSettings struct {
	FeedPauseIfActive bool
	Mode              types.TimerMode
	DayMode           types.DayMode
	SwitchingCount    int
}

type DosingPump struct {
	BaseInfo
	Channel  int
	Rate     int
	PerDay   int
	Settings TimerSettings
}

func NewDosingPump() *DosingPump {
	var pump DosingPump
	pump.Type = "Dosing"
	pump.Units = "ml/day"
	return &pump
}
