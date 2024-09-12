package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"user-service/api"
	"user-service/config"
	"user-service/db"
	"user-service/server"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	databaseTimeout = 15 * time.Second
)

type App struct {
	ctx context.Context
	log *zap.Logger

	settings config.Settings

	postgres *sqlx.DB

	server server.Server
}

func NewApp(ctx context.Context, log *zap.Logger, settings config.Settings) *App {
	return &App{
		ctx:      ctx,
		log:      log,
		settings: settings,
	}
}

func (a *App) InitDatabases() error {
	var err error

	postgresCtx, cancelPostgresCtx := context.WithTimeout(a.ctx, databaseTimeout)
	defer cancelPostgresCtx()

	a.postgres, err = db.NewPgx(postgresCtx, a.settings.Database.Postgres)
	if err != nil {
		return fmt.Errorf("could not connect to postgres: %w", err)
	}

	rootFS := os.DirFS("./")
	migrationPath := "db/migrations/postgres"
	err = db.Migrate(rootFS, a.log, a.postgres, migrationPath)
	if err != nil {
		return fmt.Errorf("could not migrate postgres: %w", err)
	}

	return nil
}

func (a *App) InitServices() {
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.ctx, a.log, a.settings)
	a.server = sb.Build()
}

func (a *App) Start() {
	a.server.Start()
}

func (a *App) Stop() {
	a.server.Stop()

	if err := a.postgres.Close(); err != nil {
		a.log.Error("could not close postgres connection", zap.Error(err))
	}
}
