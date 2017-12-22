package types

type LogicMode string

const (
	LogicModeAnd            = "And"
	LogicModeOr             = "Or"
	LogicModeInvertAnd      = "Invert And"
	LogicModeInvertOr       = "Invert Or"
	LogicModeInverted       = "Inverted"
	LogicModeEqual          = "Equal"
	LogicModeNoEqual        = "Not Equal"
	LogicModePulse          = "Pulse"
	LogicModeDelayedOn      = "Delayed On"
	LogicModeDelayedOff     = "Delayed Off"
	LogicModeFrequentPulses = "Frequent Pulses"
	LogicModeSRFlipFlop     = "SR Flip-Flop"
	LogicModeExclusiveOr    = "Exclusive Or"
)

var logicModeMap = map[int]string{
	0:  LogicModeAnd,
	1:  LogicModeOr,
	2:  LogicModeInvertAnd,
	3:  LogicModeInvertOr,
	4:  LogicModeInverted,
	5:  LogicModeEqual,
	6:  LogicModeNoEqual,
	7:  LogicModePulse,
	8:  LogicModeDelayedOn,
	9:  LogicModeDelayedOff,
	10: LogicModeFrequentPulses,
	11: LogicModeSRFlipFlop,
	12: LogicModeExclusiveOr,
}

func GetLogicMode(value int) string {
	if val, ok := logicModeMap[value]; ok {
		return val
	}

	return "Unknown"
}
