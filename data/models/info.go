package models

import (
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
	"time"
)

type Info struct {
	Maintenance     map[int]*Maintenance
	OperationMode   types.OperationMode
	Model           types.Model
	SoftwareDate    time.Time
	DeviceAddress   int
	Latitude        float64
	Longitude       float64
	MoonPhase       float64
	Alarm           types.CurrentState
	SoftwareVersion float64
	SerialNumber    int
	LastUpdate      time.Time
	Reminders       map[int]*Reminder
}

func (info Info) IsP3() bool {
	return info.Model == types.ProfiLuxIII || info.Model == types.ProfiLuxIIIEx
}

func (info *Info) Update(controller *profilux.Controller) {
	info.LastUpdate = time.Now()

	info.SoftwareVersion = controller.GetSoftwareVersion()
	info.Model = controller.GetModel()
	info.SerialNumber = controller.GetSerialNumber()
	info.SoftwareDate = controller.GetSoftwareDate()
	info.DeviceAddress = controller.GetDeviceAddress()
	info.Latitude = controller.GetLatitude()
	info.Longitude = controller.GetLongitude()
	info.MoonPhase = controller.GetMoonPhase()
	info.Alarm = controller.GetAlarm()
	info.OperationMode = controller.GetOperationMode()

	info.updateMaintenanceMode(controller, 0)
	info.updateMaintenanceMode(controller, 1)
	info.updateMaintenanceMode(controller, 2)
	info.updateMaintenanceMode(controller, 3)

	for i := 0; i < controller.GetReminderCount(); i++ {
		info.updateReminder(controller, i)
	}
}

func (info *Info) updateMaintenanceMode(controller *profilux.Controller, index int) {
	maintenance, found := info.Maintenance[index]
	if !found {
		info.Maintenance[index] = NewMaintenance(index)
		maintenance = info.Maintenance[index]
	}

	maintenance.Update(controller)
}

func (info *Info) updateReminder(controller *profilux.Controller, index int) {
	reminder, found := info.Reminders[index]
	if !controller.IsReminderActive(index) {
		if found {
			delete(info.Reminders, index)
		}
		return
	}

	if !found {
		reminder = NewReminder(index)
		info.Reminders[index] = reminder
	}

	reminder.Update(controller)
}
