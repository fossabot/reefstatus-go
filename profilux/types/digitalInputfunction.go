package types

type DigitalInputFunction string

const (
	DigitalInputFunctionNotUsed      = "Not Used"
	DigitalInputFunctionLevelSensor  = "Level Sensor"
	DigitalInputFunctionWaterChange  = "Water Change"
	DigitalInputFunctionMaintenance  = "Maintenance"
	DigitalInputFunctionFeedingPause = "Feeding Pause"
	DigitalInputFunctionThunderstorm = "Thunderstorm"
)

var digitalInputFunctionMap = map[int]string{
	0: DigitalInputFunctionNotUsed,
	1: DigitalInputFunctionLevelSensor,
	2: DigitalInputFunctionWaterChange,
	3: DigitalInputFunctionMaintenance,
	4: DigitalInputFunctionFeedingPause,
	5: DigitalInputFunctionThunderstorm,
}

func GetDigitalInputFunction(value int) string {
	if val, ok := digitalInputFunctionMap[value]; ok {
		return val
	}

	return "Unknown"
}
