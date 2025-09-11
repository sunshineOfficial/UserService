package api

import (
	"context"
	"fmt"
	"user-service/api/handlers"
	"user-service/config"
	_ "user-service/docs"
	"user-service/service"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/middleware"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/plugin"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/golog"
	swagger "github.com/swaggo/http-swagger/v2"
)

type ServerBuilder struct {
	server goserver.Server
	router *gorouter.Router
}

func NewServerBuilder(ctx context.Context, log golog.Logger, settings config.Settings) *ServerBuilder {
	return &ServerBuilder{
		server: goserver.NewHTTPServer(ctx, log, fmt.Sprintf(":%d", settings.Port)),
		router: gorouter.NewRouter(log).Use(
			middleware.Metrics(),
			middleware.Recover,
			middleware.LogError,
		),
	}
}

func (s *ServerBuilder) AddDebug() {
	s.router.Install(plugin.NewPProf(), plugin.NewMetrics())
}

func (s *ServerBuilder) AddSwagger() {
	s.router.HandleGet("/swagger/*", gorouter.WrapStdLibFunc(swagger.Handler(swagger.URL("/swagger/doc.json"))))
}

func (s *ServerBuilder) AddUser(user service.User) {
	s.router.HandleGet("/user/{id}", handlers.GetUserByIdHandler(user))
	s.router.HandleGet("/user", handlers.GetUsersHandler(user))
	s.router.HandlePost("/user", handlers.AddUserHandler(user))
	s.router.HandlePut("/user/{id}", handlers.UpdateUserHandler(user))
	s.router.HandleDelete("/user/{id}", handlers.DeleteUserHandler(user))
	s.router.HandleGet("/user/{id}/tickets", handlers.GetUserTicketsByUserIdHandler(user))
}

func (s *ServerBuilder) Build() goserver.Server {
	s.server.UseHandler(s.router)

	return s.server
}
