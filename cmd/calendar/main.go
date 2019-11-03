package main

import (
	"context"
	"fmt"
	"log"

	app "github.com/SpeedCuber73/calendar/internal/calendar-app"
	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
	flag "github.com/spf13/pflag"
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
	fmt.Println(cfg)
}
