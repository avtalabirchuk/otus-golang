package sqlstorage

import "github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/storage"

func New() storage.Storage {
	return &store{}
}
