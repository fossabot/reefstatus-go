package settings

import (
	"github.com/cjburchell/reefstatus-go/common"
)

type AlertSettings struct {
	Enabled        bool
	SendOnReminder bool

	MailUserName string
	MailPassword string
	MailServer   string
	MailTo       string
	MailFrom     string
}

func newAlertSettings() (settings AlertSettings) {
	settings.Enabled = common.GetEnvBool("ALERT_ENABLE", false)
	settings.SendOnReminder = common.GetEnvBool("ALERT_REMINDER_ENABLE", false)
	settings.MailUserName = common.GetEnv("ALERT_MAIL_USERNAME", "reefstatusalert")
	settings.MailPassword = common.GetEnv("ALERT_MAIL_PASSWORD", "")
	settings.MailServer = common.GetEnv("ALERT_MAIL_SERVER", "smtp.gmail.com:587")
	settings.MailTo = common.GetEnv("ALERT_MAIL_To", "cjburchell@yahoo.com")
	settings.MailFrom = common.GetEnv("ALERT_MAIL_From", "reefstatusalert@gmail.com")
	return
}

var Alert = newAlertSettings()
