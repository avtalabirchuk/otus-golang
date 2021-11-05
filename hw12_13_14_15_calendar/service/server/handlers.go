package server

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/avtalabirchuk/otus-golang/hw12_13_14_15_calendar/internal/repository"
)

var (
	ErrConvertDBStructToProto = status.Error(codes.InvalidArgument, "cannot convert DB struct to proto one")
	ErrFetchingDataFromDB     = status.Error(codes.Internal, "error occurred during fetching data from DB")
	ErrUnsupportedRequest     = status.Error(codes.Unimplemented, "request type is not supported")
	ErrIncomingTimeStampError = status.Error(codes.InvalidArgument, "provided timestamp is invalid")
)

type Service struct {
	r repository.CRUD
}

func New(r repository.CRUD) *Service {
	return &Service{r: r}
}

func processEvents(repo repository.CRUD, query *QueryEventsRequest) (events []repository.Event, err error) {
	startDate := time.Now()
	if ts := query.GetTs(); ts != 0 {
		startDate = time.Unix(ts, 0)
		if err != nil {
			return nil, ErrIncomingTimeStampError
		}
	}
	log.Debug().Msgf("Making query with startDate %s", startDate)
	switch query.GetType() {
	case QueryRangeType_DAY:
		events, err = repo.GetDayEvents(startDate)
	case QueryRangeType_WEEK:
		events, err = repo.GetWeekEvents(startDate)
	case QueryRangeType_MONTH:
		events, err = repo.GetMonthEvents(startDate)
	case QueryRangeType_UNKNOWN:
	default:
		return nil, ErrUnsupportedRequest
	}
	log.Debug().Msgf("Getting events %+v, err %s", events, err)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrFetchingDataFromDB, err)
	}
	return
}

func (s *Service) GetEvents(ctx context.Context, query *QueryEventsRequest) (result *EventsResponse, err error) {
	dbEvents, err := processEvents(s.r, query)
	result = &EventsResponse{}
	if err != nil {
		log.Error().Msgf("%s", err)
		return
	}

	for _, evt := range dbEvents {
		converted, err := ConvertEventToProto(evt)
		if err != nil {
			log.Error().Msgf("%s %s", ErrConvertDBStructToProto, err)
			return nil, ErrConvertDBStructToProto
		}
		result.Events = append(result.Events, converted)
	}
	return
}

func (s *Service) CreateEvent(ctx context.Context, event *Event) (result *Event, err error) {
	obj, err := ConvertEventFromProto(event)
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, "%s", err)
		return
	}
	evt, err := s.r.CreateEvent(*obj)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			err = status.Errorf(codes.InvalidArgument, "%s", err)
		}
		return
	}
	log.Debug().Msgf("[GRPC] Created Event: %+v\n", evt)
	return ConvertEventToProto(evt)
}

func (s *Service) UpdateEvent(ctx context.Context, data *UpdateEventRequest) (result *Event, err error) {
	obj, err := ConvertEventFromProto(data.Event)
	if err != nil {
		err = status.Errorf(codes.InvalidArgument, "%s", err)
		return
	}
	evt, err := s.r.UpdateEvent(data.Id, *obj)
	if err != nil {
		if _, ok := err.(validator.ValidationErrors); ok {
			err = status.Errorf(codes.InvalidArgument, "%s", err)
		}
		return
	}
	log.Debug().Msgf("[GRPC] Updated Event: %+v\n", evt)
	return ConvertEventToProto(evt)
}

func (s *Service) DeleteEvent(ctx context.Context, data *DeleteEventRequest) (*empty.Empty, error) {
	return &empty.Empty{}, s.r.DeleteEvent(data.Id)
}
