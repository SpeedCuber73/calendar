package mock

import (
	"context"
	"time"

	"github.com/bobrovka/calendar/internal/models"
	"github.com/stretchr/testify/mock"
)

// StorageMock мок хранилища
type StorageMock struct {
	mock.Mock
}

// ListEvents мокирует метод
func (m *StorageMock) ListEvents(ctx context.Context, user string, from, to time.Time) ([]*models.Event, error) {
	args := m.Called(ctx, user, from, to)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	return args.Get(0).([]*models.Event), err
}

// CreateEvent мокирует метод
func (m *StorageMock) CreateEvent(ctx context.Context, event *models.Event) (string, error) {
	args := m.Called(ctx, event)
	err := args.Error(1)
	if err != nil {
		return "", err
	}

	return args.String(0), err
}

// UpdateEvent мокирует метод
func (m *StorageMock) UpdateEvent(ctx context.Context, id string, event *models.Event) error {
	args := m.Called(ctx, id, event)
	return args.Error(0)
}

// DeleteEvent мокирует метод
func (m *StorageMock) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *StorageMock) PopNotifications(ctx context.Context) ([]*models.Event, error) {
	args := m.Called(ctx)
	err := args.Error(1)
	if err != nil {
		return nil, err
	}

	return args.Get(0).([]*models.Event), err
}
