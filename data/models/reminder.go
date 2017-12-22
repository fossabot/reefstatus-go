package models

import (
	"github.com/cjburchell/reefstatus-go/profilux"
	"time"
)

type Reminder struct {
	IsOverdue   bool
	Next        time.Time
	Text        string
	Index       int
	Period      int
	IsRepeating bool
}

func NewReminder(index int) *Reminder {
	var reminder Reminder
	reminder.Index = index
	return &reminder
}

func (reminder *Reminder) Update(controller *profilux.Controller) {
	reminder.Next = controller.GetReminderNext(reminder.Index)
	reminder.IsOverdue = reminder.Next.After(time.Now())
	reminder.Text = controller.GetReminderText(reminder.Index)
	reminder.Period = controller.GetReminderPeriod(reminder.Index)
	reminder.IsRepeating = controller.GetReminderIsRepeating(reminder.Index)
}
