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
	events, err := a.storage.ListEvents()
	if err != nil {
		return []Event{}
	}
	return events
}

// AddNewEvent добавит новое событие
func (a *App) AddNewEvent(newEvent *Event) error {
	// this should be like one transaction
	currentEvents, err := a.storage.ListEvents()
	if err != nil {
		return err
	}
	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.EndAt) {
		return ErrTimeBusy
	}
	return a.storage.CreateEvent(newEvent)
}

// RemoveEvent добавит новое событие
func (a *App) RemoveEvent(id int) error {
	return a.storage.DeleteEvent(id)
}

// ChangeEvent добавит новое событие
func (a *App) ChangeEvent(id int, newEvent *Event) error {
	// this should be like one transaction
	currentEvents, err := a.storage.ListEvents()
	if err != nil {
		return err
	}

	for i, event := range currentEvents {
		if event.ID == id {
			currentEvents = append(currentEvents[:i], currentEvents[i+1:]...)
			break
		}
	}

	if !hasFreeTime(currentEvents, newEvent.StartAt, newEvent.EndAt) {
		return ErrTimeBusy
	}
	return a.storage.UpdateEvent(id, newEvent)
}
