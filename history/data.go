package history

import (
	"time"
)

type DataInterface interface {
	Setup() error
	SaveDayData(data Data) error
	SaveWeekData(data Data) error
	SaveYearData(data Data) error
	CleanUpDay() error
	CleanUpWeek() error
	CleanUpYear() error
	GetLastTimeYearDataWasSaved() (int64, error)
	GetLastTimeWeekDataWasSaved() (int64, error)
	GetDayDataPoints(dataType string) ([]Data, error)
	GetYearDataPoints(dataType string) ([]Data, error)
	GetWeekDataPoints(dataType string) ([]Data, error)
	GetDataPointsFromLastHour(dataType string) ([]Data, error)
}

var DataInstance DataInterface = &PostgresData{}

type Data struct {
	Time  time.Time `json:"time"`
	Type  string    `json:"type"`
	Value float64   `json:"value"`
}
