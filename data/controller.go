package main

import (
	"github.com/cjburchell/reefstatus-go/data/models"
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
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
	controller.UpdateLevelSensors(profiluxController)
}

func (controller *Controller) UpdateProbes(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetSensorCount(); i++ {
		sensorType := profiluxController.GetSensorType(i)

		var mode types.SensorMode
		active := false
		if sensorType != types.SensorTypeFree && sensorType != types.SensorTypeNone {
			mode = profiluxController.GetSensorMode(i)
			active = profiluxController.GetSensorActive(i)
		}

		probe, found := controller.Probes[i]

		if active && mode == types.SensorModeNormal {
			if !found {
				probe = models.NewProbe(i)
				controller.Probes[i] = probe
			}

			probe.Update(profiluxController)

		} else {
			if found {
				delete(controller.Probes, i)
			}
		}

	}
}

func (controller *Controller) UpdateLevelSensors(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetLevelSenosrCount(); i++ {
		mode := profiluxController.GetLevelSensorMode(i)

		sensor, found := controller.LevelSensors[i]

		if mode != types.LevelSensorNotEnabled {
			if !found {
				sensor = models.NewLevelSensor(i)
				controller.LevelSensors[i] = sensor
			}

			sensor.Update(profiluxController)

		} else {
			if found {
				delete(controller.LevelSensors, i)
			}
		}

	}
}
