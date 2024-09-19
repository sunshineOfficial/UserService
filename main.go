package main

import (
	"context"
	"user-service/config"
	"user-service/os"

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

	log.Debug("Service up")

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

	if err = app.InitServices(); err != nil {
		log.Error("Failed to init services", zap.Error(err))
		return
	}

	app.InitServer()

	app.Start()

	os.WaitTerminate(mainCtx, app.Stop)

	log.Debug("Service down")
}

func configureDecimal() {
	decimal.DivisionPrecision = 2
	decimal.MarshalJSONWithoutQuotes = true
}
