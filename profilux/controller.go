package profilux

import (
	"github.com/cjburchell/reefstatus-go/profilux/code5"
	"github.com/cjburchell/reefstatus-go/profilux/types"
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
	p protocol
}

func NewController(settings ConnectionSettings) (*Controller, error) {
	con, err := newConnection(settings)
	if err != nil {
		return nil, err
	}

	var controller Controller
	controller.p.Address = settings.ControllerAddress
	controller.p.Connection = con

	return &controller, nil
}

func (controller *Controller) Disconnect() {
	controller.p.Connection.Disconnect()
}

func (controller *Controller) GetReminderCount() int {
	count, _ := controller.p.GetData(code5.GETREMINDERCOUNT)
	return count
}

// region Info

func (controller *Controller) GetSoftwareVersion() float64 {
	result, _ := controller.p.GetDataFloatAndRound(code5.SOFTWAREVERSION, 0.01, 2)
	return result
}

func (controller *Controller) GetModel() types.Model {
	result, _ := controller.p.GetDataEnum(code5.PRODUCTID, types.GetModel)
	return types.Model(result)
}

func (controller *Controller) GetSerialNumber() int {
	result, _ := controller.p.GetData(code5.SERIALNUMBER)
	return result
}

func (controller *Controller) GetSoftwareDate() time.Time {
	result, _ := controller.p.GetDataDate(code5.SOFTWAREDATE)
	return result
}

func (controller *Controller) GetDeviceAddress() int {
	result, _ := controller.p.GetData(code5.ADDRESS)
	return result
}

func (controller *Controller) GetLatitude() float64 {
	result, _ := controller.p.GetDataFloat(code5.LOC_LATITUDE, 0.1)
	return result
}

func (controller *Controller) GetLongitude() float64 {
	result, _ := controller.p.GetDataFloat(code5.LOC_LONGITUDE, 0.1)
	return result
}

func (controller *Controller) GetMoonPhase() float64 {
	result, _ := controller.p.GetDataFloat(code5.MOON_ACTPHASE, 1)
	return result
}

func (controller *Controller) GetAlarm() types.CurrentState {
	result, _ := controller.p.GetDataCurrentState(code5.ISALARM)
	return result
}

func (controller *Controller) GetOperationMode() types.OperationMode {
	result, _ := controller.p.GetDataEnum(code5.OPMODE, types.GetOperationMode)
	return types.OperationMode(result)
}

// endregion

// region Reminders

func (controller *Controller) IsReminderActive(index int) bool {
	result, _ := controller.p.GetData(code5.MEM1_NEXTMEM + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result != 0xffff
}

func (controller *Controller) GetReminderNext(index int) time.Time {
	result, _ := controller.p.GetData(code5.MEM1_NEXTMEM + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	nextReminder := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	nextReminder = nextReminder.AddDate(0, 0, result)
	return nextReminder
}

func (controller *Controller) GetReminderIsRepeating(index int) bool {
	result, _ := controller.p.GetDataBool(code5.MEM1_REPEAT + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result
}

func (controller *Controller) GetReminderPeriod(index int) int {
	result, _ := controller.p.GetData(code5.MEM1_REPEAT + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result
}

func (controller *Controller) GetReminderText(index int) string {
	result, _ := controller.p.GetDataText(code5.MEM1_TEXT01 + getOffset(index, BlockitemsReminder, BlocksizeReminder))
	return result
}

// endregion

// region Maintenance

func (controller *Controller) IsMaintenanceActive(index int) bool {
	result, _ := controller.p.GetDataBool(code5.MAINTENANCE_ISACTIVE + getOffset(index, 1, 2))
	return result
}

func (controller *Controller) GetMaintenanceDuration(index int) int {
	result, _ := controller.p.GetDataMultiplier(code5.MAINTENANCE_TIMEOUT+getOffset(index, 1, 27), 60)
	return result
}

func (controller *Controller) GetMaintenanceTimeLeft(index int) int {
	result, _ := controller.p.GetDataMultiplier(code5.MAINTENANCE_REMATINGTIME+getOffset(index, 1, 2), 60)
	return result
}
func (controller *Controller) GetMaintenanceText(index int) string {
	result, _ := controller.p.GetDataText(code5.MAINT_NAME + getOffset(index, 64, 1))
	return result
}

// endregion

// region Probes
// endregion
