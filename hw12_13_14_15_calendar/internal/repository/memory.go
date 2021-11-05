package repository

import (
	"context"
	"sync"
	"time"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/utils"
)

type (
	MemoryStorageEvents map[int64]Event
	MemoryStorageUsers  map[int64]User
)

type MemoryRepo struct {
	events MemoryStorageEvents
	users  MemoryStorageUsers
	mx     sync.RWMutex
}

func (r *MemoryRepo) Connect(ctx context.Context, url string) error {
	return nil
}

func (r *MemoryRepo) Init(ctx context.Context, url, migrationsDir string) (err error) {
	return nil
}

func (r *MemoryRepo) Close() error {
	r.ClearStorage()
	return nil
}

func (r *MemoryRepo) ClearStorage() {
	r.events = make(MemoryStorageEvents)
	r.users = make(MemoryStorageUsers)
}

func (r *MemoryRepo) GetStorageEvents() *MemoryStorageEvents {
	return &r.events
}

func (r *MemoryRepo) GetStorageUsers() *MemoryStorageUsers {
	return &r.users
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{
		events: make(MemoryStorageEvents),
		users:  make(MemoryStorageUsers),
	}
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
	for _, v := range r.events {
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
	r.events[id] = data
	return data, nil
}

func (r *MemoryRepo) UpdateEvent(id int64, data Event) (event Event, err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	event, ok := r.events[id]
	if !ok {
		return Event{}, ErrItemNotFound
	}
	err = MergeEvents(&event, data)
	r.events[id] = event
	return
}

func (r *MemoryRepo) DeleteEvent(id int64) (err error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	_, ok := r.events[id]
	if !ok {
		return ErrItemNotFound
	}
	delete(r.events, id)
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

func (r *MemoryRepo) CreateUser(data User) (User, error) {
	id, err := utils.GenerateUID()
	if err != nil {
		return User{}, ErrUserCreate
	}
	data.ID = id
	r.mx.Lock()
	defer r.mx.Unlock()
	r.users[id] = data
	return data, nil
}

func (r *MemoryRepo) GetUser(id int64) (User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()
	user, ok := r.users[id]
	if !ok {
		return User{}, ErrItemNotFound
	}
	return user, nil
}
