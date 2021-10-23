package repository

import (
	"context"
	"sync"
	"time"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/utils"
)

type MemoryRepo struct {
	storage map[int64]Event
	mx      sync.RWMutex
}

func (r *MemoryRepo) Connect(ctx context.Context, c *config.Config) error {
	return nil
}

func (r *MemoryRepo) Close() error {
	r.storage = nil
	return nil
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{storage: make(map[int64]Event)}
}

// Probably, need to display events only for particular user.
func (r *MemoryRepo) GetDayEvents(date time.Time) ([]Event, error) {
	r.mx.RLock()
	defer r.mx.RLock()
	result := []Event{}
	for _, v := range r.storage {
		if v.StartDate == date {
			result = append(result, v)
		}
	}
	return result, nil
}

func isDateInRange(date time.Time, startRange time.Time, endRange time.Time) bool {
	isAfterStart := date.Equal(startRange) || date.After(startRange)
	isBeforeEnd := date.Equal(endRange) || date.Before(endRange)
	return isAfterStart && isBeforeEnd
}

func (r *MemoryRepo) getEventsInRange(startPeriod time.Time, endPeriod time.Time) ([]Event, error) {
	r.mx.RLock()
	defer r.mx.RLock()
	result := []Event{}
	for _, v := range r.storage {
		isStartDateInRange := isDateInRange(v.StartDate, startPeriod, endPeriod)
		isEndDateInRange := isDateInRange(v.EndDate, startPeriod, endPeriod)
		if isStartDateInRange || isEndDateInRange {
			result = append(result, v)
		}
	}
	return result, nil
}

func (r *MemoryRepo) GetWeekEvents(date time.Time) ([]Event, error) {
	return r.getEventsInRange(date, date.AddDate(0, 0, 7))
}

func (r *MemoryRepo) GetMonthEvents(date time.Time) ([]Event, error) {
	return r.getEventsInRange(date, date.AddDate(0, 1, 0))
}

func (r *MemoryRepo) CreateEvent(data Event) (Event, error) {
	id, err := utils.GenerateUID()
	if err != nil {
		return Event{}, ErrEventCreate
	}
	data.ID = id
	r.mx.Lock()
	defer r.mx.Unlock()
	r.storage[id] = data
	return data, nil
}

func (r *MemoryRepo) UpdateEvent(id int64, data Event) (event Event, err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	event, ok := r.storage[id]
	if !ok {
		return Event{}, ErrEventNotFound
	}
	err = MergeEvents(&event, data)
	r.storage[id] = event
	return
}

func (r *MemoryRepo) DeleteEvent(id int64) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.storage[id]
	if !ok {
		return ErrEventNotFound
	}
	delete(r.storage, id)
	return
}
