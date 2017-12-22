package models

import "fmt"

type ProgrammableLogic struct {
	DisplayName string
	Index       int
	Function    LogicFunction
	Input1      PortMode
	Input2      PortMode
}

func NewProgrammableLogic(index int) *ProgrammableLogic {
	var programmableLogic ProgrammableLogic
	programmableLogic.Index = index
	programmableLogic.DisplayName = fmt.Sprintf("Programable Logic %i", index+1)

	return &programmableLogic
}
