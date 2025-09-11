package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"user-service/api"
	"user-service/config"
	dbuser "user-service/db/user"
	"user-service/service"
	"user-service/service/user"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const (
	databaseTimeout = 15 * time.Second
)

type App struct {
	ctx context.Context
	log golog.Logger

	settings config.Settings

	postgres *sqlx.DB

	server      goserver.Server
	userService service.User
	kafka       gokafka.Kafka
	consumer    gokafka.Consumer
}

func NewApp(ctx context.Context, log golog.Logger, settings config.Settings) *App {
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

	a.kafka = gokafka.NewKafka(a.settings.Kafka.Brokers)
	a.consumer, err = a.kafka.Consumer(a.log, func() (context.Context, context.CancelFunc) {
		return context.WithCancel(a.ctx)
	}, gokafka.WithTopic(a.settings.Kafka.Topics.UserTickets))
	if err != nil {
		return fmt.Errorf("could not create kafka consumer: %w", err)
	}

	userRepository := dbuser.NewRepository(a.postgres)

	a.userService = user.NewService(userRepository)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.ctx, a.log, a.settings)
	sb.AddDebug()
	sb.AddSwagger()
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
		a.log.Errorf("could not close kafka consumer: %v", err)
	}

	if err := a.postgres.Close(); err != nil {
		a.log.Errorf("could not close postgres connection: %v", err)
	}
}
