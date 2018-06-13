package history

import (
	"database/sql"
	"fmt"
	"github.com/cjburchell/reefstatus-go/common"
	_ "github.com/lib/pq"
)

type PostgresData struct {
	db *sql.DB
}

func (pgData *PostgresData) Setup() (err error) {

	dbUser := common.GetEnv("DB_USER", "admin")
	dbPassword := common.GetEnv("DB_PASSWORD", "admin")
	dbAddress := common.GetEnv("DB_ADDRESS", "localhost")
	dbName := common.GetEnv("DB_NAME", "reefstatus_history")

	dbInfo := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=verify-full", dbUser, dbPassword, dbAddress, dbName)
	pgData.db, err = sql.Open("postgres", dbInfo)
	return
}

func (PostgresData) SaveDayData(data Data) error {
	return nil
}
func (PostgresData) SaveWeekData(data Data) error {
	return nil
}
func (PostgresData) SaveYearData(data Data) error {
	return nil
}
func (PostgresData) CleanUpDay() error {
	return nil
}
func (PostgresData) CleanUpWeek() error {
	return nil
}
func (PostgresData) CleanUpYear() error {
	return nil
}
func (PostgresData) GetLastTimeYearDataWasSaved() (int64, error) {
	return 0, nil
}
func (PostgresData) GetLastTimeWeekDataWasSaved() (int64, error) {
	return 0, nil
}
func (PostgresData) GetDayDataPoints(dataType string) ([]Data, error) {
	return nil, nil
}
func (PostgresData) GetYearDataPoints(dataType string) ([]Data, error) {
	return nil, nil
}
func (PostgresData) GetWeekDataPoints(dataType string) ([]Data, error) {
	return nil, nil
}
func (PostgresData) GetDataPointsFromLastHour(dataType string) ([]Data, error) {
	return nil, nil
}
