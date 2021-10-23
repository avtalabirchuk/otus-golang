package calendar

import (
	"errors"
)

var (
	ErrDateBusy        = errors.New("another event starts on the same date")
	ErrStartIsAfterEnd = errors.New("start date cannot be after end date")
)
