package models

import "github.com/cjburchell/reefstatus-go/profilux/types"

type LogicFunction struct {
	Invert1   bool
	Invert2   bool
	LogicMode types.LogicMode
}
