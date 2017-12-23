package models

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/profilux"
)

type ProgrammableLogic struct {
	DisplayName string
	Index       int
	Function    profilux.LogicFunction
	Input1      profilux.PortMode
	Input2      profilux.PortMode
}

func NewProgrammableLogic(index int) *ProgrammableLogic {
	var programmableLogic ProgrammableLogic
	programmableLogic.Index = index
	programmableLogic.DisplayName = fmt.Sprintf("Programable Logic %d", index+1)
	return &programmableLogic
}

func (logic *ProgrammableLogic) Update(controller *profilux.Controller) {
	logic.Input1 = controller.GetProgramLogicInput(0, logic.Index)
	logic.Input2 = controller.GetProgramLogicInput(1, logic.Index)
	logic.Function = controller.GetProgramLogicFunction(logic.Index)
}
