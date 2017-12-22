package profilux

import (
	"github.com/cjburchell/reefstatus-go/common"
	"github.com/cjburchell/reefstatus-go/profilux/code5"
	"github.com/cjburchell/reefstatus-go/profilux/types"
	"strconv"
	"time"
)

const (
	Blockitems1To10Vint           = 10
	BlockitemsIlluminationchannel = 8
	BlockitemsProglogic           = 8
	BlockitemsReminder            = 4
	BlockitemsSensorstates        = 8
	BlockitemsSwitchplug          = 24
	BlockitemsTimer               = 12
	Blocksize1To10Vint            = 3
	BlocksizeIlluminationchannel  = 28
	BlocksizeProglogic            = 4
	BlocksizeReminder             = 12
	BlocksizeSensorstates         = 8
	BlocksizeSwitchplug           = 1
	BlocksizeTimer                = 21
	MegablockSize                 = 1000
	SfFeedpause                   = 2
	SfMaintenance                 = 1
	SfThunderstorm                = 3
	SfWaterchange                 = 0
)

func getOffset(index, blockCount, blockSize int) int {
	return ((index % blockCount) * blockSize) + ((index / blockCount) * MegablockSize)
}

type Controller struct {
	p iProtocol
}

func NewController(settings ConnectionSettings) (*Controller, error) {
	p, err := newHttpProtocol(settings)
	if err != nil {
		return nil, err
	}

	var controller Controller
	controller.p = p
	return &controller, nil
}

func (controller *Controller) Disconnect() {
	controller.p.Disconnect()
}

// region Protocol

func (controller *Controller) getDataText(code int) (string, error) {
	return controller.p.GetDataText(code)
}

func (controller *Controller) getData(code int) (int, error) {
	return controller.p.GetData(code)
}

func (controller *Controller) getDataDate(code int) (time.Time, error) {
	result, err := controller.p.GetData(code)
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
	result, err := controller.p.GetData(code)
	if err != nil {
		return "", err
	}

	return convert(result), nil
}

func (controller *Controller) getDataCurrentState(code int) (types.CurrentState, error) {
	result, err := controller.p.GetData(code)
	if err != nil {
		return "", err
	}

	return types.GetCurrentState(result), nil
}

func (controller *Controller) getDataFloat(code int, multiplier float64) (float64, error) {
	result, err := controller.p.GetData(code)
	if err != nil {
		return 0, err
	}

	return float64(result) * multiplier, nil
}

func (controller *Controller) getDataMultiplier(code int, multiplier int) (int, error) {
	result, err := controller.p.GetData(code)
	if err != nil {
		return 0, err
	}

	return result * multiplier, nil
}

func (controller *Controller) getDataBool(code int) (bool, error) {
	result, err := controller.p.GetData(code)
	if err != nil {
		return false, err
	}

	return result != 0, nil
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
	count, _ := controller.getData(code5.GETREMINDERCOUNT)
	return count
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
	result, _ := controller.getData(code5.MEM1_REPEAT + getOffset(index, BlockitemsReminder, BlocksizeReminder))
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
// endregion
