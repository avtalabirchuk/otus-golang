package repository

import (
	"github.com/go-playground/validator/v10"
)

func NewEventValidator() *validator.Validate {
	v := validator.New()
	v.RegisterStructValidation(EventStructLevelValidation, Event{})
	return v
}

func EventStructLevelValidation(sl validator.StructLevel) {
	event := sl.Current().Interface().(Event)

	if event.StartDate.Unix() == 0 {
		sl.ReportError(event.StartDate, "StartDate", "StartDate", "required", "")
	}

	if event.EndDate.Unix() == 0 {
		sl.ReportError(event.EndDate, "EndDate", "EndDate", "required", "")
	}
}
