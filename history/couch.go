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

const (
	YearRange = "reefstatus_year"
	DayRange  = "reefstatus_day"
	WeekRange = "reefstatus_week"
)

type CouchData struct {
	connection *couchdb.Connection
	auth       couchdb.BasicAuth
}

type viewRow struct {
	ID    string  `json:"id"`
	Key   string  `json:"key"`
	Value string  `json:"value"`
	Doc   DocData `json:"doc"`
}

type DocData struct {
	Data
	ID  string `json:"_id"`
	Rev string `json:"_rev,omitempty"`
}

type viewResult struct {
	Total  int       `json:"total_rows"`
	Offset int       `json:"offset"`
	Rows   []viewRow `json:"rows"`
}

func (couchData CouchData) getDatabase(databaseName string) (*couchdb.Database, error) {
	database := couchData.connection.SelectDB(databaseName, &couchData.auth)
	return database, nil
}

func (couchData *CouchData) Setup() error {
	address := common.GetEnv("COUCHDB_ADDRESS", "localhost")
	port := common.GetEnvInt("COUCHDB_PORT", 5984)

	var err error
	couchData.connection, err = couchdb.NewConnection(address, port, time.Second*30)
	if err != nil {
		return err
	}

	couchData.auth = couchdb.BasicAuth{Username: "admin", Password: "admin"}
	return nil
}

func (couchData CouchData) cleanUp(databaseName string, designDoc string, view string, cutoffDate time.Time) error {
	database, err := couchData.getDatabase(databaseName)
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

func (couchData CouchData) SaveDayData(data Data) error {
	return couchData.save(data, DayRange)
}

func (couchData CouchData) SaveWeekData(data Data) error {
	return couchData.save(data, WeekRange)
}

func (couchData CouchData) SaveYearData(data Data) error {
	return couchData.save(data, YearRange)
}

func (couchData CouchData) save(data Data, databaseName string) error {
	database, err := couchData.getDatabase(databaseName)
	if err != nil {
		return err
	}

	docData := DocData{
		Data: data,
		ID:   uuid.New().String(),
	}

	log.Debugf("Saving History %s ID: %s Value:%f", databaseName, data.Type, data.Value)
	_, err = database.Save(&data, docData.ID, "")
	return err
}

func (couchData CouchData) getDataView(dataType string, databaseName string, designDoc string, view string, filter func(item Data) bool) ([]Data, error) {
	database, err := couchData.getDatabase(databaseName)
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
		if filter(item.Doc.Data) {
			data[i] = item.Doc.Data
		}
	}

	return data, nil
}

func (couchData CouchData) CleanUpDay() error {
	cutoffDate := time.Now().Add(time.Hour * -24)
	return couchData.cleanUp(DayRange, "daylog", "active", cutoffDate)
}

func (couchData CouchData) CleanUpWeek() error {
	cutoffDate := time.Now().AddDate(0, 0, -7)
	return couchData.cleanUp(WeekRange, "weeklog", "active", cutoffDate)
}

func (couchData CouchData) CleanUpYear() error {
	cutoffDate := time.Now().AddDate(-1, 0, 0)
	return couchData.cleanUp(YearRange, "yearlog", "active", cutoffDate)
}

type lastSavedRow struct {
	Key   string `json:"key"`
	Value int64  `json:"value"`
}

type lastSavedResult struct {
	Rows []lastSavedRow `json:"rows"`
}

func (couchData CouchData) GetLastTimeYearDataWasSaved() (int64, error) {
	database, err := couchData.getDatabase(YearRange)
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

func (couchData CouchData) GetLastTimeWeekDataWasSaved() (int64, error) {
	database, err := couchData.getDatabase(WeekRange)
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

func (couchData CouchData) GetDayDataPoints(dataType string) ([]Data, error) {
	return couchData.getDataView(dataType, DayRange, "daylog", "active", nil)
}

func (couchData CouchData) GetYearDataPoints(dataType string) ([]Data, error) {
	return couchData.getDataView(dataType, YearRange, "yearlog", "active", nil)
}

func (couchData CouchData) GetWeekDataPoints(dataType string) ([]Data, error) {
	log.Debugf("Get Week logs %s", dataType)
	return couchData.getDataView(dataType, WeekRange, "weeklog", "active", nil)
}

func (couchData CouchData) GetDataPointsFromLastHour(dataType string) ([]Data, error) {
	cutoffTime := time.Now().Add(-time.Hour)
	return couchData.getDataView(dataType, DayRange, "daylog", "active", func(item Data) bool {
		return item.Time.After(cutoffTime)
	})
}
