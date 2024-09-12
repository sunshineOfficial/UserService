package main

import (
	"context"
	"os"
	"os/signal"
	"user-service/config"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func main() {
	configureDecimal()

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	defer cancelMainCtx()

	log, err := zap.NewDevelopment(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		panic(err)
	}

	settings, err := config.NewSettings()
	if err != nil {
		log.Error("Failed to load settings", zap.Error(err))
		return
	}

	app := NewApp(mainCtx, log, settings)

	if err = app.InitDatabases(); err != nil {
		log.Error("Failed to init databases", zap.Error(err))
		return
	}

	app.InitServices()
	app.InitServer()

	app.Start()
	stop := getStopSignal()
	<-stop

	app.Stop()
}

func configureDecimal() {
	decimal.DivisionPrecision = 2
	decimal.MarshalJSONWithoutQuotes = true
}

func getStopSignal() <-chan os.Signal {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	return stop
}
