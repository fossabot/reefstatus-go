package history

import (
	"github.com/cjburchell/reefstatus-go/data"
	"time"
)

func SaveDay() error {
	now := time.Now()
	for _, probe := range data.Controller.GetProbes() {
		err := DataInstance.SaveDayData(Data{Type: probe.Id, Value: probe.Value, Time: now})
		if err != nil {
			return err
		}
	}

	return DataInstance.CleanUpDay()
}

func SaveWeek() error {

	now := time.Now()
	for _, probe := range data.Controller.GetProbes() {
		average, err := getLastHourAverage(probe.Id)
		if err != nil {
			return err
		}

		err = DataInstance.SaveWeekData(Data{Type: probe.Id, Value: average, Time: now})
		if err != nil {
			return err
		}
	}
	return DataInstance.CleanUpWeek()
}

func SaveYear() error {
	now := time.Now()
	for _, probe := range data.Controller.GetProbes() {
		average, err := getLastDayAverage(probe.Id)
		if err != nil {
			return err
		}

		err = DataInstance.SaveYearData(Data{Type: probe.Id, Value: average, Time: now})
		if err != nil {
			return err
		}
	}
	return DataInstance.CleanUpYear()
}

func average(data []Data) float64 {
	if len(data) == 0 {
		return 0
	}

	sum := float64(0)
	for _, item := range data {
		sum += item.Value
	}

	return sum / float64(len(data))
}

func getLastHourAverage(dataType string) (float64, error) {
	data, err := DataInstance.GetDataPointsFromLastHour(dataType)
	if err != nil {
		return 0, err
	}
	return average(data), nil
}

func getLastDayAverage(dataType string) (float64, error) {
	data, err := DataInstance.GetDayDataPoints(dataType)
	if err != nil {
		return 0, err
	}
	return average(data), nil
}
