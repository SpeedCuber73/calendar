package service

import (
	"context"

	app "github.com/bobrovka/calendar/internal/calendar-app"
	"github.com/bobrovka/calendar/internal/models"
	"github.com/bobrovka/calendar/pkg/calendar/api"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
)

// EventService is implementation for grpc event service
type EventService struct {
	app    app.App
	logger *zap.SugaredLogger
}

// NewEventService creates new instance of grpc event service
func NewEventService(app app.App, logger *zap.SugaredLogger) *EventService {
	return &EventService{
		app:    app,
		logger: logger,
	}
}

// ListEvents method
func (es *EventService) ListEvents(ctx context.Context, request *api.ListRequest) (*api.ListResponse, error) {
	day, err := ptypes.Timestamp(request.GetDate())
	if err != nil {
		es.logger.Errorw("error time conversion", "methodName", "ListEvents", "err", err)
		return nil, err
	}

	var events []*models.Event
	switch request.GetPeriod() {
	case api.Period_DAY:
		events, err = es.app.ListDayEvents(ctx, request.User, day)
		if err != nil {
			es.logger.Errorw("error ListDayEvents", "methodName", "ListEvents", "err", err)
			return nil, err
		}
	case api.Period_WEEK:
		events, err = es.app.ListWeekEvents(ctx, request.User, day)
		if err != nil {
			es.logger.Errorw("error ListWeekEvents", "methodName", "ListEvents", "err", err)
			return nil, err
		}
	case api.Period_MONTH:
		events, err = es.app.ListMonthEvents(ctx, request.User, day)
		if err != nil {
			es.logger.Errorw("error ListMonthEvents", "methodName", "ListEvents", "err", err)
			return nil, err
		}
	}

	result := make([]*api.Event, 0, len(events))
	for _, event := range events {
		startAt, err := ptypes.TimestampProto(event.StartAt)
		if err != nil {
			es.logger.Errorw("error time conversion", "methodName", "ListEvents", "err", err)
			return nil, err
		}

		result = append(result, &api.Event{
			Uuid:         event.UUID,
			Title:        event.Title,
			StartAt:      startAt,
			Duration:     ptypes.DurationProto(event.Duration),
			Description:  event.Description,
			User:         event.User,
			NotifyBefore: ptypes.DurationProto(event.Duration),
		})
	}

	es.logger.Infow("Success ListEvents")
	return &api.ListResponse{
		Events: result,
	}, nil
}

// CreateEvent method
func (es *EventService) CreateEvent(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	newEvent := request.GetEvent()

	startAt, err := ptypes.Timestamp(newEvent.GetStartAt())
	if err != nil {
		es.logger.Errorw("error time conversion", "methodName", "CreateEvent", "err", err)
		return nil, err
	}

	duration, err := ptypes.Duration(newEvent.GetDuration())
	if err != nil {
		es.logger.Errorw("error duration conversion", "methodName", "CreateEvent", "err", err)
		return nil, err
	}

	notifyBefore, err := ptypes.Duration(newEvent.GetNotifyBefore())
	if err != nil {
		es.logger.Errorw("error duration conversion", "methodName", "CreateEvent", "err", err)
		return nil, err
	}

	uuid, err := es.app.CreateNewEvent(ctx, &models.Event{
		Title:        newEvent.GetTitle(),
		StartAt:      startAt,
		Duration:     duration,
		Description:  newEvent.GetDescription(),
		User:         newEvent.GetUser(),
		NotifyBefore: notifyBefore,
	})
	if err != nil {
		es.logger.Errorw("error CreateNewEvent", "methodName", "CreateEvent", "err", err)
		return nil, err
	}

	es.logger.Infow("Success CreateEvent", "UUID", uuid)
	return &api.CreateResponse{
		Uuid: uuid,
	}, nil
}

// UpdateEvent method
func (es *EventService) UpdateEvent(ctx context.Context, request *api.UpdateRequest) (*empty.Empty, error) {
	uuid := request.GetUuid()
	updatedEvent := request.GetEvent()

	startAt, err := ptypes.Timestamp(updatedEvent.GetStartAt())
	if err != nil {
		es.logger.Errorw("error time conversion", "methodName", "UpdateEvent", "err", err)
		return nil, err
	}

	duration, err := ptypes.Duration(updatedEvent.GetDuration())
	if err != nil {
		es.logger.Errorw("error duration conversion", "methodName", "UpdateEvent", "err", err)
		return nil, err
	}

	notifyBefore, err := ptypes.Duration(updatedEvent.GetNotifyBefore())
	if err != nil {
		es.logger.Errorw("error duration conversion", "methodName", "UpdateEvent", "err", err)
		return nil, err
	}

	err = es.app.ChangeEvent(ctx, uuid, &models.Event{
		Title:        updatedEvent.GetTitle(),
		StartAt:      startAt,
		Duration:     duration,
		Description:  updatedEvent.GetDescription(),
		User:         updatedEvent.GetUser(),
		NotifyBefore: notifyBefore,
	})
	if err != nil {
		es.logger.Errorw("error ChangeEvent", "methodName", "UpdateEvent", "err", err)
		return nil, err
	}

	es.logger.Infow("Success ChangeEvent", "UUID", uuid)
	return &empty.Empty{}, nil
}

// DeleteEvent method
func (es *EventService) DeleteEvent(ctx context.Context, request *api.DeleteRequest) (*empty.Empty, error) {
	uuid := request.GetUuid()

	err := es.app.RemoveEvent(ctx, uuid)
	if err != nil {
		es.logger.Errorw("error RemoveEvent", "methodName", "DeleteEvent", "err", err)
		return nil, err
	}

	es.logger.Infow("Success DeleteEvent", "UUID", uuid)
	return &empty.Empty{}, nil
}
