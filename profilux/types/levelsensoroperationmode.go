package types

type LevelSensorOperationMode string

const (
	LevelSensorNotEnabled               = "Not Enabled"
	LevelSensorAutoTopOff               = "Auto Top Off"
	LevelSensorMinMaxControl            = "Min-Max Control"
	LevelSensorWaterChange              = "Water Change"
	LevelSensorLeakageDetection         = "Leakage Detection"
	LevelSensorWaterChangeAndAutoTopOff = "Water Change and Auto Top Off"
	LevelSensorAutoTopOffWith2Sensors   = "Auto Top Off With 2 Sensors"
	LevelSensorReturnPump               = "Return Pump"
)

var levelSensorOperationModeMap = map[int]string{
	0: LevelSensorNotEnabled,
	1: LevelSensorAutoTopOff,
	2: LevelSensorMinMaxControl,
	3: LevelSensorWaterChange,
	4: LevelSensorLeakageDetection,
	5: LevelSensorWaterChangeAndAutoTopOff,
	6: LevelSensorAutoTopOffWith2Sensors,
	7: LevelSensorReturnPump,
}

func GetLevelSensorOperationMode(value int) string {
	if val, ok := levelSensorOperationModeMap[value]; ok {
		return val
	}

	return "Unknown"
}
