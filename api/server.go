package api

import (
	"context"
	"fmt"
	"user-service/config"
	"user-service/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type ServerBuilder struct {
	router chi.Router
	server server.Server
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
	}
}

func (s *ServerBuilder) Build() server.Server {
	s.server.UseHandler(s.router)

	return s.server
}
