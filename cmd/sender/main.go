package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	flag "github.com/spf13/pflag"

	"github.com/bobrovka/calendar/internal"
	"github.com/bobrovka/calendar/internal/models"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var configPath string

func init() {
	flag.StringVarP(&configPath, "config", "c", "", "path to config file")
}

func main() {
	flag.Parse()

	if configPath == "" {
		log.Fatal("no config file")
	}

	cfg := internal.Config{
		LogLevel: "debug",
	}

	loader := confita.NewLoader(
		file.NewBackend(configPath),
	)

	err := loader.Load(context.Background(), &cfg)
	failOnError(err, "cannot read config")
	fmt.Println(cfg)

	logCfg := zap.NewDevelopmentConfig()
	logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logCfg.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	logCfg.OutputPaths = []string{cfg.LogFileSender}

	logger, err := logCfg.Build()
	failOnError(err, "cant create logger")
	defer logger.Sync()

	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		cfg.RabbitUser,
		cfg.RabbitPassword,
		cfg.RabbitHost,
		cfg.RabbitPort,
	))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			var e models.Event
			json.Unmarshal(d.Body, &e)
			logger.Info(fmt.Sprintf("Notification to %s\n%s at %v", e.User, e.Title, e.StartAt))
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
