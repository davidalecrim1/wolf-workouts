package command

import "time"

type ScheduleTrainingCommand struct {
	UserID           string
	Notes            string
	TrainingDateTime time.Time
}
