package main

import (
	"context"
	"log"
	"net"

	app "github.com/bobrovka/calendar/internal/calendar-app"
	"github.com/bobrovka/calendar/internal/grpc/api"
	"github.com/bobrovka/calendar/internal/service"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
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
	_ = sugaredLogger

	lis, err := net.Listen("tcp", cfg.HTTPListen)
	if err != nil {
		log.Fatalf("cannot listen %s, %v", cfg.HTTPListen, err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	api.RegisterEventsServer(grpcServer, &service.EventService{})
	_ = grpcServer.Serve(lis)
}
