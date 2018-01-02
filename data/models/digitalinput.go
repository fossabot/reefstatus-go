package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
)

type DigitalInput struct {
	SensorInfo
	Value    types.CurrentState
	Function types.DigitalInputFunction
}

func NewDigitalInput(index int) *DigitalInput {
	var digitalInput DigitalInput
	digitalInput.Index = index
	digitalInput.Type = "DigitalInput"
	digitalInput.SensorType = types.SensorTypeDigitalInput
	digitalInput.Units = "State"
	digitalInput.Id = fmt.Sprintf("DigitalInput%d", 1+index)

	return &digitalInput
}

func (sensor *DigitalInput) Update(controller *profilux.Controller) {
	sensor.Function = controller.GetDigitalInputFunction(sensor.Index)
	sensor.Value = controller.GetDigitalInputState(sensor.Index)
}

func (sensor *DigitalInput) UpdateState(controller *profilux.Controller) {
	sensor.Value = controller.GetDigitalInputState(sensor.Index)
}
