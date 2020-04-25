package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bobrovka/calendar/internal"
	app "github.com/bobrovka/calendar/internal/calendar-app"
	"github.com/bobrovka/calendar/internal/scheduler"
	"github.com/bobrovka/calendar/internal/scheduler/producer"
	"github.com/bobrovka/calendar/internal/service"
	pg "github.com/bobrovka/calendar/internal/storage/storage-pg"
	"github.com/bobrovka/calendar/pkg/calendar/api"
	"github.com/go-errors/errors"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

var ErrOSTerminated = errors.New("os terminated")

func init() {
	flag.StringVarP(&configPath, "config", "c", "", "path to config file")
}

func main() {
	flag.Parse()

	cfg := getConfig()

	logCfg := zap.NewDevelopmentConfig()
	logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logCfg.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	logCfg.OutputPaths = []string{cfg.LogFile}

	logger, err := logCfg.Build()
	failOnError(err, "cant create logger")
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	storage, err := pg.NewStoragePg(cfg.PgUser, cfg.PgPassword, cfg.PgHost, cfg.PgPort, cfg.PgName)
	failOnError(err, "cannot create storage")

	app, err := app.NewCalendar(storage, sugaredLogger)
	failOnError(err, "cannot create app instance")

	eventService := service.NewEventService(app, sugaredLogger)

	// Create grpc server
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	api.RegisterEventsServer(grpcServer, eventService)

	lis, err := net.Listen("tcp", cfg.HTTPListen)
	failOnError(err, fmt.Sprint("cannot listen ", cfg.HTTPListen))

	exitChannel := make(chan error)
	go func() {
		// start grpc server
		exitChannel <- grpcServer.Serve(lis)
	}()

	go func() {
		termChan := make(chan os.Signal)
		signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

		<-termChan
		exitChannel <- ErrOSTerminated
	}()

	producer := producer.NewProducerMQ(fmt.Sprintf(
		"amqp://%s:%s@%s:%d",
		cfg.RabbitUser,
		cfg.RabbitPassword,
		cfg.RabbitHost,
		cfg.RabbitPort,
	), "event.exchange", "direct", "event.queue", "event.notification")
	sched := scheduler.NewScheduler(producer, storage, sugaredLogger)

	go func() {
		exitChannel <- sched.Run()
	}()

	err = <-exitChannel
	log.Println("stopped with err: ", err)

	grpcServer.GracefulStop()
	err = sched.Stop()
	if err != nil {
		log.Println("cannot gracefully stop scheduler, err: ", err)
	}
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
