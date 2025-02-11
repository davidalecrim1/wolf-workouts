package command

import "time"

type ScheduleTrainingCommand struct {
	UserID           string
	Username         string
	Notes            string
	TrainingDateTime time.Time
}
