package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
)

var ErrUnSupportedRepoType = errors.New("unsupported repository type")

type dbConnector interface {
	Connect(context.Context, string) error
	Close() error
}

type Base interface {
	GetDayEvents(time.Time) ([]Event, error)
	GetWeekEvents(time.Time) ([]Event, error)
	GetMonthEvents(time.Time) ([]Event, error)

	CreateEvent(Event) (Event, error)
	UpdateEvent(int64, Event) (Event, error)
	DeleteEvent(int64) error
	dbConnector
}

type Stats interface {
	GetCurrentEvents() ([]Event, error)
	MarkEventsAsSent(*[]Event) error
	MarkEventsAsProcessing(*[]Event) error
	DeleteObsoleteEvents() error
	dbConnector
}

type Event struct {
	ID          int64          `db:"id"`
	UserID      int64          `db:"user_id" validate:"required"`
	Title       string         `db:"title" validate:"required"`
	Description sql.NullString `db:"description"`
	StartDate   time.Time      `db:"start_date" validate:"required"`
	EndDate     time.Time      `db:"end_date" validate:"required,gtfield=StartDate"`
	NotifiedFor int            `db:"notified_for" validate:"gte=1"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	UpdatedAt   sql.NullTime   `db:"updated_at"`
}

func newRepo(repoType string, args ...interface{}) interface{} {
	switch repoType {
	case "psql":
		return NewPSQLRepo(args...)
	case "memory":
		return NewMemoryRepo()
	}
	return nil
}

func New(repoType string, args ...interface{}) Base {
	repo, ok := newRepo(repoType, args...).(Base)
	if !ok {
		return nil
	}
	return repo
}

func NewStats(repoType string, args ...interface{}) Stats {
	repo, ok := newRepo(repoType, args...).(Stats)
	if !ok {
		return nil
	}
	return repo
}

func GetSQLDSN(c *config.DBConfig) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", c.Host, c.Port, c.DBName, c.User, c.Pass)
}
