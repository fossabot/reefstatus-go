package history

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/common"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/google/uuid"
	"github.com/rhinoman/couchdb-go"
	"net/url"
	"time"
)

type Data struct {
	ID    string    `json:"_id"`
	Rev   string    `json:"_rev,omitempty"`
	Time  time.Time `json:"time"`
	Type  string    `json:"type"`
	Value float64   `json:"value"`
}

type viewRow struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Doc   Data   `json:"doc"`
}

type viewResult struct {
	Total  int       `json:"total_rows"`
	Offset int       `json:"offset"`
	Rows   []viewRow `json:"rows"`
}

func getDatabase(databaseName string) (*couchdb.Database, error) {

	address := common.GetEnv("COUCHDB_ADDRESS", "localhost")
	port := common.GetEnvInt("COUCHDB_PORT", 5984)

	connection, err := couchdb.NewConnection(address, port, time.Second*30)
	if err != nil {
		return nil, err
	}

	auth := couchdb.BasicAuth{Username: "admin", Password: "admin"}

	/*list, err := connection.GetDBList()
	if err != nil {
		return nil, err
	}

	found := false
	for _, item := range list {
		if item == databaseName {
			found = true
			break
		}
	}

	if !found {
		err := connection.CreateDB(databaseName, &auth)
		if err != nil {
			return nil, err
		}

	}*/

	database := connection.SelectDB(databaseName, &auth)
	return database, nil

}

const (
	YearRange = "reefstatus_year"
	DayRange  = "reefstatus_day"
	WeekRange = "reefstatus_week"
)

func Save(data Data, databaseName string) error {
	database, err := getDatabase(databaseName)
	if err != nil {
		return err
	}

	data.ID = uuid.New().String()

	log.Debugf("Saving History %s ID: %s Value:%f", databaseName, data.Type, data.Value)
	_, err = database.Save(&data, data.ID, "")
	return err
}

func cleanUp(databaseName string, designDoc string, view string, cutoffDate time.Time) error {
	database, err := getDatabase(databaseName)
	if err != nil {
		return err
	}

	var data viewResult
	options := url.Values{}
	options.Add("include_docs", "true")
	err = database.GetView(designDoc, view, &data, &options)
	if err != nil {
		return err
	}

	itemsRemoved := false
	for _, item := range data.Rows {
		if item.Doc.Time.Before(cutoffDate) {
			log.Debugf("Removing %s %s %s", databaseName, item.Doc.Type, item.Doc.Time.String())
			itemsRemoved = true
			_, err := database.Delete(item.Doc.ID, item.Doc.Rev)
			if err != nil {
				return err
			}
		}
	}

	if itemsRemoved {
		database.Compact()
	}

	return nil
}

func CleanUpDay() error {
	cutoffDate := time.Now().Add(time.Hour * -24)
	return cleanUp(DayRange, "daylog", "active", cutoffDate)
}

func CleanUpWeek() error {
	cutoffDate := time.Now().AddDate(0, 0, -7)
	return cleanUp(WeekRange, "weeklog", "active", cutoffDate)
}

func CleanUpYear() error {
	cutoffDate := time.Now().AddDate(-1, 0, 0)
	return cleanUp(YearRange, "yearlog", "active", cutoffDate)
}

type lastSavedRow struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type lastSavedResult struct {
	Rows []lastSavedRow `json:"rows"`
}

func GetLastTimeYearDataWasSaved() (int64, error) {
	database, err := getDatabase(YearRange)
	if err != nil {
		return 0, err
	}

	var result lastSavedResult
	options := url.Values{}
	options.Add("reduce", "true")
	err = database.GetView("yearlog", "lastsavedtime", &result, &options)
	if err != nil {
		return 0, err
	}

	if len(result.Rows) != 0 {
		return result.Rows[0].Value, nil
	}

	return 0, nil
}

func GetLastTimeWeekDataWasSaved() (int64, error) {
	database, err := getDatabase(WeekRange)
	if err != nil {
		return 0, err
	}

	var result lastSavedResult
	options := url.Values{}
	options.Add("reduce", "true")
	err = database.GetView("weeklog", "lastsavedtime", &result, &options)
	if err != nil {
		return 0, err
	}

	if len(result.Rows) != 0 {
		return result.Rows[0].Value, nil
	}

	return 0, nil
}

func getDataView(dataType string, databaseName string, designDoc string, view string, filter func(item Data) bool) ([]Data, error) {
	database, err := getDatabase(databaseName)
	if err != nil {
		return nil, err
	}

	if filter == nil {
		filter = func(item Data) bool {
			return true
		}
	}

	var result viewResult
	options := url.Values{}
	options.Add("include_docs", "true")
	options.Add("key", fmt.Sprintf("\"%s\"", dataType))
	err = database.GetView(designDoc, view, &result, &options)
	if err != nil {
		return nil, err
	}

	data := make([]Data, len(result.Rows))
	for i, item := range result.Rows {
		if filter(item.Doc) {
			data[i] = item.Doc
		}
	}

	return data, nil
}

func GetDayDataPoints(dataType string) ([]Data, error) {
	return getDataView(dataType, DayRange, "daylog", "active", nil)
}

func GetYearDataPoints(dataType string) ([]Data, error) {
	return getDataView(dataType, YearRange, "yearlog", "active", nil)
}

func GetWeekDataPoints(dataType string) ([]Data, error) {
	log.Debugf("Get Week logs %s", dataType)
	return getDataView(dataType, WeekRange, "weeklog", "active", nil)
}

func GetDataPointsFromLastHour(dataType string) ([]Data, error) {
	cutoffTime := time.Now().Add(-time.Hour)
	return getDataView(dataType, DayRange, "daylog", "active", func(item Data) bool {
		return item.Time.After(cutoffTime)
	})
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

func GetLastHourAvrage(dataType string) (float64, error) {
	data, err := GetDataPointsFromLastHour(dataType)
	if err != nil {
		return 0, err
	}
	return average(data), nil
}

func GetLastDayAvrage(dataType string) (float64, error) {
	data, err := GetDayDataPoints(dataType)
	if err != nil {
		return 0, err
	}
	return average(data), nil
}
