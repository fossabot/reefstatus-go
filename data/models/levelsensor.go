package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type LevelSensor struct {
	SensorInfo
	OperationMode types.LevelSensorOperationMode
	Value         types.CurrentState
	WaterMode     types.WaterMode
	SecondSensor  types.CurrentState
}

func NewLevelSensor() *LevelSensor {
	var sensor LevelSensor
	sensor.Type = "LevelSensor"
	sensor.SensorType = types.SensorTypeLevel
	sensor.Units = "State"
	return &sensor
}

func (sensor LevelSensor) HasTwoInputs() bool {
	return sensor.OperationMode == types.LevelSensorAutoTopOffWith2Sensors ||
		sensor.OperationMode == types.LevelSensorWaterChangeAndAutoTopOff ||
		sensor.OperationMode == types.LevelSensorWaterChange ||
		sensor.OperationMode == types.LevelSensorMinMaxControl
}

func (sensor LevelSensor) HasWaterChange() bool {
	return sensor.OperationMode == types.LevelSensorWaterChangeAndAutoTopOff ||
		sensor.OperationMode == types.LevelSensorWaterChange
}
