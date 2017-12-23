package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
)

type LevelSensor struct {
	SensorInfo
	OperationMode     types.LevelSensorOperationMode
	Value             types.CurrentState
	SensorIndex       int
	WaterMode         types.WaterMode
	SecondSensor      types.CurrentState
	SecondSensorIndex int
	HasTwoInputs      bool
	HasWaterChange    bool
}

func NewLevelSensor(index int) *LevelSensor {
	var sensor LevelSensor
	sensor.Index = index
	sensor.Type = "LevelSensor"
	sensor.SensorType = types.SensorTypeLevel
	sensor.Units = "State"
	sensor.Id = fmt.Sprintf("Level%d", 1+index)
	return &sensor
}

func (sensor LevelSensor) hasTwoInputs() bool {
	return sensor.OperationMode == types.LevelSensorAutoTopOffWith2Sensors ||
		sensor.OperationMode == types.LevelSensorWaterChangeAndAutoTopOff ||
		sensor.OperationMode == types.LevelSensorWaterChange ||
		sensor.OperationMode == types.LevelSensorMinMaxControl
}

func (sensor LevelSensor) hasWaterChange() bool {
	return sensor.OperationMode == types.LevelSensorWaterChangeAndAutoTopOff ||
		sensor.OperationMode == types.LevelSensorWaterChange
}

func (sensor *LevelSensor) Update(controller *profilux.Controller) {
	sensor.OperationMode = controller.GetLevelSensorMode(sensor.Index)
	sensor.HasTwoInputs = sensor.hasTwoInputs()
	sensor.HasWaterChange = sensor.hasWaterChange()

	sensor.DisplayName = controller.GetLevelName(sensor.Index)

	state := controller.GetLevelSensorState(sensor.Index)
	sensor.AlarmState = state.Alarm
	sensor.WaterMode = state.WaterMode

	source1 := controller.GetLevelSource1(sensor.Index)
	sensorState := controller.GetLevelSensorCurrentState(source1)
	sensor.Value = sensorState.Undelayed
	sensor.SensorIndex = source1

	if sensor.HasTwoInputs {
		source2 := controller.GetLevelSource2(sensor.Index)
		sensorState2 := controller.GetLevelSensorCurrentState(source2)
		sensor.SecondSensor = sensorState2.Undelayed
		sensor.SecondSensorIndex = source2
	}
}
