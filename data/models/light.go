package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type Light struct {
	BaseInfo
	Value          float64
	Channel        int
	IsDimmable     bool
	OperationHours int
	IsLightOn      bool
}

func NewLight(index int) *Light {
	var light Light
	light.Channel = index
	light.Type = "Light"
	light.Units = "%"
	light.Id = fmt.Sprintf("Light%d", 1+index)
	return &light
}

func (light *Light) Update(controller *profilux.Controller) {
	light.IsDimmable = controller.GetIsLightDimmable(light.Channel)
	light.OperationHours = controller.GetLightOperationHours(light.Channel)
	light.Value = controller.GetLightValue(light.Channel)
	light.IsLightOn = light.Value != 0
	light.DisplayName = controller.GetLightName(light.Channel)
}
