package app

import (
	"AggregationService/internal/config"
	"AggregationService/internal/pkg/logger"
	sloglogger "AggregationService/internal/pkg/logger/slog-logger"
	middleware2 "AggregationService/internal/pkg/middleware"
	"context"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const environment = "ENV"

type App struct {
	httpServer *http.Server
	provider   *Provider
}

func New(ctx context.Context, cfg *config.Config, provider *Provider) *App {
	subHandler := provider.handler
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware2.LoggerMW)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))
	r.Use(middleware2.HeadersMiddleware)

	// --- ROUTES ---
	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", subHandler.Create)
		r.Get("/", subHandler.GetAll)
		r.Get("/cost", subHandler.CalculateCost)
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", subHandler.GetByID)
			r.Put("/", subHandler.Update)
			r.Delete("/", subHandler.Delete)
		})
	})

	srv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: r,
	}

	return &App{httpServer: srv}
}

func (a *App) Run(ctx context.Context) error {
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = a.httpServer.Shutdown(shutdownCtx)
		close(idleConnsClosed)
	}()

	fmt.Printf("Server started at %s\n", a.httpServer.Addr)
	if err := a.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	<-idleConnsClosed
	return nil
}

func InitContextWithLogger(ctx context.Context) context.Context {
	env := os.Getenv(environment)
	return logger.ContextWithLogger(ctx, sloglogger.New(env))
}
