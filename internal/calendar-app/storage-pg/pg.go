package pg

import (
	"context"
	"time"

	"github.com/bobrovka/calendar/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StoragePg struct {
	db *sqlx.DB
}

func NewStoragePg(db *sqlx.DB) (*StoragePg, error) {
	return &StoragePg{
		db: db,
	}, nil
}

func (pg *StoragePg) ListEvents(ctx context.Context, user string, from time.Time, to time.Time) ([]*models.Event, error) {
	rows, err := pg.db.QueryxContext(ctx, `SELECT uuid, title, start_at, duration, descr, user_name, notify_before
	FROM events
	WHERE user_name=$1 AND $2<start_at AND start_at<$3`, user, from, to)
	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for rows.Next() {
		var e models.Event
		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}
		events = append(events, &e)
	}

	return events, nil
}

func (pg *StoragePg) CreateEvent(ctx context.Context, event *models.Event) (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	_, err = pg.db.ExecContext(ctx, `INSERT INTO events 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, uuid.String(), event.Title, event.StartAt, event.Duration, event.Description, event.User, event.NotifyBefore)
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

func (pg *StoragePg) UpdateEvent(ctx context.Context, id string, event *models.Event) error {
	panic("not implemented") // TODO: Implement
}

func (pg *StoragePg) DeleteEvent(ctx context.Context, id string) error {
	panic("not implemented") // TODO: Implement
}
