package models

type CurrentPump struct {
	BaseInfo
	Value int
	Index int
}

func NewCurrentPump() *CurrentPump {
	var pump CurrentPump
	pump.Type = "CurrentPump"
	pump.Units = "%"

	return &pump
}
