package app

import (
	"testing"
	"time"

	mock "github.com/bobrovka/calendar/internal/calendar-app/storage-mock"
	"github.com/bobrovka/calendar/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestApp_CreateEvent(t *testing.T) {
	type testCase struct {
		newEvent           *models.Event
		listEventsResponse []*models.Event
		expUuid            string
		expErr             error
	}

	testCases := make(map[string]testCase)

	testCases["Event for free time"] = testCase{
		newEvent: &models.Event{
			UUID:         "1",
			Title:        "first",
			StartAt:      time.Date(2020, time.February, 29, 15, 30, 0, 0, time.UTC), // 15:30
			Duration:     2 * time.Hour,
			Description:  "cool meeting",
			User:         "Kira",
			NotifyBefore: 3 * time.Hour,
		},
		expUuid: "100",
	}
	testCases["Event for busy time"] = testCase{
		newEvent: &models.Event{
			UUID:         "1",
			Title:        "first",
			StartAt:      time.Date(2020, time.February, 29, 15, 30, 0, 0, time.UTC), // 15:30
			Duration:     2 * time.Hour,
			Description:  "cool meeting",
			User:         "Kira",
			NotifyBefore: 3 * time.Hour,
		},
		listEventsResponse: []*models.Event{
			&models.Event{
				UUID:         "2",
				Title:        "second",
				StartAt:      time.Date(2020, time.February, 29, 16, 30, 0, 0, time.UTC), // 16:30
				Duration:     2 * time.Hour,
				Description:  "boring meeting",
				User:         "Kira",
				NotifyBefore: 3 * time.Hour,
			},
		},
		expErr: ErrTimeBusy,
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			storage := &mock.StorageMock{}
			app, err := NewApp(storage)
			assert.NoError(t, err)

			storage.On("ListEvents", v.newEvent.StartAt, v.newEvent.StartAt.AddDate(0, 0, 1)).Return(v.listEventsResponse, nil)
			if v.expErr == nil {
				storage.On("CreateEvent", v.newEvent).Return(v.expUuid, nil)
			}
			uuid, err := app.CreateNewEvent(v.newEvent)
			if err != nil {
				assert.Equal(t, v.expErr, err)
			} else {
				assert.Equal(t, v.expUuid, uuid)
			}

			storage.AssertExpectations(t)
		})
	}
}

func TestApp_ChangeEvent(t *testing.T) {
	type testCase struct {
		uuid               string
		newEvent           *models.Event
		listEventsResponse []*models.Event
		expErr             error
	}

	testCases := make(map[string]testCase)

	testCases["Event not found"] = testCase{
		uuid: "1",
		newEvent: &models.Event{
			Title:        "first",
			StartAt:      time.Date(2020, time.February, 29, 15, 30, 0, 0, time.UTC), // 15:30
			Duration:     2 * time.Hour,
			Description:  "cool meeting",
			User:         "Kira",
			NotifyBefore: 3 * time.Hour,
		},
		listEventsResponse: []*models.Event{
			&models.Event{
				UUID:         "2",
				Title:        "second",
				StartAt:      time.Date(2020, time.February, 29, 16, 30, 0, 0, time.UTC), // 16:30
				Duration:     2 * time.Hour,
				Description:  "boring meeting",
				User:         "Kira",
				NotifyBefore: 3 * time.Hour,
			},
		},
		expErr: ErrNotFound,
	}

	testCases["Event time busy"] = testCase{
		uuid: "1",
		newEvent: &models.Event{
			Title:        "first",
			StartAt:      time.Date(2020, time.February, 29, 15, 30, 0, 0, time.UTC), // 15:30
			Duration:     2 * time.Hour,
			Description:  "cool meeting",
			User:         "Kira",
			NotifyBefore: 3 * time.Hour,
		},
		listEventsResponse: []*models.Event{
			&models.Event{
				UUID:         "1",
				Title:        "first",
				StartAt:      time.Date(2020, time.February, 29, 16, 30, 0, 0, time.UTC), // 16:30
				Duration:     2 * time.Hour,
				Description:  "boring meeting",
				User:         "Kira",
				NotifyBefore: 3 * time.Hour,
			},
			&models.Event{
				UUID:         "2",
				Title:        "second",
				StartAt:      time.Date(2020, time.February, 29, 16, 30, 0, 0, time.UTC), // 16:30
				Duration:     2 * time.Hour,
				Description:  "boring meeting",
				User:         "Kira",
				NotifyBefore: 3 * time.Hour,
			},
		},
		expErr: ErrTimeBusy,
	}

	testCases["Event successfull update"] = testCase{
		uuid: "1",
		newEvent: &models.Event{
			Title:        "first",
			StartAt:      time.Date(2020, time.February, 29, 15, 30, 0, 0, time.UTC), // 15:30
			Duration:     2 * time.Hour,
			Description:  "cool meeting",
			User:         "Kira",
			NotifyBefore: 3 * time.Hour,
		},
		listEventsResponse: []*models.Event{
			&models.Event{
				UUID:         "1",
				Title:        "first",
				StartAt:      time.Date(2020, time.February, 29, 16, 30, 0, 0, time.UTC), // 16:30
				Duration:     2 * time.Hour,
				Description:  "boring meeting",
				User:         "Kira",
				NotifyBefore: 3 * time.Hour,
			},
			&models.Event{
				UUID:         "2",
				Title:        "second",
				StartAt:      time.Date(2020, time.February, 29, 11, 30, 0, 0, time.UTC), // 16:30
				Duration:     2 * time.Hour,
				Description:  "boring meeting",
				User:         "Kira",
				NotifyBefore: 3 * time.Hour,
			},
		},
	}

	for k, v := range testCases {
		t.Run(k, func(t *testing.T) {
			storage := &mock.StorageMock{}
			app, err := NewApp(storage)
			assert.NoError(t, err)

			storage.On("ListEvents", v.newEvent.StartAt, v.newEvent.StartAt.AddDate(0, 0, 1)).Return(v.listEventsResponse, nil)
			if v.expErr == nil {
				storage.On("UpdateEvent", v.uuid, v.newEvent).Return(nil)
			}
			err = app.ChangeEvent(v.uuid, v.newEvent)
			assert.Equal(t, v.expErr, err)

			storage.AssertExpectations(t)
		})
	}
}
