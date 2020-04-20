package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	cfg := getConfig()

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

	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case msg := <-msgs:
				var e models.Event
				json.Unmarshal(msg.Body, &e)
				logger.Info(fmt.Sprintf("Notification to %s\n%s at %v", e.User, e.Title, e.StartAt))
				msg.Ack(false)
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan

	cancel()
	wg.Wait()

	fmt.Println("Sender stopped")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getConfig() *internal.Config {
	if configPath == "" {
		log.Fatal("no config file")
	}

	cfg := &internal.Config{
		HTTPListen: "127.0.0.1:50051",
		LogLevel:   "debug",
	}

	loader := confita.NewLoader(
		file.NewBackend(configPath),
	)

	err := loader.Load(context.Background(), cfg)
	failOnError(err, "cannot read config")
	fmt.Println(cfg)
	return cfg
}
