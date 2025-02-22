package domain

import "errors"

var (
	ErrNotFullHour      = errors.New("the hour schedule must be a fullhour. e.g. 2011-10-05T14:00:00.000Z.")
	ErrTooDistantDate   = errors.New("the hour schedule has too distant date.")
	ErrPastHour         = errors.New("the hour schedule can't be in the past")
	ErrHourNotAvailable = errors.New("the selected hour is not available")
)
