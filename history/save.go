package history

import (
	"github.com/cjburchell/reefstatus-go/data"
	"time"
)

func SaveDay() error {
	now := time.Now()
	for _, probe := range data.Controller.GetProbes() {
		err := Save(Data{Type: probe.Id, Value: probe.Value, Time: now}, DayRange)
		if err != nil {
			return err
		}
	}

	return CleanUpDay()
}

func SaveWeek() error {

	now := time.Now()
	for _, probe := range data.Controller.GetProbes() {
		average, err := GetLastHourAvrage(probe.Id)
		if err != nil {
			return err
		}

		err = Save(Data{Type: probe.Id, Value: average, Time: now}, WeekRange)
		if err != nil {
			return err
		}
	}
	return CleanUpWeek()
}

func SaveYear() error {
	now := time.Now()
	for _, probe := range data.Controller.GetProbes() {
		average, err := GetLastDayAvrage(probe.Id)
		if err != nil {
			return err
		}

		err = Save(Data{Type: probe.Id, Value: average, Time: now}, YearRange)
		if err != nil {
			return err
		}
	}
	return CleanUpYear()
}
