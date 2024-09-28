package api

import (
	"context"
	"fmt"
	"user-service/api/handlers"
	"user-service/config"
	_ "user-service/docs"
	"user-service/server"
	"user-service/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	swagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type ServerBuilder struct {
	router chi.Router
	server server.Server
	log    *zap.Logger
}

func NewServerBuilder(ctx context.Context, log *zap.Logger, settings config.Settings) *ServerBuilder {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Heartbeat("/ping"))
	router.Mount("/debug", middleware.Profiler())

	return &ServerBuilder{
		router: router,
		server: server.NewHTTPServer(ctx, log, fmt.Sprintf(":%d", settings.Port)),
		log:    log,
	}
}

func (s *ServerBuilder) AddSwagger() {
	s.router.Get("/swagger/*", swagger.Handler(swagger.URL("/swagger/doc.json")))
}

func (s *ServerBuilder) AddUser(user service.User) {
	s.router.Get("/user/{id}", handlers.GetUserByIdHandler(user, s.log))
	s.router.Get("/user", handlers.GetUsersHandler(user, s.log))
	s.router.Post("/user", handlers.AddUserHandler(user, s.log))
	s.router.Put("/user/{id}", handlers.UpdateUserHandler(user, s.log))
	s.router.Delete("/user/{id}", handlers.DeleteUserHandler(user, s.log))
	s.router.Get("/user/{id}/tickets", handlers.GetUserTicketsByUserIdHandler(user, s.log))
}

func (s *ServerBuilder) Build() server.Server {
	s.server.UseHandler(s.router)

	return s.server
}
