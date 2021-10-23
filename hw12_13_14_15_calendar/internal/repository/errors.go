package repository

import (
	"errors"
)

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrEventValidation   = errors.New("validation error")
	ErrEventCreate       = errors.New("error happened during creating the event")
	ErrEventCreateFailed = errors.New("event was not created")
	ErrEventUpdate       = errors.New("error happened during updating the event")
	ErrEventUpdateFailed = errors.New("event was not updated")
	ErrEventDelete       = errors.New("error happened during removing the event")
	ErrEventDeleteFailed = errors.New("event was not removed")
)
