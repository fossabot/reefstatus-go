package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type DigitalInput struct {
	SensorInfo
	Value    types.CurrentState
	Function types.DigitalInputFunction
}

func NewDigitalInput() *DigitalInput {
	var digitalInput DigitalInput
	digitalInput.Type = "DigitalInput"
	digitalInput.SensorType = types.SensorTypeDigitalInput
	digitalInput.Units = "State"
	return &digitalInput
}
