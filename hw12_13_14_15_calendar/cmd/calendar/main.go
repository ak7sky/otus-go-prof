package main

import (
	"context"
	"flag"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/config"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/server/http"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.toml", "Path to yml configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	logg := logger.NewPreconfigured()

	conf, err := config.NewConfig(configFile)
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}

	logg = logger.NewConfigured(conf.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	strg, err := storage.New(ctx, conf.Storage)
	if err != nil {
		logg.Error(err.Error())
		os.Exit(1)
	}

	calendar := app.New(logg, strg)

	server := internalhttp.NewServer(logg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
