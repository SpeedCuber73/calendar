package service

import (
	"context"

	"github.com/bobrovka/calendar/internal/grpc/api"
)

type EventService struct{}

func (es *EventService) GetEventsForDay(_ context.Context, _ *api.Day) (*api.EventsResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (es *EventService) GetEventsForWeek(_ context.Context, _ *api.Day) (*api.EventsResponse, error) {
	panic("not implemented") // TODO: Implement
}

func (es *EventService) GetEventsForMonth(_ context.Context, _ *api.Day) (*api.EventsResponse, error) {
	panic("not implemented") // TODO: Implement
}
