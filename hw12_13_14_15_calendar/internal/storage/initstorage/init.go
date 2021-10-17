package initstorage

import (
	"context"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/storage/sql"
)

func New(ctx context.Context, inmem bool, connect string) (storage.Storage, error) {
	var db storage.Storage
	if inmem {
		db = memorystorage.New()
	} else {
		db = sqlstorage.New()
	}
	err := db.Connect(ctx, connect)
	return db, err
}
