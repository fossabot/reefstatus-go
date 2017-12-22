package types

type TimerMode string

const (
	TimerModeNormal       = "Normal"
	TimerModeShort        = "Short"
	TimerModeAutoDosing   = "AutoDosing"
	TimerModeManualDosing = "Manual Dosing"
	TimerModeStartEvent   = "Start Event"
	TimerModeCyclic       = "Cyclic"
)

var timerModeMap = map[int]string{
	0: TimerModeNormal,
	1: TimerModeShort,
	2: TimerModeAutoDosing,
	3: TimerModeManualDosing,
	4: TimerModeStartEvent,
	5: TimerModeCyclic,
}

func GetTimerMode(value int) string {
	if val, ok := timerModeMap[value]; ok {
		return val
	}

	return "Unknown"
}
