package profilux

import (
	"github.com/cjburchell/reefstatus-go/common"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/profilux/code5"
	"github.com/cjburchell/reefstatus-go/profilux/protocol"
	"github.com/cjburchell/reefstatus-go/profilux/protocol/http"
	"github.com/cjburchell/reefstatus-go/profilux/protocol/native"
	"github.com/cjburchell/reefstatus-go/profilux/settings"
	"github.com/cjburchell/reefstatus-go/profilux/types"
	"strconv"
	"time"
)

const (
	//Blockitems1To10Vint           = 10
	//BlockitemsIlluminationchannel = 8
	//BlockitemsProglogic           = 8
	BlockitemsReminder = 4
	//BlockitemsSensorstates        = 8
	//BlockitemsSwitchplug          = 24
	//BlockitemsTimer               = 12
	//Blocksize1To10Vint            = 3
	//BlocksizeIlluminationchannel  = 28
	//BlocksizeProglogic            = 4
	BlocksizeReminder = 12
	//BlocksizeSensorstates         = 8
	//BlocksizeSwitchplug           = 1
	//BlocksizeTimer                = 21
	MegablockSize  = 1000
	SfFeedPause    = 2
	SfMaintenance  = 1
	SfThunderstorm = 3
	SfWaterChange  = 0
)

func getOffset(index, blockCount, blockSize int) int {
	return ((index % blockCount) * blockSize) + ((index / blockCount) * MegablockSize)
}

type Controller struct {
	p                 protocol.IProtocol
	reminderCount     *int
	sensorCount       *int
	levelSensorCount  *int
	digitalInputCount *int
	timerCount        *int
	lightCount        *int
	pumpCount         *int
	logicCount        *int
	sPortCount        *int
	lPortCount        *int
	callCount         int
}

func NewController() (*Controller, error) {

	var p protocol.IProtocol
	var err error
	if settings.Connection.Protocol == settings.ProtocolHTTP {
		p, err = http.NewProtocol(settings.Connection)
		if err != nil {
			return nil, err
		}
	} else {
		p, err = native.NewProtocol(settings.Connection)
		if err != nil {
			return nil, err
		}
	}

	var controller Controller
	controller.p = p
	return &controller, nil
}

func (controller *Controller) Disconnect() {
	controller.p.Disconnect()
}

func (controller *Controller) ResetStats() {
	controller.callCount = 0
}

func (controller Controller) GetCallCount() int {
	return controller.callCount
}

// region Protocol

func (controller *Controller) getDataText(code int) (string, error) {
	controller.callCount++
	return controller.p.GetDataText(code)
}

func (controller *Controller) getData(code int) (int, error) {
	controller.callCount++
	return controller.p.GetData(code)
}

func (controller *Controller) getDataDate(code int) (time.Time, error) {
	result, err := controller.getData(code)
	if err != nil {
		return time.Now(), err
	}

	timeString := strconv.Itoa(result)

	if len(timeString) == 6 {
		yearValue, _ := strconv.Atoi(timeString[len(timeString)-2:])
		monthValue, _ := strconv.Atoi(timeString[len(timeString)-4 : len(timeString)-2])
		dateValue, _ := strconv.Atoi(timeString[:len(timeString)-4])
		return time.Date(yearValue+2000, time.Month(monthValue), dateValue, 0, 0, 0, 0, time.UTC), nil
	} else if len(timeString) == 7 {
		yearValue, _ := strconv.Atoi(timeString[len(timeString)-3:])
		monthValue, _ := strconv.Atoi(timeString[len(timeString)-5 : len(timeString)-3])
		dateValue, _ := strconv.Atoi(timeString[:len(timeString)-5])
		return time.Date(yearValue+2000, time.Month(monthValue), dateValue, 0, 0, 0, 0, time.UTC), nil
	}

	return time.Now(), err
}

func (controller *Controller) getDataEnum(code int, convert func(int) string) (string, error) {
	result, err := controller.getData(code)
	if err != nil {
		return "", err
	}

	return convert(result), nil
}

func (controller *Controller) getDataCurrentState(code int) (types.CurrentState, error) {
	result, err := controller.getData(code)
	if err != nil {
		return "", err
	}

	return types.GetCurrentState(result), nil
}

func (controller *Controller) getDataCurrentStateConvert(code int, convert func(int) int) (types.CurrentState, error) {
	result, err := controller.getData(code)
	if err != nil {
		return "", err
	}

	return types.GetCurrentState(convert(result)), nil
}

func (controller *Controller) getDataFloat(code int, multiplier float64) (float64, error) {
	result, err := controller.getData(code)
	if err != nil {
		return 0, err
	}

	return float64(result) * multiplier, nil
}

func (controller *Controller) getDataMultiplier(code int, multiplier int) (int, error) {
	result, err := controller.getData(code)
	if err != nil {
		return 0, err
	}

	return result * multiplier, nil
}

func (controller *Controller) getDataBool(code int) (bool, error) {
	result, err := controller.getData(code)
	if err != nil {
		return false, err
	}

	return result != 0, nil
}

func (controller *Controller) getDataBoolConvert(code int, convert func(int) int) (bool, error) {
	result, err := controller.getData(code)
	if err != nil {
		return false, err
	}

	return convert(result) != 0, nil
}

func (controller *Controller) getDataFloatAndRound(code int, multiplier float64, digits int) (float64, error) {
	result, err := controller.getDataFloat(code, multiplier)
	if err != nil {
		return 0, err
	}

	return common.Round(result, digits), nil
}

// endregion

// region Info

func (controller *Controller) GetSoftwareVersion() float64 {
	result, _ := controller.getDataFloatAndRound(code5.SOFTWAREVERSION, 0.01, 2)
	return result
}

func (controller *Controller) GetModel() types.Model {
	result, _ := controller.getDataEnum(code5.PRODUCTID, types.GetModel)
	return types.Model(result)
}

func (controller *Controller) GetSerialNumber() int {
	result, _ := controller.getData(code5.SERIALNUMBER)
	return result
}

func (controller *Controller) GetSoftwareDate() time.Time {
	result, _ := controller.getDataDate(code5.SOFTWAREDATE)
	return result
}

func (controller *Controller) GetDeviceAddress() int {
	result, _ := controller.getData(code5.ADDRESS)
	return result
}

func (controller *Controller) GetLatitude() float64 {
	result, _ := controller.getDataFloat(code5.LOC_LATITUDE, 0.1)
	return result
}

func (controller *Controller) GetLongitude() float64 {
	result, _ := controller.getDataFloat(code5.LOC_LONGITUDE, 0.1)
	return result
}

func (controller *Controller) GetMoonPhase() float64 {
	result, _ := controller.getDataFloat(code5.MOON_ACTPHASE, 1)
	return result
}

func (controller *Controller) GetAlarm() types.CurrentState {
	result, _ := controller.getDataCurrentState(code5.ISALARM)
	return result
}

func (controller *Controller) GetOperationMode() types.OperationMode {
	result, _ := controller.getDataEnum(code5.OPMODE, types.GetOperationMode)
	return types.OperationMode(result)
}

// endregion

// region Reminders

func (controller *Controller) GetReminderCount() int {

	if controller.reminderCount == nil {
		count, _ := controller.getData(code5.GETREMINDERCOUNT)
		controller.reminderCount = &count
	}

	return *controller.reminderCount
}

func (controller *Controller) IsReminderActive(index int) bool {
	result, _ := controller.getData(code5.MEM1_NEXTMEM + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result != 0xffff
}

func (controller *Controller) GetReminderNext(index int) time.Time {
	result, _ := controller.getData(code5.MEM1_NEXTMEM + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	nextReminder := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	nextReminder = nextReminder.AddDate(0, 0, result)
	return nextReminder
}

func (controller *Controller) GetReminderIsRepeating(index int) bool {
	result, _ := controller.getDataBool(code5.MEM1_REPEAT + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result
}

func (controller *Controller) GetReminderPeriod(index int) int {
	result, _ := controller.getData(code5.MEM1_DAYS + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result
}

func (controller *Controller) GetReminderText(index int) string {
	result, _ := controller.getDataText(code5.MEM1_TEXT01 + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result
}

// endregion

// region Maintenance

func (controller *Controller) IsMaintenanceActive(index int) bool {
	result, _ := controller.getDataBool(code5.MAINTENANCE_ISACTIVE + getOffset(index, 1, 2))
	return result
}

func (controller *Controller) GetMaintenanceDuration(index int) int {
	result, _ := controller.getDataMultiplier(code5.MAINTENANCE_TIMEOUT+getOffset(index, 1, 27), 60)
	return result
}

func (controller *Controller) GetMaintenanceTimeLeft(index int) int {
	result, _ := controller.getDataMultiplier(code5.MAINTENANCE_REMATINGTIME+getOffset(index, 1, 2), 60)
	return result
}
func (controller *Controller) GetMaintenanceText(index int) string {
	result, _ := controller.getDataText(code5.MAINT_NAME + getOffset(index, 64, 1))
	return result
}

// endregion

// region Probes

func getSensorOffset(index int) int {
	return getOffset(index, 8, 24)
}

func (controller *Controller) GetSensorCount() int {

	if controller.sensorCount == nil {
		count, _ := controller.getData(code5.GETSENSORCOUNT)
		controller.sensorCount = &count
	}

	return *controller.sensorCount
}

func (controller *Controller) GetSensorType(index int) types.SensorType {
	result, _ := controller.getDataEnum(code5.SENSORPARA1_SENSORTYPE+getSensorOffset(index), types.GetSensorType)
	return types.SensorType(result)
}

func (controller *Controller) GetSensorMode(index int) types.SensorMode {
	result, _ := controller.getDataEnum(code5.SENSORPARA1_CAL1ADC+getSensorOffset(index), func(value int) string {
		return types.GetSensorMode(value >> 12)
	})
	return types.SensorMode(result)
}

func (controller *Controller) GetSensorActive(index int) bool {
	props, _ := controller.getData(code5.SENSORPARA1_PROPS + getSensorOffset(index))
	return props&0x1 == 1
}

func (controller *Controller) GetProbeName(index int) string {
	name, _ := controller.getDataText(code5.SENSOR1_NAME + getOffset(index, 32, 1))
	return name
}

func (controller *Controller) GetSensorNominalValue(index int, multiplier float64) float64 {
	result, _ := controller.getDataFloat(code5.SENSORPARA1_DESVALUE+getSensorOffset(index), multiplier)
	return result
}

func (controller *Controller) GetSensorAlarmDeviation(index int, multiplier float64) float64 {
	result, _ := controller.getDataFloat(code5.SENSORPARA1_ALARMDEVIATION+getSensorOffset(index), multiplier)
	return result
}

func (controller *Controller) GetSensorAlarmEnable(index int) bool {
	result, _ := controller.getDataBool(code5.SENSORPARA1_ALARMMODE + getSensorOffset(index))
	return result
}

func (controller *Controller) GetSensorAlarm(index int) types.CurrentState {
	result, _ := controller.getDataCurrentStateConvert(code5.SENSORPARA1_ACTSTATE+getOffset(index, 8, 8), func(i int) int {
		return i & 0x100
	})
	return result
}

func (controller *Controller) GetSensorValue(index int, multiplier float64) float64 {
	result, _ := controller.getDataFloat(code5.SENSORPARA1_ACTVALUE+getOffset(index, 8, 8), multiplier)
	return result
}

func (controller *Controller) GetProbeOperationHours(index int) int {
	result, _ := controller.getData(code5.SENSORPARA1_OHM + getOffset(index, 8, 8))
	return result
}

func (controller *Controller) GetSensorFormat(index int) int {
	result, _ := controller.getData(code5.SENSORPARA1_DISPLAYMODE + getSensorOffset(index))
	return result & 0x0f
}

// endregion

// region LevelSensor

func (controller *Controller) GetLevelSenosrCount() int {
	if controller.levelSensorCount == nil {
		count, _ := controller.getData(code5.GETLEVELSENSORCOUNT)
		controller.levelSensorCount = &count
	}

	return *controller.levelSensorCount
}

func (controller *Controller) GetLevelSensorMode(index int) types.LevelSensorOperationMode {
	result, _ := controller.getDataEnum(code5.LEVEL1_PROPS+getOffset(index, 3, 3), func(value int) string {

		return types.GetLevelSensorOperationMode(value >> 13)
	})
	return types.LevelSensorOperationMode(result)
}

func (controller *Controller) GetLevelName(index int) string {
	name, _ := controller.getDataText(code5.LEVEL_NAME + getOffset(index, 64, 1))
	return name
}

type LevelState struct {
	WaterMode types.WaterMode
	Drain     types.CurrentState
	Fill      types.CurrentState
	Alarm     types.CurrentState
}

func (controller *Controller) GetLevelSensorState(index int) LevelState {
	// 7654 3210
	// AFDW WWWR
	// A - Alarm
	// F - Fill
	// D - Drain
	// W - Water State
	// R - Reserverd
	var levelState LevelState
	state, _ := controller.getData(code5.LEVEL1_ACTSTATE + getOffset(index, 3, 1))
	state >>= 1
	levelState.WaterMode = types.GetWaterMode(state & 0xf)
	state >>= 4
	levelState.Drain = types.GetCurrentState(state & 0x1)
	state >>= 1
	levelState.Fill = types.GetCurrentState(state & 0x1)
	state >>= 1
	levelState.Alarm = types.GetCurrentState(state & 0x1)

	return levelState
}

func (controller *Controller) GetLevelSource1(index int) int {
	result, _ := controller.getData(code5.LEVEL1_SOURCES + getOffset(index, 3, 3))
	return result & 0xf

}

func (controller *Controller) GetLevelSource2(index int) int {
	result, _ := controller.getData(code5.LEVEL1_SOURCES + getOffset(index, 3, 3))
	result >>= 4
	return result & 0xf
}

type LevelInputState struct {
	Delayed   types.CurrentState
	Previous  types.CurrentState
	Undelayed types.CurrentState
}

func (controller *Controller) GetLevelSensorCurrentState(index int) LevelInputState {
	// 7654 3210
	// UPDR RRRR
	// U - undelayed
	// P-  Previous
	// D - Delayed
	// R - Reserverd
	var levelState LevelInputState
	state, _ := controller.getData(code5.LEVEL1_INPUT_STATE + getOffset(index, 4, 1))
	state >>= 5
	levelState.Delayed = types.GetCurrentState(state & 0x1)
	state >>= 1
	levelState.Previous = types.GetCurrentState(state & 0x1)
	state >>= 1
	levelState.Undelayed = types.GetCurrentState(state & 0x1)

	return levelState
}

// endregion

// region DigitalInput
func (controller *Controller) GetDigitalInputCount() int {

	if controller.digitalInputCount == nil {
		count, _ := controller.getData(code5.GETDIGITALINPUTCOUNT)
		controller.digitalInputCount = &count
	}

	return *controller.digitalInputCount
}

func (controller *Controller) GetDigitalInputFunction(index int) types.DigitalInputFunction {
	result, _ := controller.getDataEnum(code5.DIGITALINPUT1_FUNCTION+getOffset(index, 4, 1), types.GetDigitalInputFunction)
	return types.DigitalInputFunction(result)
}

func (controller *Controller) GetDigitalInputState(index int) types.CurrentState {
	state, _ := controller.getData(code5.DIGITALINPUTSSTATE + getOffset(index, 4, 0))
	switch index % 4 {
	case 0:
		return types.GetCurrentState(state % 0x1)
	case 1:
		return types.GetCurrentState(state % 0x2)
	case 2:
		return types.GetCurrentState(state % 0x4)
	case 3:
		return types.GetCurrentState(state % 0x8)
	}

	return types.CurrentStateOff
}

// endregion

// region Timer

func (controller *Controller) GetTimerCount() int {

	if controller.timerCount == nil {
		count, _ := controller.getData(code5.GETTIMERCOUNT)
		controller.timerCount = &count
	}

	return *controller.timerCount
}

type TimerSettings struct {
	FeedPauseIfActive bool
	Mode              types.TimerMode
	DayMode           types.DayMode
	SwitchingCount    int
}

func (controller *Controller) GetTimerSettings(index int) TimerSettings {
	config, _ := controller.getData(code5.TIMER1_PROPS + getOffset(index, 12, 21))

	return TimerSettings{
		Mode:              types.GetTimerMode((config >> 7) & 0x7),
		SwitchingCount:    config >> 11,
		FeedPauseIfActive: (config & 0x10) != 0,
		DayMode:           types.GetDayMode((config >> 10) & 0x1),
	}
}

func (controller *Controller) GetDosingRate(index int) int {
	result, _ := controller.getData(code5.TIMER1_RATEPERDOSING + getOffset(index, 12, 21))
	return result
}

// endregion

// region Light
func (controller *Controller) GetLightCount() int {
	if controller.lightCount == nil {
		count, _ := controller.getData(code5.GETILLUMINATIONCOUNT)
		controller.lightCount = &count
	}

	return *controller.lightCount
}

func (controller *Controller) GetIsLightActive(index int) bool {
	result, _ := controller.getDataBoolConvert(code5.ILLUMINATION1_PROPS+getOffset(index, 8, 28), func(config int) int {
		return (config >> 4) & 0xF
	})
	return result
}
func (controller *Controller) GetIsLightDimmable(index int) bool {
	result, _ := controller.getDataBoolConvert(code5.ILLUMINATION1_PROPS+getOffset(index, 8, 28), func(config int) int {
		return config & 0x8
	})
	return result
}
func (controller *Controller) GetLightOperationHours(index int) int {
	result, _ := controller.getData(code5.ILLUMINATION1_OHM + getOffset(index, 8, 4))
	return result
}
func (controller *Controller) GetLightValue(index int) float64 {
	result, _ := controller.getDataFloat(code5.ILLUMINATION1_ACTPERCENT+getOffset(index, 8, 4), 1)
	return result
}
func (controller *Controller) GetLightName(index int) string {
	result, _ := controller.getDataText(code5.ILLUMINATION1_NAME + getOffset(index, 32, 1))
	return result
}

// endregion

// region Pumps

func (controller *Controller) GetCurrentPumpCount() int {

	if controller.pumpCount == nil {
		count, _ := controller.getData(code5.GETCURRENTPUMPCOUNT)
		controller.pumpCount = &count
	}

	return *controller.pumpCount
}

func (controller *Controller) GetIsCurrentPumpAssigned(index int) bool {
	group1Mask, _ := controller.getData(code5.CURRENTCONTROL_GROUP1PUMPCOUNT)
	i := uint(index)
	if (group1Mask >> i & 0x1) == 1 {
		return true
	}

	group2Mask, _ := controller.getData(code5.CURRENTCONTROL_GROUP2PUMPCOUNT)
	if (group2Mask >> i & 0x1) == 1 {
		return true
	}

	return false
}

func (controller *Controller) GetCurrentPumpValue(index int) int {
	result, _ := controller.getData(code5.CURRENTPUMP1_ACTPERCENT + getOffset(index, 4, 1))
	return result
}

// endregion

// region ProgrammableLogic
func (controller *Controller) GetProgrammableLogicCount() int {

	if controller.logicCount == nil {
		count, _ := controller.getData(code5.GETPROGLOGICCOUNT)
		controller.logicCount = &count
	}

	return *controller.logicCount
}

func (controller *Controller) GetProgramLogicInput(input, index int) PortMode {
	mode, _ := controller.getData(code5.PROGLOGIC1_INPUT1 + input + getOffset(index, 8, 4))
	// mode format
	// 1234 1234 1234 1234
	// RRRR RRPP PPPT TTTT
	// R = Reserved
	// P = Port Number
	// T = Type
	var portMode PortMode
	mode >>= 6
	portMode.Port = (mode & 0x1F) + 1
	mode >>= 5
	portMode.DeviceMode = types.GetDeviceMode(mode & 0x1F)
	portMode.IsProbe = getIsProbe(portMode.DeviceMode)
	return portMode
}

type LogicFunction struct {
	Invert1   bool
	Invert2   bool
	LogicMode types.LogicMode
}

func (controller *Controller) GetProgramLogicFunction(index int) LogicFunction {
	mode, _ := controller.getData(code5.PROGLOGIC1_FUNCTION + getOffset(index, 8, 4))

	// mode format
	// 1234 1234
	// RRMM MMMM
	// R = Reserved
	// P = Port Number
	// T = Type
	var function LogicFunction
	function.Invert2 = (mode & 0x1) == 1
	mode >>= 1
	function.Invert1 = (mode & 0x1) == 1
	mode >>= 1
	function.LogicMode = types.GetLogicMode(mode & 0x3F)

	return function
}

// endregion

// region sport

func (controller *Controller) GetSPortCount() int {

	if controller.sPortCount == nil {
		count, _ := controller.getData(code5.GETSWITCHCOUNT)
		controller.sPortCount = &count
	}

	return *controller.sPortCount
}

type PortMode struct {
	DeviceMode types.DeviceMode
	Port       int
	BlackOut   int
	Invert     bool
	Id         string
	IsProbe    bool
}

func getIsProbe(mode types.DeviceMode) bool {
	return mode == types.DeviceModeDecrease ||
		mode == types.DeviceModeIncrease ||
		mode == types.DeviceModeSubstrate ||
		mode == types.DeviceModeProbeAlarm
}

func (controller *Controller) GetSPortFunction(index int) PortMode {
	mode, _ := controller.getData(code5.SWITCHPLUG1_FUNCTION + getOffset(index, 24, 1))

	// mode format
	// 1234 1234 1234 1234
	// IBBB BBPP PPPT TTTT
	// I = Invert
	// B = Blackout time
	// P = Port Number
	// T = Type
	var portMode PortMode
	portMode.Invert = (mode & 0x01) != 0
	mode >>= 1
	portMode.BlackOut = mode & 0x1F
	mode >>= 5
	portMode.Port = (mode & 0x1F) + 1
	mode >>= 5
	portMode.DeviceMode = types.GetDeviceMode(mode & 0x1F)
	portMode.IsProbe = getIsProbe(portMode.DeviceMode)
	return portMode
}

func (controller *Controller) GetSPortValue(index int) types.CurrentState {
	value, _ := controller.getDataCurrentState(code5.SP1_STATE + getOffset(index, 24, 1))
	return value
}

func (controller *Controller) GetSPortName(index int) string {
	result, _ := controller.getDataText(code5.SWITCHPLUG1_NAME + getOffset(index, 64, 1))
	return result
}

// endregion

// region lport

func (controller *Controller) GetLPortCount() int {

	if controller.lPortCount == nil {
		count, _ := controller.getData(code5.GETONTETOTENVINTCOUNT)
		controller.lPortCount = &count
	}

	return *controller.lPortCount
}

func (controller *Controller) GetLPortFunction(index int) PortMode {
	mode, _ := controller.getData(code5.L1TO10VINT1_FUNCTION + getOffset(index, 10, 3))

	// mode format
	// 1234 1234
	// PPPT TTTT
	// I = Invert
	// B = Blackout time
	// P = Port Number
	// T = Type
	var portMode PortMode
	portMode.BlackOut = mode & 0x003F
	mode >>= 6
	portMode.Port = (mode & 0x001F) + 1
	mode >>= 5
	portMode.DeviceMode = types.LPortModeToSocketType(mode & 0x003F)
	portMode.IsProbe = getIsProbe(portMode.DeviceMode)
	return portMode
}

const (
	LValueMin = 18.0
	LValueMax = 255.00
)

func (controller *Controller) GetLPortValue(index int) float64 {

	value, _ := controller.getData(code5.L1TO10VINT1_PWMVALUE + getOffset(index, 10, 1))

	if value < LValueMin {
		return 0
	} else {
		return (float64(value) - LValueMin) / (LValueMax - LValueMin) * 100.0
	}
}

// endregion

func (controller *Controller) FeedPause(index int, activate bool) error {
	command := (index << 16) | (0 << 8) | SfFeedPause
	if activate {
		command = (index << 16) | (0xFF << 8) | SfFeedPause
	}

	err := controller.p.SendData(code5.INVOKESPECIALFUNCTION, command)
	if err != nil {
		log.Errorf(err, "FeedPause: %s", err.Error())
		return err
	}

	return nil
}

func (controller *Controller) ClearReminder(reminder int) error {
	offset := getOffset(reminder, BlockitemsReminder, BlocksizeReminder)
	err := controller.p.SendData(code5.MEM1_NEXTMEM+offset, 0xFFFF)
	if err != nil {
		log.Errorf(err, "ClearReminder: %s", err.Error())
		return err
	}

	return nil
}

func (controller *Controller) ResetReminder(reminder int, period int) error {
	nextReminder := time.Now().AddDate(0, 0, period)
	span := nextReminder.Sub(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))
	data := int(span.Hours() / 24)

	offset := getOffset(reminder, BlockitemsReminder, BlocksizeReminder)
	err := controller.p.SendData(code5.MEM1_NEXTMEM+offset, data)
	if err != nil {
		log.Errorf(err, "ClearReminder: %s", err.Error())
		return err
	}

	return nil
}

func (controller *Controller) ClearLevelAlarm(index int) error {
	err := controller.p.SendData(code5.LEVEL1_ACTSTATE+getOffset(index, 3, 1), 0)

	if err != nil {
		log.Errorf(err, "ClearLevelAlarm: %s", err.Error())
		return err
	}

	return nil
}

func (controller *Controller) Thunderstorm(duration int) error {
	command := (duration << 8) | SfThunderstorm
	err := controller.p.SendData(code5.INVOKESPECIALFUNCTION, command)
	if err != nil {
		log.Errorf(err, "Thunderstorm: %s", err.Error())
		return err
	}

	return nil
}

func (controller *Controller) WaterChange(index int) error {
	command := index<<16 | (0xFF << 8) | SfWaterChange
	err := controller.p.SendData(code5.INVOKESPECIALFUNCTION, command)
	if err != nil {
		log.Errorf(err, "WaterChange: %s", err.Error())
		return err
	}

	return nil
}

func (controller *Controller) Maintenance(activate bool, index int) error {
	command := (index << 16) | (0 << 8) | SfMaintenance
	if activate {
		command = (index << 16) | (0xFF << 8) | SfMaintenance
	}

	err := controller.p.SendData(code5.INVOKESPECIALFUNCTION, command)

	if err != nil {
		log.Errorf(err, "Maintenance: %s", err.Error())
		return err
	}

	return nil
}
