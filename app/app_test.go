package app

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestBusyStorage struct {
	exists []Event
}

func CreateTestBusyStorage(events []Event) (*TestBusyStorage, error) {
	return &TestBusyStorage{
		exists: events,
	}, nil
}

func (s *TestBusyStorage) ListEvents() ([]Event, error) {
	return s.exists, nil
}

func (s *TestBusyStorage) CreateEvent(event *Event) error {
	return nil
}

func (s *TestBusyStorage) UpdateEvent(id int, event *Event) error {
	return nil
}

func (s *TestBusyStorage) DeleteEvent(id int) error {
	return nil
}

func TestApp(t *testing.T) {
	assert := assert.New(t)

	start, _ := time.Parse(time.RFC3339, "2019-11-11T15:00:00Z")
	end, _ := time.Parse(time.RFC3339, "2019-11-11T17:00:00Z")
	predefinedEvents := []Event{
		Event{
			ID:      1,
			StartAt: start,
			EndAt:   end,
		},
	}

	myStorage, _ := CreateTestBusyStorage(predefinedEvents)
	myApp, _ := NewApp(myStorage)

	// uncross event, expect no problems
	start, _ = time.Parse(time.RFC3339, "2019-11-11T17:00:00Z")
	end, _ = time.Parse(time.RFC3339, "2019-11-11T19:00:00Z")
	uncrossEvent := Event{
		ID:      1,
		StartAt: start,
		EndAt:   end,
	}
	err := myApp.AddNewEvent(&uncrossEvent)
	assert.Nil(err)

	// cross event, expect ErrTimeBusy (high half intersection)
	start, _ = time.Parse(time.RFC3339, "2019-11-11T14:00:00Z")
	end, _ = time.Parse(time.RFC3339, "2019-11-11T16:00:00Z")
	crossEvent := Event{
		ID:      1,
		StartAt: start,
		EndAt:   end,
	}
	err = myApp.AddNewEvent(&crossEvent)
	assert.Equal(ErrTimeBusy, err)

	// cross event, expect ErrTimeBusy (low half intersection)
	start, _ = time.Parse(time.RFC3339, "2019-11-11T16:00:00Z")
	end, _ = time.Parse(time.RFC3339, "2019-11-11T18:00:00Z")
	crossEvent = Event{
		ID:      1,
		StartAt: start,
		EndAt:   end,
	}
	err = myApp.AddNewEvent(&crossEvent)
	assert.Equal(ErrTimeBusy, err)

	// cross event, expect ErrTimeBusy (full intersection)
	start, _ = time.Parse(time.RFC3339, "2019-11-11T15:00:00Z")
	end, _ = time.Parse(time.RFC3339, "2019-11-11T17:00:00Z")
	crossEvent = Event{
		ID:      1,
		StartAt: start,
		EndAt:   end,
	}
	err = myApp.AddNewEvent(&crossEvent)
	assert.Equal(ErrTimeBusy, err)
}
