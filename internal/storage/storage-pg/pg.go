package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/bobrovka/calendar/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// insert into events(uuid, title,start_at,duration,user_name,notify_at,notified,descr) VALUES ('1','tit','2020-01-01',2131312,'Kira','2020-01-02',false,'some description');
type event struct {
	UUID        string
	Title       string
	StartAt     time.Time `db:"start_at"`
	Duration    time.Duration
	Description string    `db:"descr"`
	User        string    `db:"user_name"`
	NotifyAt    time.Time `db:"notify_at"`
}

// StoragePg ...
type StoragePg struct {
	db *sqlx.DB
}

// NewStoragePg ...
func NewStoragePg(user, password, host string, port int, name string) (*StoragePg, error) {
	db, err := sqlx.Connect("pgx", fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		name,
	))
	if err != nil {
		return nil, err
	}

	return &StoragePg{
		db: db,
	}, nil
}

// ListEvents ...
func (pg *StoragePg) ListEvents(ctx context.Context, user string, from time.Time, to time.Time) ([]*models.Event, error) {
	rows, err := pg.db.QueryxContext(ctx, `SELECT uuid, title, start_at, duration, descr, user_name, notify_at
	FROM events
	WHERE user_name=$1 AND $2<start_at AND start_at<$3`, user, from, to)
	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for rows.Next() {
		var e event
		err = rows.StructScan(&e)
		if err != nil {
			return nil, err
		}

		events = append(events, toEventModel(&e))
	}

	return events, nil
}

// CreateEvent ...
func (pg *StoragePg) CreateEvent(ctx context.Context, event *models.Event) (string, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return "", nil
	}

	_, err = pg.db.ExecContext(ctx, `INSERT INTO events 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`, uuid.String(), event.Title, event.StartAt, event.Duration, event.Description, event.User, event.StartAt.Add(-event.NotifyBefore))
	if err != nil {
		return "", err
	}

	return uuid.String(), nil
}

// UpdateEvent ...
func (pg *StoragePg) UpdateEvent(ctx context.Context, uuid string, event *models.Event) error {
	_, err := pg.db.ExecContext(ctx, `UPDATE events 
	SET title=$1, 
	start_at=$2, 
	duration=$3, 
	descr=$4, 
	user_name=$5, 
	notify_at=$6 
	WHERE uuid=$7`, event.Title, event.StartAt, event.Duration, event.Description, event.User, event.StartAt.Add(-event.NotifyBefore), uuid)
	if err != nil {
		return err
	}

	return nil
}

// DeleteEvent ...
func (pg *StoragePg) DeleteEvent(ctx context.Context, uuid string) error {
	_, err := pg.db.ExecContext(ctx, `DELETE FROM events 
	WHERE uuid=$1`, uuid)
	if err != nil {
		return err
	}

	return nil
}

func (pg *StoragePg) PopNotifications(ctx context.Context) ([]*models.Event, error) {
	tx, err := pg.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	rows, err := tx.QueryxContext(ctx, `SELECT uuid, title, start_at, duration, descr, user_name 
	FROM events
	WHERE notified != true AND notify_at<now()`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var events []*models.Event
	for rows.Next() {
		var e models.Event
		err = rows.StructScan(&e)
		if err != nil {
			fmt.Println("is it? ", err)
			tx.Rollback()
			return nil, err
		}

		events = append(events, &e)
	}

	for _, e := range events {
		_, err := tx.ExecContext(ctx, "UPDATE events SET notified=true WHERE uuid=$1", e.UUID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return events, nil
}

func toEventModel(e *event) *models.Event {
	return &models.Event{
		UUID:         e.UUID,
		Title:        e.Title,
		StartAt:      e.StartAt,
		Duration:     e.Duration,
		Description:  e.Description,
		User:         e.User,
		NotifyBefore: e.StartAt.Sub(e.NotifyAt),
	}
}
