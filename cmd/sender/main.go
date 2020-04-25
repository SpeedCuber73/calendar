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

	"github.com/bobrovka/calendar/internal"
	"github.com/bobrovka/calendar/internal/consumer"
	"github.com/bobrovka/calendar/internal/models"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	flag "github.com/spf13/pflag"
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

	mqConsumer := consumer.NewConsumer("", fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		cfg.RabbitUser,
		cfg.RabbitPassword,
		cfg.RabbitHost,
		cfg.RabbitPort,
	), "event.exchange", "direct", "event.queue", "event.notification")

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	go func() {
		err := mqConsumer.Handle(ctx, wg, func(msgs <-chan amqp.Delivery) {
			for {
				select {
				case msg, ok := <-msgs:
					// если закроется канал, завершить обработчик
					if !ok {
						return
					}
					var e models.Event
					err := json.Unmarshal(msg.Body, &e)
					if err != nil {
						logger.Warn(fmt.Sprintf("got invalid message %s", msg.Body))
						msg.Reject(false)
					} else {
						logger.Info(fmt.Sprintf("Notification to %s\n%s at %v", e.User, e.Title, e.StartAt))
						msg.Ack(false)
					}
				case <-ctx.Done():
					// если завершается программа, завершить обработчик
					return
				}
			}
		}, 3)
		failOnError(err, "handling error")
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
