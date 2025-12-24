package http

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/hoshina-dev/gapi/internal/adapters/graph"
	"github.com/hoshina-dev/gapi/internal/adapters/infrastructure"
	"github.com/vektah/gqlparser/v2/ast"
)

func SetupRouter(resolver *graph.Resolver, cfg infrastructure.Config) *fiber.App {
	app := fiber.New()

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CorsOrigins,
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/health", healthCheck)
	app.Get("/", playgroundHandler())
	app.All("/query", graphQLHandler(resolver))

	return app
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok", "time": time.Now()})
}

func graphQLHandler(resolver *graph.Resolver) fiber.Handler {
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return adaptor.HTTPHandler(srv)
}

func playgroundHandler() fiber.Handler {
	return adaptor.HTTPHandler(playground.Handler("GraphQL playground", "/query"))
}
