package service

import (
	grpc "github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/service/server"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

type Service interface {
	RunHTTP(string, string) error
	RunGRPC(string) error
}

func New(r repository.Base) Service {
	return grpc.New(r)
}
