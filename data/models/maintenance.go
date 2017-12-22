package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type Maintenance struct {
	DisplayName string
	Index       int
	IsActive    bool
	Duration    int
	TimeLeft    int
}

func NewMaintenance(index int) *Maintenance {
	var maintenance Maintenance
	maintenance.Index = index
	maintenance.DisplayName = fmt.Sprintf("Maintenance%d", index+1)
	return &maintenance
}

func (maintenance *Maintenance) Update(controller *profilux.Controller) {
	maintenance.IsActive = controller.IsMaintenanceActive(maintenance.Index)
	maintenance.Duration = controller.GetMaintenanceDuration(maintenance.Index)
	maintenance.TimeLeft = controller.GetMaintenanceTimeLeft(maintenance.Index)

	name := controller.GetMaintenanceText(maintenance.Index)
	if len(name) != 0 {
		maintenance.DisplayName = name
	}
}
