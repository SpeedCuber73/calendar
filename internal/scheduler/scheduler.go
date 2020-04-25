package scheduler

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/bobrovka/calendar/internal/scheduler/producer"
	"go.uber.org/zap"
)

type Scheduler struct {
	producer *producer.ProducerMQ
	wg       *sync.WaitGroup
	finished chan error
	storage  EventStorage
	logger   *zap.SugaredLogger
}

func NewScheduler(producer *producer.ProducerMQ, storage EventStorage, logger *zap.SugaredLogger) *Scheduler {
	return &Scheduler{
		producer: producer,
		wg:       &sync.WaitGroup{},
		finished: make(chan error),
		storage:  storage,
		logger:   logger,
	}
}

func (s *Scheduler) sendNotifications() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	events, err := s.storage.PopNotifications(ctx)
	if err != nil {
		s.logger.Warnw("error get notifications", "MethodName", "sendNotifications", "err", err)
		return err
	}

	for _, e := range events {
		body, err := json.Marshal(e)
		if err != nil {
			s.logger.Warnw("error marshal event", "MethodName", "sendNotifications", "err", err)
			return err
		}

		err = s.producer.Publish(body)
		if err != nil {
			s.logger.Warnw("error publish event", "MethodName", "sendNotifications", "err", err)
			return err
		}
	}

	return nil
}

func (s *Scheduler) Run() error {
	go func() {
		err := s.producer.KeepConnection()
		if err != nil {
			s.logger.Warnw("connection ends", "err", err)
		}
		s.finished <- err
		close(s.finished)
	}()

	ticker := time.NewTicker(5 * time.Second)

	s.wg.Add(1)
	defer s.wg.Done()
	for {
		select {
		case err := <-s.finished:
			return err
		case <-ticker.C:
			s.sendNotifications()
		}
	}
}

func (s *Scheduler) Stop() error {
	err := s.producer.GracefulStop()
	s.wg.Wait()
	return err
}
