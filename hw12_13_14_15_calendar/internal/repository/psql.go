package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	// is used for init postgres.
	_ "github.com/lib/pq"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
)

var ErrDBOpen = errors.New("database open error")

type PSQLRepo struct {
	db            *sqlx.DB
	itemsPerQuery int
}

var insertQs = `INSERT INTO events
(user_id, title, description, start_date, start_time, end_date, end_time, notified_at)
VALUES
(:user_id, :title, :description, :start_date, :start_time, :end_date, :end_time, :notified_at)
RETURNING id`

var updateQs = `UPDATE events
	SET (title, description, start_date, start_time, end_date, end_time, notified_at) = (:title, :description, :start_date, :start_time, :end_date, :end_time, :notified_at)
	WHERE id = :id
	RETURNING id`

var deleteQs = `DELETE FROM events where id = $1 RETURNING id`

func getDSN(c *config.Config) string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", c.DBHost, c.DBPort, c.DBName, c.DBUser, c.DBPass)
}

func (r *PSQLRepo) Connect(ctx context.Context, c *config.Config) (err error) {
	r.db, err = sqlx.Connect("postgres", getDSN(c))
	if err != nil {
		return fmt.Errorf("%s: %w", ErrDBOpen, err)
	}
	if c.DBMaxConn != 0 {
		r.db.SetMaxOpenConns(c.DBMaxConn)
	}
	if c.DBItemsPerQuery != 0 {
		r.itemsPerQuery = c.DBItemsPerQuery
	}
	return r.db.PingContext(ctx)
}

func (r *PSQLRepo) Close() error {
	return r.db.Close()
}

func NewPSQLRepo() *PSQLRepo {
	return &PSQLRepo{itemsPerQuery: 100}
}

func (r *PSQLRepo) getEventsBetween(startDate time.Time, endDate time.Time) (result []Event, err error) {
	query := "SELECT * FROM events WHERE start_date >= $1 or end_date <= $2 ORDER BY start_date ASC LIMIT $3"
	err = r.db.Select(&result, query, startDate, endDate, r.itemsPerQuery)
	return
}

func (r *PSQLRepo) getEventByID(id int64) (result Event, err error) {
	err = r.db.Get(&result, "SELECT * FROM events WHERE id = $1", id)
	return
}

func (r *PSQLRepo) GetDayEvents(date time.Time) (result []Event, err error) {
	query := "SELECT * FROM events WHERE start_date = $1 ORDER BY start_date ASC LIMIT $2"
	err = r.db.Select(&result, query, date, r.itemsPerQuery)
	return
}

func (r *PSQLRepo) GetWeekEvents(date time.Time) (result []Event, err error) {
	return r.getEventsBetween(date, date.AddDate(0, 0, 7))
}

func (r *PSQLRepo) GetMonthEvents(date time.Time) (result []Event, err error) {
	return r.getEventsBetween(date, date.AddDate(0, 1, 0))
}

func (r *PSQLRepo) DeleteEvent(id int64) (err error) {
	result, err := r.db.Exec(deleteQs, id)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrEventDelete, err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", ErrEventDelete, err)
	}
	if n == 0 {
		err = ErrEventNotFound
	}
	return
}

func (r *PSQLRepo) CreateEvent(data Event) (event Event, err error) {
	stmt, err := r.db.PrepareNamed(insertQs)
	if err != nil {
		err = fmt.Errorf("%s: %w", ErrEventCreate, err)
		return
	}
	row := stmt.QueryRow(data)
	if err = row.Err(); err != nil {
		err = fmt.Errorf("%s: %w", ErrEventCreate, err)
		return
	}
	evt := Event{}
	err = row.StructScan(&evt)
	if err != nil {
		err = fmt.Errorf("%s: %w", ErrEventCreate, err)
		return
	}
	if evt.ID == 0 {
		err = ErrEventNotFound
		return
	}
	return r.getEventByID(evt.ID)
}

func (r *PSQLRepo) UpdateEvent(id int64, data Event) (event Event, err error) {
	data.ID = id
	stmt, err := r.db.PrepareNamed(updateQs)
	if err != nil {
		err = fmt.Errorf("%s: %w", ErrEventUpdate, err)
		return
	}
	row := stmt.QueryRow(data)
	if err = row.Err(); err != nil {
		err = fmt.Errorf("%s: %w", ErrEventUpdate, err)
		return
	}
	evt := Event{}
	err = row.StructScan(&evt)
	if err != nil {
		err = fmt.Errorf("%s: %w", ErrEventCreate, err)
		return
	}
	if evt.ID == 0 {
		err = ErrEventNotFound
		return
	}
	return r.getEventByID(evt.ID)
}
