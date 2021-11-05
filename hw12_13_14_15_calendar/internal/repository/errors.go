package repository

import (
	"errors"
)

var (
	ErrConfigInvalid   = errors.New("provided config is invalid")
	ErrItemNotFound    = errors.New("item not found")
	ErrEventValidation = errors.New("validation error")
	ErrEventCreate     = errors.New("error happened during creating the event")
	ErrEventUpdate     = errors.New("error happened during updating the event")
	ErrEventDelete     = errors.New("error happened during removing the event")
	ErrUserCreate      = errors.New("error happened during creating the user")
)
