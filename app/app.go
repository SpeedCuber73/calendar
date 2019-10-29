package app

import (
	"github.com/SpeedCuber73/calendar/storage"
)

type App struct {
	storage *storage.EventStorage
}

func CreateApp(storage *storage.EventStorage) (*App, error) {
	return &App{
		storage: storage.EventStorage,
	}, nil
}

// ListEvents вернет список событий
func (a *App) ListEvents() []storage.Event {
	events, _ := a.storage.ListEvents()
	return events
}
