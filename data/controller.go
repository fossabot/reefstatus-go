package main

import (
	"fmt"
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

func newController() *Controller {

	var controller Controller
	controller.Info = *models.NewInfo()
	controller.DigitalInputs = make(map[int]*models.DigitalInput)
	controller.DosingPumps = make(map[int]*models.DosingPump)
	controller.LevelSensors = make(map[int]*models.LevelSensor)
	controller.Lights = make(map[int]*models.Light)
	controller.LPorts = make(map[int]*models.LPort)
	controller.Probes = make(map[int]*models.Probe)
	controller.ProgrammableLogic = make(map[int]*models.ProgrammableLogic)
	controller.Pumps = make(map[int]*models.CurrentPump)
	controller.SPorts = make(map[int]*models.SPort)

	return &controller
}

func (controller *Controller) Update(profiluxController *profilux.Controller) {
	controller.Info.Update(profiluxController)
	controller.UpdateProbes(profiluxController)
	controller.UpdateLevelSensors(profiluxController)
	controller.UpdateDigitalInputs(profiluxController)
	controller.UpdateDosingPumps(profiluxController)
	controller.UpdateLights(profiluxController)
	controller.UpdateCurrentPumps(profiluxController)
	controller.UpdateProgrammableLogic(profiluxController)
	controller.UpdateSPorts(profiluxController)
	controller.UpdateLPorts(profiluxController)
	controller.UpdateAssociations()
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

func (controller *Controller) UpdateDigitalInputs(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetDigitalInputCount(); i++ {
		mode := profiluxController.GetDigitalInputFunction(i)

		sensor, found := controller.DigitalInputs[i]

		if mode != types.DigitalInputFunctionNotUsed {
			if !found {
				sensor = models.NewDigitalInput(i)
				controller.DigitalInputs[i] = sensor
			}

			sensor.Update(profiluxController)

		} else {
			if found {
				delete(controller.DigitalInputs, i)
			}
		}
	}
}

func (controller *Controller) UpdateDosingPumps(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetTimerCount(); i++ {
		pump, found := controller.DosingPumps[i]

		settings := profiluxController.GetTimerSettings(i)
		if settings.Mode == types.TimerModeAutoDosing {
			if !found {
				pump = models.NewDosingPump(i)
				controller.DosingPumps[i] = pump
			}

			pump.Update(profiluxController)

		} else {
			if found {
				delete(controller.DosingPumps, i)
			}
		}
	}
}

func (controller *Controller) UpdateLights(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetLightCount(); i++ {
		light, found := controller.Lights[i]
		if profiluxController.GetIsLightActive(i) {
			if !found {
				light = models.NewLight(i)
				controller.Lights[i] = light
			}

			light.Update(profiluxController)

		} else {
			if found {
				delete(controller.Lights, i)
			}
		}
	}
}

func (controller *Controller) UpdateCurrentPumps(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetCurrentPumpCount(); i++ {
		pump, found := controller.Pumps[i]
		if profiluxController.GetIsCurrentPumpAssigned(i) {
			if !found {
				pump = models.NewCurrentPump(i)
				controller.Pumps[i] = pump
			}

			pump.Update(profiluxController)

		} else {
			if found {
				delete(controller.Pumps, i)
			}
		}
	}
}

func (controller *Controller) UpdateProgrammableLogic(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetProgrammableLogicCount(); i++ {
		logic, found := controller.ProgrammableLogic[i]

		var input1 = profiluxController.GetProgramLogicInput(0, i)
		var input2 = profiluxController.GetProgramLogicInput(1, i)

		if input1.DeviceMode != types.DeviceModeAlwaysOff && input2.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				logic = models.NewProgrammableLogic(i)
				controller.ProgrammableLogic[i] = logic
			}

			logic.Update(profiluxController)

		} else {
			if found {
				delete(controller.ProgrammableLogic, i)
			}
		}
	}
}

func (controller *Controller) UpdateSPorts(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetSPortCount(); i++ {
		port, found := controller.SPorts[i]

		mode := profiluxController.GetSPortFunction(i)

		if mode.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				port = models.NewSPort(i)
				controller.SPorts[i] = port
			}

			port.Update(profiluxController)

		} else {
			if found {
				delete(controller.SPorts, i)
			}
		}
	}
}

func (controller *Controller) UpdateLPorts(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetLPortCount(); i++ {
		port, found := controller.LPorts[i]

		mode := profiluxController.GetLPortFunction(i)

		if mode.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				port = models.NewLPort(i)
				controller.LPorts[i] = port
			}

			port.Update(profiluxController)

		} else {
			if found {
				delete(controller.LPorts, i)
			}
		}
	}
}

func (controller Controller) getAssociatedModeItem(mode profilux.PortMode) string {
	index := mode.Port - 1
	if mode.IsProbe {
		probe, found := controller.Probes[index]
		if found {
			return probe.Id
		}
	}

	switch mode.DeviceMode {
	case types.DeviceModeLights:
		light, found := controller.Lights[index]
		if found {
			return light.Id
		}

	case types.DeviceModeTimer:
		timer, found := controller.DosingPumps[index]
		if found {
			return timer.Id
		}
	case types.DeviceModeWater:
		level, found := controller.LevelSensors[index]
		if found {
			return level.Id
		}

	case types.DeviceModeDrainWater:
		level, found := controller.LevelSensors[index]
		if found {
			return level.Id
		}

	case types.DeviceModeWaterChange:
		level, found := controller.LevelSensors[index]
		if found {
			return level.Id
		}

	case types.DeviceModeCurrentPump:
		pump, found := controller.Pumps[index]
		if found {
			return pump.Id
		}

	case types.DeviceModeProgrammableLogic:
		logic, found := controller.ProgrammableLogic[index]
		if found {
			return fmt.Sprintf("%d", logic.Index)
		}
	}

	return ""
}

func (controller *Controller) UpdateAssociations() {

	for _, logic := range controller.ProgrammableLogic {
		logic.Input1.Id = controller.getAssociatedModeItem(logic.Input1)
		logic.Input1.Id = controller.getAssociatedModeItem(logic.Input2)
	}

	for _, port := range controller.SPorts {
		port.Mode.Id = controller.getAssociatedModeItem(port.Mode)
	}

	for _, port := range controller.LPorts {
		port.Mode.Id = controller.getAssociatedModeItem(port.Mode)
	}
}
