package types

type DeviceMode string

const (
	DeviceModeLights               = "Lights"
	DeviceModeTimer                = "Timer"
	DeviceModeDecrease             = "Decrease"
	DeviceModeIncrease             = "Increase"
	DeviceModeSubstrate            = "Substrate"
	DeviceModeProbeAlarm           = "ProbeAlarm"
	DeviceModeWater                = "Water"
	DeviceModeCurrentPump          = "CurrentPump"
	DeviceModeProgrammableLogic    = "ProgrammableLogic"
	DeviceModeVariableIllumination = "VariableIllumination"
	DeviceModeTempPTC              = "TempPTC"
	DeviceModeDigitalInput         = "DigitalInput"
	DeviceModeMaintenance          = "Maintenance"
	DeviceModeThunderStorm         = "ThunderStorm"
	DeviceModeWaterChange          = "WaterChange"
	DeviceModeFilter               = "Filter"
	DeviceModeAlarm                = "Alarm"
	DeviceModeAlwaysOn             = "AlwaysOn"
	DeviceModeAlwaysOff            = "AlwaysOff"
	DeviceModeThunder              = "Thunder"
)

var deviceModeMap = map[int]string{
	0:  DeviceModeLights,
	1:  DeviceModeTimer,
	2:  DeviceModeDecrease,
	3:  DeviceModeIncrease,
	4:  DeviceModeSubstrate,
	5:  DeviceModeProbeAlarm,
	6:  DeviceModeWater,
	7:  DeviceModeCurrentPump,
	9:  DeviceModeProgrammableLogic,
	10: DeviceModeVariableIllumination,
	11: DeviceModeTempPTC,
	12: DeviceModeDigitalInput,
	13: DeviceModeMaintenance,
	25: DeviceModeThunderStorm,
	26: DeviceModeWaterChange,
	27: DeviceModeFilter,
	28: DeviceModeAlarm,
	29: DeviceModeAlwaysOn,
	30: DeviceModeAlwaysOff,
	31: DeviceModeThunder,
}

func GetDeviceMode(mode int) string {
	if val, ok := deviceModeMap[mode]; ok {
		return val
	}

	return "Unknown"
}
