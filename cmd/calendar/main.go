package main

import (
	"context"
	"fmt"
	"log"
	"net"

	app "github.com/bobrovka/calendar/internal/calendar-app"
	pg "github.com/bobrovka/calendar/internal/calendar-app/storage-pg"
	"github.com/bobrovka/calendar/internal/service"
	"github.com/bobrovka/calendar/pkg/calendar/api"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	flag "github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	cfg := app.Config{
		HTTPListen: "127.0.0.1:50051",
		LogLevel:   "debug",
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
	logCfg.OutputPaths = []string{cfg.LogFile}

	logger, err := logCfg.Build()
	failOnError(err, "cant create logger")
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	lis, err := net.Listen("tcp", cfg.HTTPListen)
	failOnError(err, fmt.Sprint("cannot listen ", cfg.HTTPListen))

	db, err := sqlx.Connect("pgx", fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgName,
	))
	failOnError(err, "cannot connect to db")

	storage, err := pg.NewStoragePg(db)
	failOnError(err, "cannot create storage")

	app, err := app.NewCalendar(storage, sugaredLogger)
	failOnError(err, "cannot create app instance")

	ch := getMQChannel()
	err = app.RunScheduler(context.Background(), ch)
	failOnError(err, "cant start scheduler")

	eventService := service.NewEventService(app, sugaredLogger)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	api.RegisterEventsServer(grpcServer, eventService)
	_ = grpcServer.Serve(lis)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "cant connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "cant get channel on RabbitMQ")

	return ch
}
