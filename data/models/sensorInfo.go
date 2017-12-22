package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type SensorInfo struct {
	BaseInfo
	Format     int
	SensorType types.SensorType
	Index      int
	AlarmState types.CurrentState
}

func (sensor *SensorInfo) SetFormat(format int) {
	sensor.Format = format
	sensor.updateUnits()
}

func (sensor *SensorInfo) SetSensorType(sensorType types.SensorType) {
	sensor.SensorType = sensorType
	sensor.updateUnits()
}

func (sensor *SensorInfo) updateUnits() {
	switch sensor.SensorType {
	case types.SensorTypeLevel:
		sensor.Units = "State"
		break
	case types.SensorTypePH:
		sensor.Units = "PH"
		break
	case types.SensorTypeAirTemperature:
	case types.SensorTypeTemperature:
		if sensor.Format == 1 {
			sensor.Units = "°F"
		} else {
			sensor.Units = "°C"
		}

		break
	case types.SensorTypeConductivityF:
		sensor.Units = "μS"
		break
	case types.SensorTypeConductivity:
		if sensor.Format == 1 {
			sensor.Units = "ppt/PSU"
		} else if sensor.Format == 2 {
			sensor.Units = "SG"
		} else {
			sensor.Units = "mS"
		}

		break
	case types.SensorTypeRedox:
		sensor.Units = "mV"
		break
	case types.SensorTypeOxygen:
	case types.SensorTypeHumidity:
		sensor.Units = "%"
		break
	case types.SensorTypeVoltage:
		sensor.Units = "V"
		break
	}
}
