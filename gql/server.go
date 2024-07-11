package gql

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cosmossdk.io/log"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
)

// Server configures and starts the GQL server.
func Server(ctx context.Context, clientCtx client.Context, logger log.Logger) error {
	if !viper.GetBool("gql-server") {
		return nil
	}

	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		Debug:          false,
	}).Handler)

	logFile := viper.GetString("log-file")

	port := viper.GetString("gql-port")

	srv := handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: &Resolver{
		ctx:     clientCtx,
		logFile: logFile,
	}}))

	router.Handle("/", PlaygroundHandler("/api"))

	if viper.GetBool("gql-playground") {
		apiBase := viper.GetString("gql-playground-api-base")

		router.Handle("/webui", PlaygroundHandler(apiBase+"/api"))
		router.Handle("/console", PlaygroundHandler(apiBase+"/graphql"))
	}

	router.Handle("/api", srv)
	router.Handle("/graphql", srv)

	errCh := make(chan error)

	go func() {
		logger.Info(fmt.Sprintf("Connect to GraphQL playground url: http://localhost:%s", port))
		server := &http.Server{
			Addr:              ":" + port,
			Handler:           router,
			ReadHeaderTimeout: 3 * time.Second,
		}
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		// Gracefully stop the GQL server.
		logger.Info("Stopping GQL server...")
		return nil
	case err := <-errCh:
		logger.Error(fmt.Sprintf("Failed to start GQL server: %s", err))
		return err
	}
}
