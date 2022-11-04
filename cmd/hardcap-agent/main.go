package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/hardcaporg/hardcap/internal/config"
	"github.com/hardcaporg/hardcap/internal/logging"
	"github.com/hardcaporg/hardcap/internal/middleware"
	"github.com/hardcaporg/hardcap/internal/random"
	"github.com/hardcaporg/hardcap/internal/routes"
	"github.com/hardcaporg/hardcap/internal/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

func main() {
	// ctx := context.Background()
	random.SeedGlobal()
	config.Initialize("config/agent.env")
	logging.Initialize()

	rootRouter := chi.NewRouter()
	rootRouter.Use(middleware.NewPatternMiddleware(version.Hostname))
	rootRouter.Use(middleware.VersionMiddleware)
	rootRouter.Use(middleware.TraceID)
	rootRouter.Use(middleware.LoggerMiddleware(&log.Logger))

	tmplRouter := chi.NewRouter()

	// Mount paths
	routes.MountTemplateEndpoint(tmplRouter)
	rootRouter.Mount("/t", tmplRouter)

	// Routes for metrics
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle(config.Prometheus.Path, promhttp.Handler())

	log.Info().Msgf("Starting new instance on port %d with prometheus on %d", config.Application.Port, config.Prometheus.Port)
	apiServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Application.Port),
		Handler: rootRouter,
	}

	metricsServer := http.Server{
		Addr:    fmt.Sprintf(":%d", config.Prometheus.Port),
		Handler: metricsRouter,
	}

	waitForSignal := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint
		if err := apiServer.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Main service shutdown error")
		}
		if err := metricsServer.Shutdown(context.Background()); err != nil {
			log.Fatal().Err(err).Msg("Metrics service shutdown error")
		}
		close(waitForSignal)
	}()

	go func() {
		if err := metricsServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatal().Err(err).Msg("Metrics service listen error")
			}
		}
	}()

	if err := apiServer.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Main service listen error")
		}
	}

	<-waitForSignal

	log.Info().Msg("Shutdown finished, exiting")
}
