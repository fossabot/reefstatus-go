package data

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/data/models"
	"github.com/cjburchell/reefstatus-go/profilux"
	"github.com/cjburchell/reefstatus-go/profilux/types"
)

type IController interface {
	Update()
	UpdateState()
	GetInfo() models.Info
	GetDigitalInputs() []models.DigitalInput
	GetDosingPumps() []models.DosingPump
	GetLevelSensors() []models.LevelSensor
	GetLights() []models.Light
	GetLPorts() []models.LPort
	GetProbes() []models.Probe
	GetProgrammableLogic() []models.ProgrammableLogic
	GetCurrentPumps() []models.CurrentPump
	GetSPorts() []models.SPort

	FeedPause(enable bool)
	Thunderstorm(duration int)
	ResetReminder(index int)
	Maintenance(index int, enable bool)
	ClearLevelAlarm(id string)
	WaterChange(id string)
}

type controller struct {
	info              models.Info
	digitalInputs     map[int]*models.DigitalInput
	dosingPumps       map[int]*models.DosingPump
	levelSensors      map[int]*models.LevelSensor
	lights            map[int]*models.Light
	lPorts            map[int]*models.LPort
	probes            map[int]*models.Probe
	programmableLogic map[int]*models.ProgrammableLogic
	pumps             map[int]*models.CurrentPump
	sPorts            map[int]*models.SPort
}

func newController() IController {

	var controller controller
	controller.info = *models.NewInfo()
	controller.digitalInputs = make(map[int]*models.DigitalInput)
	controller.dosingPumps = make(map[int]*models.DosingPump)
	controller.levelSensors = make(map[int]*models.LevelSensor)
	controller.lights = make(map[int]*models.Light)
	controller.lPorts = make(map[int]*models.LPort)
	controller.probes = make(map[int]*models.Probe)
	controller.programmableLogic = make(map[int]*models.ProgrammableLogic)
	controller.pumps = make(map[int]*models.CurrentPump)
	controller.sPorts = make(map[int]*models.SPort)

	return &controller
}

var Controller = newController()

func (controller *controller) FeedPause(bool bool) {
	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	profiluxController.FeedPause(0, bool)
	controller.info.UpdateState(profiluxController)
}

func (controller *controller) Thunderstorm(duration int) {
	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	profiluxController.Thunderstorm(duration)
	controller.info.UpdateState(profiluxController)
}
func (controller *controller) ResetReminder(index int) {

	var reminder *models.Reminder
	for _, item := range controller.info.Reminders {
		if item.Index == index {
			reminder = item
			break
		}
	}

	if reminder == nil {
		log.Warnf("unable to find reminder")
		return
	}

	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
		return
	}

	defer profiluxController.Disconnect()

	if reminder.IsRepeating {
		profiluxController.ResetReminder(index, reminder.Period)
	} else {
		profiluxController.ClearReminder(index)
	}

	reminder.UpdateState(profiluxController)
}

func (controller *controller) Maintenance(index int, enable bool) {
	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	var maintenance *models.Maintenance
	for _, item := range controller.info.Maintenance {
		if item.Index == index {
			maintenance = item
			break
		}
	}

	if maintenance == nil {
		log.Warnf("unable to find reminder")
		return
	}

	profiluxController.Maintenance(enable, index)
	controller.info.UpdateState(profiluxController)
}

func (controller *controller) ClearLevelAlarm(id string) {

	var sensor *models.LevelSensor
	for _, level := range controller.levelSensors {
		if level.Id == id {
			sensor = level
			break
		}
	}

	if sensor == nil {
		log.Warnf("unable to find level snsoro %s", id)
		return
	}

	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	profiluxController.ClearLevelAlarm(sensor.Index)
	sensor.UpdateState(profiluxController)
	controller.info.UpdateState(profiluxController)
}
func (controller *controller) WaterChange(id string) {
	var sensor *models.LevelSensor
	for _, level := range controller.levelSensors {
		if level.Id == id {
			sensor = level
			break
		}
	}

	if sensor == nil {
		log.Warnf("unable to find level snsoro %s", id)
		return
	}

	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	profiluxController.WaterChange(sensor.Index)
	sensor.UpdateState(profiluxController)
}

func (controller *controller) Update() {

	log.Debug("RefreshSettings - Start")
	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	profiluxController.ResetStats()
	controller.info.Update(profiluxController)
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
	log.Debugf("Call Count %d", profiluxController.GetCallCount())
	log.Debug("RefreshSettings - End")
}

func (controller *controller) UpdateState() {
	log.Debug("RefreshSettings - Start")
	profiluxController, err := profilux.NewController()
	if err != nil {
		log.Errorf(err, "unable to connect")
	}

	defer profiluxController.Disconnect()

	profiluxController.ResetStats()
	controller.info.UpdateState(profiluxController)
	for _, item := range controller.probes {
		item.UpdateState(profiluxController)
	}

	for _, item := range controller.levelSensors {
		item.UpdateState(profiluxController)
	}

	for _, item := range controller.digitalInputs {
		item.UpdateState(profiluxController)
	}

	for _, item := range controller.lights {
		item.UpdateState(profiluxController)
	}

	for _, item := range controller.pumps {
		item.Update(profiluxController)
	}

	for _, item := range controller.sPorts {
		item.UpdateState(profiluxController)
	}

	for _, item := range controller.lPorts {
		item.UpdateState(profiluxController)
	}

	log.Debugf("Call Count %d", profiluxController.GetCallCount())
	log.Debug("RefreshSettings - End")
}

func (controller *controller) UpdateProbes(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetSensorCount(); i++ {
		sensorType := profiluxController.GetSensorType(i)

		var mode types.SensorMode
		active := false
		if sensorType != types.SensorTypeFree && sensorType != types.SensorTypeNone {
			mode = profiluxController.GetSensorMode(i)
			active = profiluxController.GetSensorActive(i)
		}

		probe, found := controller.probes[i]

		if active && mode == types.SensorModeNormal {
			if !found {
				probe = models.NewProbe(i)
				controller.probes[i] = probe
			}

			probe.Update(profiluxController)

		} else {
			if found {
				delete(controller.probes, i)
			}
		}

	}
}

func (controller *controller) UpdateLevelSensors(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetLevelSenosrCount(); i++ {
		mode := profiluxController.GetLevelSensorMode(i)

		sensor, found := controller.levelSensors[i]

		if mode != types.LevelSensorNotEnabled {
			if !found {
				sensor = models.NewLevelSensor(i)
				controller.levelSensors[i] = sensor
			}

			sensor.Update(profiluxController)

		} else {
			if found {
				delete(controller.levelSensors, i)
			}
		}

	}
}

func (controller *controller) UpdateDigitalInputs(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetDigitalInputCount(); i++ {
		mode := profiluxController.GetDigitalInputFunction(i)

		sensor, found := controller.digitalInputs[i]

		if mode != types.DigitalInputFunctionNotUsed {
			if !found {
				sensor = models.NewDigitalInput(i)
				controller.digitalInputs[i] = sensor
			}

			sensor.Update(profiluxController)

		} else {
			if found {
				delete(controller.digitalInputs, i)
			}
		}
	}
}

func (controller *controller) UpdateDosingPumps(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetTimerCount(); i++ {
		pump, found := controller.dosingPumps[i]

		settings := profiluxController.GetTimerSettings(i)
		if settings.Mode == types.TimerModeAutoDosing {
			if !found {
				pump = models.NewDosingPump(i)
				controller.dosingPumps[i] = pump
			}

			pump.Update(profiluxController)

		} else {
			if found {
				delete(controller.dosingPumps, i)
			}
		}
	}
}

func (controller *controller) UpdateLights(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetLightCount(); i++ {
		light, found := controller.lights[i]
		if profiluxController.GetIsLightActive(i) {
			if !found {
				light = models.NewLight(i)
				controller.lights[i] = light
			}

			light.Update(profiluxController)

		} else {
			if found {
				delete(controller.lights, i)
			}
		}
	}
}

func (controller *controller) UpdateCurrentPumps(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetCurrentPumpCount(); i++ {
		pump, found := controller.pumps[i]
		if profiluxController.GetIsCurrentPumpAssigned(i) {
			if !found {
				pump = models.NewCurrentPump(i)
				controller.pumps[i] = pump
			}

			pump.Update(profiluxController)

		} else {
			if found {
				delete(controller.pumps, i)
			}
		}
	}
}

func (controller *controller) UpdateProgrammableLogic(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetProgrammableLogicCount(); i++ {
		logic, found := controller.programmableLogic[i]

		var input1 = profiluxController.GetProgramLogicInput(0, i)
		var input2 = profiluxController.GetProgramLogicInput(1, i)

		if input1.DeviceMode != types.DeviceModeAlwaysOff && input2.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				logic = models.NewProgrammableLogic(i)
				controller.programmableLogic[i] = logic
			}

			logic.Update(profiluxController)

		} else {
			if found {
				delete(controller.programmableLogic, i)
			}
		}
	}
}

func (controller *controller) UpdateSPorts(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetSPortCount(); i++ {
		port, found := controller.sPorts[i]

		mode := profiluxController.GetSPortFunction(i)

		if mode.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				port = models.NewSPort(i)
				controller.sPorts[i] = port
			}

			port.Update(profiluxController)

		} else {
			if found {
				delete(controller.sPorts, i)
			}
		}
	}
}

func (controller *controller) UpdateLPorts(profiluxController *profilux.Controller) {
	for i := 0; i < profiluxController.GetLPortCount(); i++ {
		port, found := controller.lPorts[i]

		mode := profiluxController.GetLPortFunction(i)

		if mode.DeviceMode != types.DeviceModeAlwaysOff {
			if !found {
				port = models.NewLPort(i)
				controller.lPorts[i] = port
			}

			port.Update(profiluxController)

		} else {
			if found {
				delete(controller.lPorts, i)
			}
		}
	}
}

func (controller controller) getAssociatedModeItem(mode profilux.PortMode) string {
	index := mode.Port - 1
	if mode.IsProbe {
		probe, found := controller.probes[index]
		if found {
			return probe.Id
		}
	}

	switch mode.DeviceMode {
	case types.DeviceModeLights:
		light, found := controller.lights[index]
		if found {
			return light.Id
		}

	case types.DeviceModeTimer:
		timer, found := controller.dosingPumps[index]
		if found {
			return timer.Id
		}
	case types.DeviceModeWater:
		level, found := controller.levelSensors[index]
		if found {
			return level.Id
		}

	case types.DeviceModeDrainWater:
		level, found := controller.levelSensors[index]
		if found {
			return level.Id
		}

	case types.DeviceModeWaterChange:
		level, found := controller.levelSensors[index]
		if found {
			return level.Id
		}

	case types.DeviceModeCurrentPump:
		pump, found := controller.pumps[index]
		if found {
			return pump.Id
		}

	case types.DeviceModeProgrammableLogic:
		logic, found := controller.programmableLogic[index]
		if found {
			return fmt.Sprintf("%d", logic.Index)
		}
	}

	return ""
}

func (controller *controller) UpdateAssociations() {

	for _, logic := range controller.programmableLogic {
		logic.Input1.Id = controller.getAssociatedModeItem(logic.Input1)
		logic.Input1.Id = controller.getAssociatedModeItem(logic.Input2)
	}

	for _, port := range controller.sPorts {
		port.Mode.Id = controller.getAssociatedModeItem(port.Mode)
	}

	for _, port := range controller.lPorts {
		port.Mode.Id = controller.getAssociatedModeItem(port.Mode)
	}
}

func (controller *controller) GetInfo() models.Info {
	return controller.info
}

func (controller *controller) GetDigitalInputs() []models.DigitalInput {
	v := make([]models.DigitalInput, len(controller.digitalInputs))
	idx := 0
	for _, value := range controller.digitalInputs {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetDosingPumps() []models.DosingPump {
	v := make([]models.DosingPump, len(controller.dosingPumps))
	idx := 0
	for _, value := range controller.dosingPumps {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetLevelSensors() []models.LevelSensor {
	v := make([]models.LevelSensor, len(controller.levelSensors))
	idx := 0
	for _, value := range controller.levelSensors {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetLights() []models.Light {
	v := make([]models.Light, len(controller.lights))
	idx := 0
	for _, value := range controller.lights {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetLPorts() []models.LPort {
	v := make([]models.LPort, len(controller.lPorts))
	idx := 0
	for _, value := range controller.lPorts {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetProbes() []models.Probe {
	v := make([]models.Probe, len(controller.probes))
	idx := 0
	for _, value := range controller.probes {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetProgrammableLogic() []models.ProgrammableLogic {
	v := make([]models.ProgrammableLogic, len(controller.programmableLogic))
	idx := 0
	for _, value := range controller.programmableLogic {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetCurrentPumps() []models.CurrentPump {
	v := make([]models.CurrentPump, len(controller.pumps))
	idx := 0
	for _, value := range controller.pumps {
		v[idx] = *value
		idx++
	}

	return v
}
func (controller *controller) GetSPorts() []models.SPort {
	v := make([]models.SPort, len(controller.sPorts))
	idx := 0
	for _, value := range controller.sPorts {
		v[idx] = *value
		idx++
	}

	return v
}
