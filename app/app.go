package app

// App сущность, описывающая бизнес-логику сервиса
type App struct {
	storage EventStorage
}

// CreateApp создает новый инстанс приложения
func CreateApp(storage EventStorage) (*App, error) {
	return &App{
		storage: storage,
	}, nil
}

// ListAllEvents вернет список событий
func (a *App) ListAllEvents() []Event {
	events, _ := a.storage.ListEvents()
	return events
}

// AddNewEvent добавит новое событие
func (a *App) AddNewEvent(newEvent *Event) error {
	return a.storage.CreateEvent(newEvent)
}

// RemoveEvent добавит новое событие
func (a *App) RemoveEvent(id int) error {
	return a.storage.DeleteEvent(id)
}

// ChangeEvent добавит новое событие
func (a *App) ChangeEvent(id int, newEvent *Event) error {
	return a.storage.UpdateEvent(id, newEvent)
}
