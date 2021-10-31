package server

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

var (
	ErrInvalidNullTime = errors.New("field value is invalid")
	ErrObjectIsNil     = errors.New("object is nil")
)

func ConvertTimeToTimestamp(ntime sql.NullTime) (result *timestamp.Timestamp, err error) {
	if !ntime.Valid {
		err = ErrInvalidNullTime
		return
	}
	value, err := ntime.Value()
	if err != nil {
		return nil, err
	}
	tvalue, ok := value.(time.Time)
	if !ok {
		err = ErrInvalidNullTime
		return
	}
	return ptypes.TimestampProto(tvalue)
}

// Didn't want to use reflection.
func ConvertEventToProto(evt repository.Event) (*Event, error) {
	result := &Event{
		ID:          evt.ID,
		UserID:      evt.UserID,
		Title:       evt.Title,
		NotifiedFor: int64(evt.NotifiedFor),
	}
	if value, err := ptypes.TimestampProto(evt.StartDate); err != nil {
		return nil, err
	} else {
		result.StartDate = value
	}
	if value, err := ptypes.TimestampProto(evt.EndDate); err != nil {
		return nil, err
	} else {
		result.EndDate = value
	}
	return result, nil
}

// Didn't want to use reflection.
func ConvertEventFromProto(evt *Event) (*repository.Event, error) {
	if evt == nil {
		return nil, ErrObjectIsNil
	}
	result := repository.Event{}
	if evt.ID != 0 {
		result.ID = evt.ID
	}
	if evt.UserID != 0 {
		result.UserID = evt.UserID
	}
	if evt.NotifiedFor != 0 {
		result.NotifiedFor = int(evt.NotifiedFor)
	}
	if evt.Title != "" {
		result.Title = evt.Title
	}
	if evt.StartDate != nil {
		result.StartDate = time.Unix(evt.StartDate.GetSeconds(), int64(evt.StartDate.GetNanos()))
	}
	if evt.EndDate != nil {
		result.EndDate = time.Unix(evt.EndDate.GetSeconds(), int64(evt.EndDate.GetNanos()))
	}
	return &result, nil
}
