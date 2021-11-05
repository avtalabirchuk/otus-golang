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

	// is used for init postgres.
	_ "github.com/lib/pq"
)

var ErrDBOpen = errors.New("database open error")

type PSQLRepo struct {
	db            *sqlx.DB
	itemsPerQuery int
	maxConn       int
	validator     *validator.Validate
}

var selectCurrentEventsQs = `SELECT e.* from events e join events_status es on e.id = es.event_id
WHERE DATE_PART('day', start_date::timestamp - NOW()) <= e.notified_for and e.end_date > NOW() and es.status = 'New';`

var selectObsoleteEventsQs = `SELECT id FROM events WHERE DATE_PART('year', NOW()) - DATE_PART('year', end_date::timestamp) >= 1`

var insertQs = `INSERT INTO events
(user_id, title, description, start_date, end_date, notified_for)
VALUES
(:user_id, :title, :description, :start_date, :end_date, :notified_for)
RETURNING id`

var insertStatusQs = `INSERT INTO events_status (event_id) VALUES ($1)`

var updateQs = `UPDATE events
	SET (title, description, start_date, end_date, notified_for) = (:title, :description, :start_date, :end_date, :notified_for)
	WHERE id = :id
	RETURNING id`

var updateEventsStatusQs = `UPDATE events_status SET status = :status WHERE id IN (:ids);`

var (
	deleteQs               = `DELETE FROM events where id = $1 RETURNING id`
	deleteStatusQs         = `DELETE FROM events_status where event_id = $1`
	deleteObsoleteEventsQs = fmt.Sprintf(`DELETE FROM events WHERE id IN (%s);`, selectObsoleteEventsQs)
	deleteObsoleteStatusQs = fmt.Sprintf(`DELETE FROM events_status WHERE event_id IN (%s);`, selectObsoleteEventsQs)
)

func (r *PSQLRepo) Connect(ctx context.Context, url string) (err error) {
	log.Debug().Msgf("Connecting to %s", url)
	r.db, err = sqlx.Connect("postgres", url)
	if err != nil {
		return fmt.Errorf("%s: %w", ErrDBOpen, err)
	}
	if r.maxConn != 0 {
		r.db.SetMaxOpenConns(r.maxConn)
	}
	return r.db.PingContext(ctx)
}

func (r *PSQLRepo) Close() error {
	return r.db.Close()
}

func NewPSQLRepo(args ...interface{}) *PSQLRepo {
	nums := make([]int, len(args))
	for i, el := range args {
		if n, ok := el.(int); ok {
			nums[i] = n
		}
	}
	return &PSQLRepo{
		itemsPerQuery: nums[0],
		maxConn:       nums[1],
		validator:     NewEventValidator(),
	}
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
	return r.getEventsBetween(date, date.AddDate(0, 0, 1))
}

func (r *PSQLRepo) GetWeekEvents(date time.Time) (result []Event, err error) {
	return r.getEventsBetween(date, date.AddDate(0, 0, 7))
}

func (r *PSQLRepo) GetMonthEvents(date time.Time) (result []Event, err error) {
	return r.getEventsBetween(date, date.AddDate(0, 1, 0))
}

func (r *PSQLRepo) DeleteEvent(id int64) (err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return
	}
	_, err = tx.Exec(deleteStatusQs, id)
	if err != nil {
		return
	}
	result, err := tx.Exec(deleteQs, id)
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
	if err := tx.Commit(); err != nil {
		return err
	}
	return
}

func (r *PSQLRepo) CreateEvent(data Event) (event Event, err error) {
	err = r.validator.Struct(data)
	if err != nil {
		return
	}
	tx, err := r.db.Beginx()
	if err != nil {
		return
	}
	stmt, err := tx.PrepareNamed(insertQs)
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
	_, err = tx.Exec(insertStatusQs, evt.ID)
	if err != nil {
		return
	}

	if err := tx.Commit(); err != nil {
		return evt, err
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

func (r *PSQLRepo) GetCurrentEvents() (result []Event, err error) {
	err = r.db.Select(&result, selectCurrentEventsQs)
	return
}

func (r *PSQLRepo) processEvents(events *[]Event, queryString string, status string) error {
	if len(*events) == 0 {
		return nil
	}
	ids := make([]int64, len(*events))
	for i, el := range *events {
		ids[i] = el.ID
	}
	arg := map[string]interface{}{
		"ids":    ids,
		"status": status,
	}
	query, args, err := sqlx.Named(queryString, arg)
	if err != nil {
		return err
	}
	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)
	row, err := r.db.Queryx(query, args...)
	if err != nil {
		return err
	}
	defer row.Close()

	if err != nil {
		return err
	}
	return row.Err()
}

func (r *PSQLRepo) MarkEventsAsProcessing(events *[]Event) error {
	return r.processEvents(events, updateEventsStatusQs, "Processing")
}

func (r *PSQLRepo) MarkEventsAsSent(events *[]Event) error {
	return r.processEvents(events, updateEventsStatusQs, "Sent")
}

func (r *PSQLRepo) DeleteObsoleteEvents() (err error) {
	tx, err := r.db.Beginx()
	if err != nil {
		return
	}
	_, err = tx.Exec(deleteObsoleteStatusQs)
	if err != nil {
		return
	}
	result, err := tx.Exec(deleteObsoleteEventsQs)
	if err != nil {
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		return
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return
}
