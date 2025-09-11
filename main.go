package main

import (
	"context"
	"user-service/config"

	"github.com/shopspring/decimal"
	"github.com/sunshineOfficial/golib/golog"
	"github.com/sunshineOfficial/golib/goos"
)

// @title		user-service API
// @version		1.0
// @description	Микросервис пользователей.
func main() {
	configureDecimal()

	log := golog.NewLogger("user-service")
	log.Debug("service up")

	mainCtx, cancelMainCtx := context.WithCancel(context.Background())
	defer cancelMainCtx()

	settings, err := config.Parse()
	if err != nil {
		log.Errorf("failed to load settings: %v", err)
		return
	}

	app := NewApp(mainCtx, log, settings)

	if err = app.InitDatabases(); err != nil {
		log.Errorf("failed to init databases: %v", err)
		return
	}

	if err = app.InitServices(); err != nil {
		log.Errorf("failed to init services: %v", err)
		return
	}

	app.InitServer()

	app.Start()

	goos.WaitTerminate(mainCtx, app.Stop)

	log.Debug("service down")
}

func configureDecimal() {
	decimal.DivisionPrecision = 2
	decimal.MarshalJSONWithoutQuotes = true
}
