package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/bobrovka/calendar/internal/models"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// App интерфейс приложения
type App interface {
	ListDayEvents(ctx context.Context, user string, date time.Time) ([]*models.Event, error)
	ListWeekEvents(ctx context.Context, user string, date time.Time) ([]*models.Event, error)
	ListMonthEvents(ctx context.Context, user string, date time.Time) ([]*models.Event, error)
	CreateNewEvent(ctx context.Context, newEvent *models.Event) (string, error)
	RemoveEvent(ctx context.Context, uuid string) error
	ChangeEvent(ctx context.Context, uuid string, newEvent *models.Event) error

	RunScheduler(ctx context.Context, ch *amqp.Channel) error
}

// Calendar сущность, описывающая бизнес-логику сервиса
type Calendar struct {
	storage EventStorage
	logger  *zap.SugaredLogger
}

// NewCalendar создает новый инстанс приложения
func NewCalendar(storage EventStorage, logger *zap.SugaredLogger) (App, error) {
	return &Calendar{
		storage: storage,
		logger:  logger,
	}, nil
}

// ListDayEvents вернет список событий на день
func (a *Calendar) ListDayEvents(ctx context.Context, user string, date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(ctx, user, date, date.AddDate(0, 0, 1))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListWeekEvents вернет список событий на неделю
func (a *Calendar) ListWeekEvents(ctx context.Context, user string, date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(ctx, user, date, date.AddDate(0, 0, 7))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// ListMonthEvents вернет список событий на месяц
func (a *Calendar) ListMonthEvents(ctx context.Context, user string, date time.Time) ([]*models.Event, error) {
	events, err := a.storage.ListEvents(ctx, user, date, date.AddDate(0, 1, 0))
	if err != nil {
		return nil, err
	}
	return events, nil
}

// CreateNewEvent добавит новое событие
func (a *Calendar) CreateNewEvent(ctx context.Context, newEvent *models.Event) (string, error) {
	currentEvents, err := a.storage.ListEvents(ctx, newEvent.User, time.Unix(0, 0), time.Unix(67098285000, 0))
	if err != nil {
		return "", err
	}

	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.StartAt.Add(newEvent.Duration)) {
		return "", ErrTimeBusy
	}

	return a.storage.CreateEvent(ctx, newEvent)
}

// RemoveEvent удалит событие
func (a *Calendar) RemoveEvent(ctx context.Context, uuid string) error {
	return a.storage.DeleteEvent(ctx, uuid)
}

// ChangeEvent изменит событие
func (a *Calendar) ChangeEvent(ctx context.Context, uuid string, newEvent *models.Event) error {
	// get events on this day
	currentEvents, err := a.storage.ListEvents(ctx, newEvent.User, time.Unix(0, 0), time.Unix(67098285000, 0))
	if err != nil {
		return err
	}

	var found bool
	// delete an event that is being modified
	for i, event := range currentEvents {
		if event.UUID == uuid {
			currentEvents = append(currentEvents[:i], currentEvents[i+1:]...)
			found = true
			break
		}
	}
	if !found {
		return ErrNotFound
	}

	// if no free time - abort changing
	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.StartAt.Add(newEvent.Duration)) {
		return ErrTimeBusy
	}

	return a.storage.UpdateEvent(ctx, uuid, newEvent)
}

func (a *Calendar) RunScheduler(ctx context.Context, ch *amqp.Channel) error {
	ticker := time.NewTicker(5 * time.Second)

	q, err := ch.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				func() {
					events, err := a.storage.PopNotifications(ctx)
					if err != nil {
						a.logger.Warnw("error get notifications", "MethodName", "RunScheduler", "err", err)
						return
					}

					for _, e := range events {
						body, err := json.Marshal(e)
						if err != nil {
							a.logger.Warnw("error marshal event", "MethodName", "RunScheduler", "err", err)
							return
						}

						err = ch.Publish(
							"",     // exchange
							q.Name, // routing key
							false,  // mandatory
							false,  // immediate
							amqp.Publishing{
								ContentType: "application/json",
								Body:        []byte(body),
							},
						)
					}
				}()
			}
		}
	}()

	return nil
}

func hasFreeTime(existingEvents []*models.Event, start, end time.Time) bool {
	for _, event := range existingEvents {
		eventEndAt := event.StartAt.Add(event.Duration)
		if (event.StartAt.Before(start) || event.StartAt.Equal(start)) && eventEndAt.After(start) {
			return false
		}
		if event.StartAt.After(start) && event.StartAt.Before(end) {
			return false
		}
	}

	return true
}
