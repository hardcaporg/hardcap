package main

import (
    "context"
    "errors"
    "fmt"
    "github.com/go-chi/chi/v5"
    "github.com/hardcaporg/hardcap/internal/config"
    "github.com/hardcaporg/hardcap/internal/db"
    "github.com/hardcaporg/hardcap/internal/logging"
    "github.com/hardcaporg/hardcap/internal/middleware"
    "github.com/hardcaporg/hardcap/internal/random"
    "github.com/hardcaporg/hardcap/internal/rpc/server"
    "github.com/hardcaporg/hardcap/internal/srv"
    "github.com/hardcaporg/hardcap/internal/version"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/rs/zerolog/log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
)

func main() {
    ctx := context.Background()
	random.SeedGlobal()
	config.Initialize("config/agent.env")
	logging.Initialize()
    config.PrintConfig(ctx)

	db.Initialize()

    err := server.Initialize(ctx)
    if err != nil {
        log.Fatal().Err(err).Msg("Cannot initialize RPC server")
        return
    }

	rootRouter := chi.NewRouter()
	rootRouter.Use(middleware.NewPatternMiddleware(version.Hostname))
	rootRouter.Use(middleware.VersionMiddleware)
	rootRouter.Use(middleware.TraceID)
	rootRouter.Use(middleware.LoggerMiddleware(&log.Logger))

	// Routes
	tmplRouter := chi.NewRouter()
	tmplRouter.Route("/ks", func(r chi.Router) {
		r.Get("/", srv.KickstartTemplateService)
	})
	restRouter := chi.NewRouter()
	restRouter.Route("/host_register", func(r chi.Router) {
		r.Post("/", srv.RegisterHostService)
	})
	rootRouter.Mount("/t", tmplRouter)
	rootRouter.Mount("/r", restRouter)

	// Routes for metrics
	metricsRouter := chi.NewRouter()
	metricsRouter.Handle(config.Prometheus.Path, promhttp.Handler())

	log.Info().Msgf("Starting new %s instance with prometheus on %d", config.Application.HttpListenAddress, config.Prometheus.Port)
	apiServer := http.Server{
		Addr:    config.Application.HttpListenAddress,
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
