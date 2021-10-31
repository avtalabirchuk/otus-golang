package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/config"
	// is used for init postgres.
	_ "github.com/lib/pq"
)

var ErrDBOpen = errors.New("database open error")

type PSQLRepo struct {
	db            *sqlx.DB
	itemsPerQuery int
	validator     *validator.Validate
}

var insertQs = `INSERT INTO events
(user_id, title, description, start_date, end_date, notified_at)
VALUES
(:user_id, :title, :description, :start_date, :end_date, :notified_at)
RETURNING id`

var updateQs = `UPDATE events
	SET (title, description, start_date, end_date, notified_at) = (:title, :description, :start_date, :end_date, :notified_at)
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
	return &PSQLRepo{itemsPerQuery: 100, validator: NewEventValidator()}
}

func (r *PSQLRepo) getEventsBetween(startPeriod time.Time, endPeriod time.Time) (result []Event, err error) {
	query := "SELECT * FROM events WHERE (start_date <= $1 and end_date >= $1) or (start_date <= $2 and end_date >= $2) or (start_date >= $1 and end_date <= $2) ORDER BY start_date ASC LIMIT $3"
	err = r.db.Select(&result, query, startPeriod, endPeriod, r.itemsPerQuery)
	return
}

func (r *PSQLRepo) getEventByID(id int64) (result Event, err error) {
	err = r.db.Get(&result, "SELECT * FROM events WHERE id = $1", id)
	if err != nil {
		log.Debug().Msgf("[DB] getEventByID Err %d, %+v, %s", id, result, err)
	}
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
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		return
	}
	if n == 0 {
		err = sql.ErrNoRows
	}
	return
}

func (r *PSQLRepo) CreateEvent(data Event) (event Event, err error) {
	err = r.validator.Struct(data)
	if err != nil {
		return
	}
	stmt, err := r.db.PrepareNamed(insertQs)
	if err != nil {
		return
	}
	row := stmt.QueryRow(data)
	if err = row.Err(); err != nil {
		return
	}
	evt := Event{}
	err = row.StructScan(&evt)
	if err != nil {
		return
	}
	if evt.ID == 0 {
		err = ErrEventNotFound
		return
	}
	return r.getEventByID(evt.ID)
}

func (r *PSQLRepo) UpdateEvent(id int64, data Event) (event Event, err error) {
	// Fetch Event from DB
	event, err = r.getEventByID(id)
	if err != nil {
		return
	}
	// Merge data with event
	err = MergeEvents(&event, data)
	if err != nil {
		return
	}
	log.Debug().Msgf("Merged Object %+v", event)
	// Validate received object
	err = r.validator.Struct(event)
	if err != nil {
		return
	}
	// Validate received object
	stmt, err := r.db.PrepareNamed(updateQs)
	if err != nil {
		return
	}
	row := stmt.QueryRow(event)
	if err = row.Err(); err != nil {
		return
	}
	return r.getEventByID(id)
}
