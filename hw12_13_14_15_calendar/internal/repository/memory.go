package repository

import (
	"context"
	"sync"
	"time"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/utils"
)

type MemoryStorage map[int64]Event

type MemoryRepo struct {
	storage MemoryStorage
	mx      sync.RWMutex
}

func (r *MemoryRepo) Connect(ctx context.Context, url string) error {
	return nil
}

func (r *MemoryRepo) Init(ctx context.Context, url string) (err error) {
	return nil
}

func (r *MemoryRepo) Close() error {
	r.ClearStorage()
	return nil
}

func (r *MemoryRepo) ClearStorage() {
	r.storage = make(MemoryStorage)
}

func (r *MemoryRepo) GetStorage() *MemoryStorage {
	return &r.storage
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{storage: make(MemoryStorage)}
}

func isDateAfter(date time.Time, base time.Time) bool {
	return date.Equal(base) || date.After(base)
}

func isDateBefore(date time.Time, base time.Time) bool {
	return date.Equal(base) || date.Before(base)
}

func (r *MemoryRepo) getEventsInRange(startPeriod time.Time, endPeriod time.Time) ([]Event, error) {
	r.mx.RLock()
	defer r.mx.RUnlock()
	result := []Event{}
	for _, v := range r.storage {
		isStartPeriodInRange := isDateAfter(startPeriod, v.StartDate) && isDateBefore(startPeriod, v.EndDate)
		isEndPeriodInRange := isDateAfter(endPeriod, v.StartDate) && isDateBefore(endPeriod, v.EndDate)
		isPeriodOutRange := isDateAfter(endPeriod, v.EndDate) && isDateBefore(startPeriod, v.StartDate)
		if isStartPeriodInRange || isEndPeriodInRange || isPeriodOutRange {
			result = append(result, v)
		}
	}
	return result, nil
}

// Probably, need to display events only for particular user.
func (r *MemoryRepo) GetDayEvents(date time.Time) ([]Event, error) {
	return r.getEventsInRange(date, date.AddDate(0, 0, 1))
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

func (r *MemoryRepo) GetCurrentEvents() (result []Event, err error) {
	return r.GetDayEvents(time.Now())
}

func (r *MemoryRepo) MarkEventsAsSent(ids *[]int64) error {
	return nil
}

func (r *MemoryRepo) MarkEventsAsProcessed(ids *[]int64) error {
	return nil
}
