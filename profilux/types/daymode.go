package types

type DayMode string

const (
	DayModeDaysOfWeek  = "Days of the Week"
	DayModeDayInterval = "Interval"
)

var dayModeMap = map[int]string{
	0: DayModeDaysOfWeek,
	1: DayModeDayInterval,
}

func GetDayMode(value int) string {
	if val, ok := dayModeMap[value]; ok {
		return val
	}

	return "Unknown"
}
