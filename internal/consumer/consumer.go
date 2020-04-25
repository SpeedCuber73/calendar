package consumer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v3"

	"github.com/streadway/amqp"
)

// Consumer ...
type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	consumerTag  string
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
}

func NewConsumer(consumerTag, uri, exchangeName, exchangeType, queue, bindingKey string) *Consumer {
	return &Consumer{
		consumerTag:  consumerTag,
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queue:        queue,
		bindingKey:   bindingKey,
		done:         make(chan error),
	}
}

func (c *Consumer) reConnect() (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		select {
		case <-time.After(d):
			if err := c.connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}
			msgs, err := c.announceQueue()
			if err != nil {
				fmt.Printf("Couldn't connect: %+v", err)
				continue
			}

			return msgs, nil
		}
	}
}

func (c *Consumer) connect() error {

	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		errMQ := <-c.conn.NotifyClose(make(chan *amqp.Error))
		log.Printf("closing: %s", errMQ)
		if errMQ != nil {
			// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
			c.done <- errors.New("Channel Closed")
		} else {
			c.done <- nil
		}
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	if err = c.channel.ExchangeDeclare(
		c.exchangeName,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	return nil
}

// Задекларировать очередь, которую будем слушать.
func (c *Consumer) announceQueue() (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.channel.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	msgs, err := c.channel.Consume(
		queue.Name,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	return msgs, nil
}

func (c *Consumer) Handle(ctx context.Context, wg *sync.WaitGroup, fn func(<-chan amqp.Delivery), threads int) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	defer func() {
		err := c.gracefulStop()
		if err != nil {
			log.Println("connot gracefully stop consumer ", err)
		}
	}()

	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	for {
		for i := 0; i < threads; i++ {
			wg.Add(1)
			go func() {
				fn(msgs)
				wg.Done()
			}()
		}

		select {
		case err := <-c.done:
			if err != nil {
				msgs, err = c.reConnect()
				if err != nil {
					return fmt.Errorf("Reconnecting Error: %s", err)
				}
			} else {
				return nil
			}
		case <-ctx.Done():
			return nil
		}

		fmt.Println("Reconnected...")
	}
}

func (c *Consumer) gracefulStop() error {
	err := c.channel.Close()
	if err != nil {
		return err
	}
	err = c.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
