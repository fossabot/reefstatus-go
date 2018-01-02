package alert

import (
	"fmt"
	"github.com/cjburchell/reefstatus-go/alert/settings"
	"github.com/cjburchell/reefstatus-go/common/log"
	"github.com/cjburchell/reefstatus-go/data"
	"github.com/cjburchell/reefstatus-go/profilux/types"
	"net"
	"net/mail"
	"net/smtp"
	"time"
)

func Check() {
	if !settings.Alert.Enabled {
		return
	}

	checkReminders()
	checkAlarm()
}

func checkReminders() {
	if !settings.Alert.SendOnReminder {
		return
	}

	for _, reminder := range data.Controller.GetInfo().Reminders {
		if reminder.IsOverdue {
			if reminderSent[reminder.Index] {
				continue
			}

			var subject = fmt.Sprintf("Reef Status [Reminder] %s", reminder.Text)
			var body = fmt.Sprintf("Reminder \"%s\" is now overdue", reminder.Text)
			go sendEmail(subject, body)

			reminderSent[reminder.Index] = true
		} else {
			if reminderSent[reminder.Index] {
				reminderSent[reminder.Index] = false
			}
		}
	}
}

var sentAlarmEmail = false
var reminderSent = make(map[int]bool)

func checkAlarm() {
	if data.Controller.GetInfo().Alarm == types.CurrentStateOn {
		if !sentAlarmEmail {
			sentAlarmEmail = true
			sendAlarmEmail(false)
		}
	} else {
		if sentAlarmEmail {
			sentAlarmEmail = false
			sendAlarmEmail(true)
		}
	}
}

func sendAlarmEmail(cleared bool) {
	statusTable := createStatusTable()
	if !cleared {
		go sendEmail("Reef Status [Alarm]", fmt.Sprintf("Alarm Detected %s\nReasons\n%s\n%s", time.Now().String(), findAlarmReason(), statusTable))
	} else {
		go sendEmail("Reef Status [Alarm Cleared]", fmt.Sprintf("Alarm Cleared %s\n%s", time.Now().String(), statusTable))
	}
}

func findAlarmReason() (reason string) {
	for _, probe := range data.Controller.GetProbes() {
		if probe.AlarmState == types.CurrentStateOn && probe.AlarmEnable {
			if probe.Value > probe.NominalValue+probe.AlarmDeviation {
				reason += fmt.Sprintf("%s is too high\n", probe.DisplayName)
			} else if probe.Value < probe.NominalValue-probe.AlarmDeviation {
				reason += fmt.Sprintf("%s is too low\n", probe.DisplayName)
			} else {
				reason += fmt.Sprintf("Alarm on %s\n", probe.DisplayName)
			}
		}
	}

	for _, sensor := range data.Controller.GetLevelSensors() {
		if sensor.AlarmState != types.CurrentStateOn {
			continue
		}

		reason += fmt.Sprintf("Level Timeout %s\n", sensor.DisplayName)
	}

	if len(reason) == 0 {
		reason = "Unknown"
	}

	return
}

func createStatusTable() (table string) {
	for _, probe := range data.Controller.GetProbes() {
		table += fmt.Sprintf("%s\t%f%s\n", probe.DisplayName, probe.ConvertedValue, probe.Units)
	}

	for _, sensor := range data.Controller.GetLevelSensors() {
		table += fmt.Sprintf("%s\t%s\t%s\n", sensor.DisplayName, sensor.OperationMode, sensor.Value)
	}

	return
}

func sendEmail(subject, body string) {
	log.Warnf("Email Error %s: %s", subject, body)
	from := mail.Address{Name: "Reef Status", Address: settings.Alert.MailFrom}
	to := mail.Address{Name: "", Address: settings.Alert.MailTo}
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	serverName := settings.Alert.MailServer
	host, _, _ := net.SplitHostPort(serverName)

	auth := smtp.PlainAuth("", settings.Alert.MailUserName, settings.Alert.MailPassword, host)

	err := smtp.SendMail(serverName, auth, from.Address, []string{to.Address}, []byte(message))
	if err != nil {
		log.Error(err)
		return
	}
}
