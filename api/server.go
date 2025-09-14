package api

import (
	"context"
	"fmt"
	"user-service/api/handlers"
	"user-service/config"
	_ "user-service/docs"
	"user-service/graph"
	"user-service/service"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/middleware"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/plugin"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/golog"
	swagger "github.com/swaggo/http-swagger/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

type ServerBuilder struct {
	server      goserver.Server
	router      *gorouter.Router
	graphServer *handler.Server
}

func NewServerBuilder(ctx context.Context, log golog.Logger, settings config.Settings, resolver *graph.Resolver) *ServerBuilder {
	graphServer := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	graphServer.AddTransport(transport.Options{})
	graphServer.AddTransport(transport.GET{})
	graphServer.AddTransport(transport.POST{})

	graphServer.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	graphServer.Use(extension.Introspection{})
	graphServer.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return &ServerBuilder{
		server: goserver.NewHTTPServer(ctx, log, fmt.Sprintf(":%d", settings.Port)),
		router: gorouter.NewRouter(log).Use(
			middleware.Metrics(),
			middleware.Recover,
			middleware.LogError,
		),
		graphServer: graphServer,
	}
}

func (s *ServerBuilder) AddDebug() {
	s.router.Install(plugin.NewPProf(), plugin.NewMetrics())
}

func (s *ServerBuilder) AddSwagger() {
	s.router.PathPrefix("/swagger/").Handler(gorouter.WrapStdLibFunc(swagger.Handler(swagger.URL("/swagger/doc.json"))))
}

func (s *ServerBuilder) AddUser(user service.User) {
	s.router.HandleGet("/user/{id}", handlers.GetUserByIdHandler(user))
	s.router.HandleGet("/user", handlers.GetUsersHandler(user))
	s.router.HandlePost("/user", handlers.AddUserHandler(user))
	s.router.HandlePut("/user/{id}", handlers.UpdateUserHandler(user))
	s.router.HandleDelete("/user/{id}", handlers.DeleteUserHandler(user))
	s.router.HandleGet("/user/{id}/tickets", handlers.GetUserTicketsByUserIdHandler(user))
}

func (s *ServerBuilder) AddGraphQL() {
	s.router.Path("/").Handler(gorouter.WrapStdLibFunc(playground.Handler("GraphQL playground", "/query")))
	s.router.Path("/query").Handler(gorouter.WrapStdLib(s.graphServer))
}

func (s *ServerBuilder) Build() goserver.Server {
	s.server.UseHandler(s.router)

	return s.server
}
