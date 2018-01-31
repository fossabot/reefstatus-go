package main

import (
	"github.com/cjburchell/reefstatus-go/alert"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/communication"
	"github.com/cjburchell/reefstatus-go/data"
	"github.com/cjburchell/reefstatus-go/history"
	"strconv"
	"time"
)

const LogRate = time.Second * 30
const HourLogRate = time.Hour
const DayLogRate = time.Hour * 24

func UpdateController() {
	session := communication.NewSession()
	defer session.Close()
	data.Controller.Update()

	session.Publish(communication.UpdateMessage, []byte(strconv.FormatBool(true)))

	updateCount := 0
	for {
		log.Debugf("RefreshSettings Sleeping for %s", LogRate.String())
		<-time.After(LogRate)
		if updateCount%20 == 19 {
			data.Controller.Update()
		} else {
			data.Controller.UpdateState()
		}

		session.Publish(communication.UpdateMessage, []byte(strconv.FormatBool(false)))
		updateCount++
	}
}

func UpdateAlerts() {
	session := communication.NewSession()
	defer session.Close()
	ch := session.Subscribe(communication.UpdateMessage)
	for {
		<-ch
		alert.Check()
	}
}

func UpdateHistory() {
	session := communication.NewSession()
	defer session.Close()
	ch := session.Subscribe(communication.UpdateMessage)

	for {
		<-ch
		history.SaveDay()
	}
}

func UpdateWeekHistory() {
	session := communication.NewSession()
	defer session.Close()
	ch := session.Subscribe(communication.UpdateMessage)
	<-ch

	lastHourSavedTime, err := history.GetLastTimeWeekDataWasSaved()
	if err != nil {
		log.Error(err)
		return
	}

	timeSinceLastHourSaved := time.Duration(int64(time.Second) * (time.Now().Unix() - lastHourSavedTime/1000))
	duration := HourLogRate
	if timeSinceLastHourSaved > HourLogRate {
		err = history.SaveWeek()
		if err != nil {
			log.Error(err)
			return
		}
	} else if lastHourSavedTime != 0 {
		duration = HourLogRate - timeSinceLastHourSaved
	} else {
		err = history.SaveWeek()
		if err != nil {
			log.Error(err)
			return
		}
	}

	for {
		log.Debugf("SaveWeekHistory Sleeping for %s", duration.String())
		<-time.After(duration)
		err := history.SaveWeek()
		if err != nil {
			log.Error(err)
			return
		}
		duration = HourLogRate
	}
}

func UpdateYearHistory() {
	session := communication.NewSession()
	defer session.Close()
	ch := session.Subscribe(communication.UpdateMessage)
	<-ch

	lastHourSavedTime, err := history.GetLastTimeYearDataWasSaved()
	if err != nil {
		return
	}

	timeSinceLastHourSaved := time.Duration(int64(time.Millisecond) * (time.Now().Unix() - lastHourSavedTime/1000))
	duration := DayLogRate
	if timeSinceLastHourSaved > DayLogRate {
		err = history.SaveYear()
		if err != nil {
			log.Error(err)
			return
		}
	} else if lastHourSavedTime != 0 {
		duration = DayLogRate - timeSinceLastHourSaved
	} else {
		err = history.SaveYear()
		if err != nil {
			log.Error(err)
			return
		}
	}

	for {
		log.Debugf("SaveYearHistory Sleeping for %s", duration.String())
		<-time.After(duration)
		err := history.SaveYear()
		if err != nil {
			log.Error(err)
			return
		}
		duration = DayLogRate
	}
}
