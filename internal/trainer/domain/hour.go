package domain

import "time"

type Availability int

const (
	Available Availability = iota
	NotAvailable
	TrainingScheduled
)

var allowedDaysToScheduleFromNow = 7

type Hour struct {
	Hour         time.Time
	availability Availability
}

func NewAvailableHour(hour time.Time) (*Hour, error) {
	if err := validateTime(hour); err != nil {
		return nil, err
	}

	return &Hour{
		Hour:         hour,
		availability: Available,
	}, nil
}

func validateTime(hour time.Time) error {
	if !hour.Round(time.Hour).Equal(hour) {
		return ErrNotFullHour
	}

	if hour.After(time.Now().AddDate(0, 0, allowedDaysToScheduleFromNow)) {
		return ErrTooDistantDate
	}

	currentHour := time.Now().Truncate(time.Hour)

	if hour.Before(currentHour) || hour.Equal(currentHour) {
		return ErrPastHour
	}

	return nil
}

func (h *Hour) ScheduleTraining() error {
	if h.availability != Available {
		return ErrHourNotAvailable
	}

	h.availability = TrainingScheduled
	return nil
}

func (h *Hour) GetAvailability() Availability {
	return h.availability
}
