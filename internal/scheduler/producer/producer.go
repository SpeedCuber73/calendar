package producer

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/streadway/amqp"
)

type ProducerMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	uri          string
	done         chan error
	exchangeName string
	exchangeType string
	queue        string
	routingKey   string
}

func NewProducerMQ(uri, exchangeName, exchangeType, queue, routingKey string) *ProducerMQ {
	return &ProducerMQ{
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		routingKey:   routingKey,
		done:         make(chan error),
	}
}

func (p *ProducerMQ) reConnect() error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}

		select {
		case <-time.After(d):
			if err := p.connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}

			return nil
		}
	}
}

func (p *ProducerMQ) connect() error {
	var err error

	p.conn, err = amqp.Dial(p.uri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	go func() {
		errMQ := <-p.conn.NotifyClose(make(chan *amqp.Error))
		log.Printf("closing: %s", errMQ)
		if errMQ != nil {
			// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
			p.done <- errors.New("Channel Closed")
		} else {
			p.done <- nil
		}
	}()

	p.channel, err = p.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	if err = p.channel.ExchangeDeclare(
		p.exchangeName,
		p.exchangeType,
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

func (p *ProducerMQ) Publish(msg []byte) error {
	return p.channel.Publish(
		p.exchangeName, // exchange
		p.routingKey,   // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
}

func (p *ProducerMQ) KeepConnection() error {
	if err := p.connect(); err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	for {
		select {
		case err := <-p.done:
			if err != nil {
				err = p.reConnect()
				if err != nil {
					return fmt.Errorf("Reconnecting Error: %s", err)
				}
			} else {
				log.Println("Finishing ProducerMQ...")
				return nil
			}

			fmt.Println("Reconnected...")
		}
	}
}

func (p *ProducerMQ) GracefulStop() error {
	err := p.channel.Close()
	if err != nil {
		return err
	}
	err = p.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
