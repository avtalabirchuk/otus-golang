package service

import (
	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"

	grpc "github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/service/server"
)

type Service interface {
	RunHTTP(string, string) error
	RunGRPC(string) error
}

func New(r repository.CRUD) Service {
	return grpc.New(r)
}
