package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	app "github.com/bobrovka/calendar/internal/calendar-app"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	flag "github.com/spf13/pflag"
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

	cfg := app.Config{
		HTTPListen: "127.0.0.1:9000",
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

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		defer sugaredLogger.Sync()
		sugaredLogger.Infow("hello world handler", "Method", r.Method, "URL", r.URL)
		fmt.Fprintf(w, "Hello, world!")
	})
	http.ListenAndServe(cfg.HTTPListen, nil)
}
