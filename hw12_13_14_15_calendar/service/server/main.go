package server

import (
	"context"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/logmiddleware"
)

func (s *Service) RunGRPC(addr string) error {
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_middleware.ChainUnaryServer(logmiddleware.ApplyGRPC()),
		),
	)

	RegisterCalendarServer(server, s)

	log.Info().Msgf("Starting GRPC server on %s", lsn.Addr().String())

	return server.Serve(lsn)
}

func (s *Service) RunHTTP(grpcAddr, addr string) error {
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer grpcConn.Close()

	grpcGwMux := runtime.NewServeMux()

	err = RegisterCalendarHandler(
		context.Background(),
		grpcGwMux,
		grpcConn,
	)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", logmiddleware.ApplyHTTP(grpcGwMux))

	log.Info().Msgf("Starting HTTP server on %s", addr)
	return http.ListenAndServe(addr, mux)
}
