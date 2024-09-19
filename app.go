package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"user-service/api"
	"user-service/config"
	"user-service/db"
	dbuser "user-service/db/user"
	"user-service/kafka"
	"user-service/server"
	"user-service/service"
	"user-service/service/user"

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

	server      server.Server
	userService service.User
	kafka       kafka.Kafka
	consumer    kafka.Consumer
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

func (a *App) InitServices() error {
	var err error

	a.kafka = kafka.NewKafka(a.settings.Kafka.Brokers)
	a.consumer, err = a.kafka.Consumer(a.log, func() (context.Context, context.CancelFunc) {
		return context.WithCancel(a.ctx)
	}, kafka.WithTopic(a.settings.Kafka.Topics.UserTickets))
	if err != nil {
		return fmt.Errorf("could not create kafka consumer: %w", err)
	}

	userRepository := dbuser.NewRepository(a.postgres)

	a.userService = user.NewService(userRepository)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.ctx, a.log, a.settings)
	sb.AddUser(a.userService)
	a.server = sb.Build()
}

func (a *App) Start() {
	a.server.Start()
	a.consumer.Subscribe(a.userService.CreateSubscriberForBookMessage(a.ctx, a.log))
}

func (a *App) Stop(ctx context.Context) {
	a.server.Stop()

	if err := a.consumer.Close(ctx); err != nil {
		a.log.Error("could not close kafka consumer", zap.Error(err))
	}

	if err := a.postgres.Close(); err != nil {
		a.log.Error("could not close postgres connection", zap.Error(err))
	}
}
