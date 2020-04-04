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
	if err != nil {
		log.Fatal("cannot read config ", err)
	}
	fmt.Println(cfg)

	logCfg := zap.NewDevelopmentConfig()
	logCfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logCfg.EncoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	logCfg.OutputPaths = []string{cfg.LogFile}

	logger, err := logCfg.Build()
	if err != nil {
		log.Fatal("cant create logger ", err)
	}
	defer logger.Sync()

	sugaredLogger := logger.Sugar()

	lis, err := net.Listen("tcp", cfg.HTTPListen)
	if err != nil {
		log.Fatalf("cannot listen %s, %v", cfg.HTTPListen, err)
	}

	db, err := sqlx.Connect("pgx", fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PgUser,
		cfg.PgPassword,
		cfg.PgHost,
		cfg.PgPort,
		cfg.PgName,
	))
	if err != nil {
		log.Fatalf("cannot connect to db %v", err)
	}

	storage, err := pg.NewStoragePg(db)
	if err != nil {
		log.Fatalf("cannot create storage %v", err)
	}

	app, err := app.NewCalendar(storage)
	if err != nil {
		log.Fatalf("cannot create app instance, %v", err)
	}
	err = app.RunScheduler(context.Background())
	if err != nil {
		log.Fatal("cant start scheduler ", err)
	}

	eventService := service.NewEventService(app, sugaredLogger)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	api.RegisterEventsServer(grpcServer, eventService)
	_ = grpcServer.Serve(lis)
}
