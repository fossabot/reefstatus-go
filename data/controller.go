package main

import (
	"github.com/cjburchell/reefstatus-go/data/models"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type Controller struct {
	Info              models.Info
	DigitalInputs     map[int]*models.DigitalInput
	DosingPumps       map[int]*models.DosingPump
	LevelSensors      map[int]*models.LevelSensor
	Lights            map[int]*models.Light
	LPorts            map[int]*models.LPort
	Probes            map[int]*models.Probe
	ProgrammableLogic map[int]*models.ProgrammableLogic
	Pumps             map[int]*models.CurrentPump
	SPorts            map[int]*models.SPort
}

func (controller *Controller) Update(profiluxController *profilux.Controller) {
	controller.Info.Update(profiluxController)
	controller.UpdateProbes(profiluxController)
}

func (controller *Controller) UpdateProbes(profiluxController *profilux.Controller) {
	/*for i := 0; i < profiluxController.GetSensorCount(); i++ {
		sensorType := profiluxController.GetSensorType(i)
	}*/
}
