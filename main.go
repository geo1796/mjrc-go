package main

import (
	"context"
	"errors"
	"mjrc/core/env"
	"mjrc/core/logger"
	"mjrc/core/postgres"
	"mjrc/core/runtime"
	"mjrc/core/security"
	"mjrc/web/api"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	os.Exit(run())
}

func run() int {
	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	environ, err := env.Load()
	if err != nil {
		logger.Error("failed to load env", logger.Err(err))
		return 1
	}

	if environ.IsProd() {
		logger.SetGlobalLogger(logger.New(logger.LevelInfo))
	}

	db, err := postgres.New(rootCtx, environ.PostgresConfig().DSN,
		environ.PostgresConfig().ConnMaxLifetime, environ.PostgresConfig().ConnMaxIdleTime,
		environ.PostgresConfig().MaxOpenConns, environ.PostgresConfig().MaxIdleConns)
	if err != nil {
		logger.Error("failed to create database", logger.Err(err))
		return 1
	}
	defer db.Close()

	deps := runtime.NewBuilder().
		WithDB(db).
		WithJWT(security.NewJWT(
			environ.SecurityConfig().AccessCookieName,
			environ.SecurityConfig().AccessTokenSecret,
			environ.SecurityConfig().AccessTokenTTL)).
		WithAdminAuthenticator(security.NewAuthenticator(environ.SecurityConfig().AdminPassword)).
		WithAPIKeyAuthenticator(security.NewAuthenticator(environ.SecurityConfig().APIKey)).
		Build()

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Get("/livez", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	})

	router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := db.Pool().Ping(ctx); err != nil {
			logger.Warn("healthz: db ping failed", logger.Err(err))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	api.Register(router, deps)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		<-rootCtx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Error("error shutting down server", logger.Err(err))
		}
	}()

	logger.Info("server started")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server error", logger.Err(err))
		return 1
	}

	return 0
}
