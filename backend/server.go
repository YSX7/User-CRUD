package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"golang.org/x/net/context"
	"vuegolang/graph"
	"vuegolang/pkg/middleware"
	"vuegolang/pkg/sessions"
)

const defaultPort = "8080"

func main() {
	var ()

	log.SetFlags(log.Lshortfile)

	db := InitDb()
	defer db.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(
		cors.New(
			cors.Options{
				AllowedOrigins:      []string{"http://localhost:9000"},
				AllowCredentials:    true,
				AllowPrivateNetwork: true,
				AllowedHeaders: []string{
					"Accept",
					"Content-Type",
					"Content-Length",
					"Accept-Encoding",
					"Authorization",
					"Set-Cookie",
				},
				AllowedMethods: []string{
					"GET", "POST", "PUT", "DELETE", "UPDATE", "PATCH",
				},
				Debug: true,
				ExposedHeaders: []string{
					"Set-Cookie",
				},
			},
		).Handler,
	)
	router.Use(middleware.AuthMiddleware(middleware.HandleAuthHTTP))

	resolver := &graph.Resolver{
		Db:       db,
		Sessions: sessions.New(),
	}

	graphConfig := graph.Config{
		Resolvers: resolver,
	}
	graphConfig.Directives.Auth = resolver.Auth
	srv := CreateServer(graphConfig)

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	httpSrv := http.Server{Addr: ":" + port, Handler: router}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		log.Print("Shutting down application...")

		// We received an interrupt signal, shut down.
		if errShutdown := httpSrv.Shutdown(context.Background()); errShutdown != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", errShutdown)
		}
		close(idleConnsClosed)
	}()

	if errListen := httpSrv.ListenAndServe(); errListen != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", errListen)
	}

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, router))

	// Block current goroutine until shutdown
	<-idleConnsClosed

	// Received OS signal: syscall.SIGTERM or syscall.SIGINT
	log.Print("Shutting down application...")
}

func CreateServer(config graph.Config) *handler.Server {
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(config),
	)
	srv.AddTransport(
		&transport.Websocket{
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					// Check against your desired domains here
					return r.Host == "localhost"
				},
				ReadBufferSize:  1024,
				WriteBufferSize: 1024,
			},
		},
	)
	return srv
}
