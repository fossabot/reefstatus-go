package models

import (
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
	"time"
)

type Info struct {
	Maintenance     []*Maintenance
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
	Reminders       []*Reminder
}

func NewInfo() *Info {
	var info Info
	info.Maintenance = make([]*Maintenance, 0)
	info.Reminders = make([]*Reminder, 0)

	return &info
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

func (info *Info) UpdateState(controller *profilux.Controller) {
	info.LastUpdate = time.Now()
	info.Alarm = controller.GetAlarm()
	info.OperationMode = controller.GetOperationMode()
	info.MoonPhase = controller.GetMoonPhase()

	for _, item := range info.Maintenance {
		item.UpdateState(controller)
	}

	for _, item := range info.Reminders {
		item.UpdateState(controller)
	}
}

func (info *Info) updateMaintenanceMode(controller *profilux.Controller, index int) {
	var maintenance *Maintenance = nil
	for _, item := range info.Maintenance {
		if item.Index == index {
			maintenance = item
			break
		}
	}

	if maintenance == nil {
		maintenance = NewMaintenance(index)
		info.Maintenance = append(info.Maintenance, maintenance)
	}

	maintenance.Update(controller)
}

func (info *Info) updateReminder(controller *profilux.Controller, index int) {
	var reminder *Reminder = nil
	var reminderIndex = 0
	for i, item := range info.Reminders {
		if item.Index == index {
			reminder = item
			reminderIndex = i
			break
		}
	}

	if !controller.IsReminderActive(index) {
		if reminder != nil {
			info.Reminders = append(info.Reminders[:reminderIndex], info.Reminders[reminderIndex+1:]...)
		}
		return
	}

	if reminder == nil {
		reminder = NewReminder(index)
		info.Reminders = append(info.Reminders, reminder)
	}

	reminder.Update(controller)
}
