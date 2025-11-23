package app

import (
	"AvitoTask2025/config"
	"AvitoTask2025/db"
	"AvitoTask2025/generated/api/pr_service"
	"AvitoTask2025/internal/controller"
	usecase "AvitoTask2025/internal/usecase/pr_service"
	"AvitoTask2025/internal/usecase/repository"
	"context"
	"errors"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func Run(logger *zap.Logger, cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.PG.URL)
	if err != nil {
		logger.Error("can not create pgxpool", zap.Error(err))
		return
	}
	defer dbPool.Close()
	db.SetupPostgres(dbPool, logger)

	repo := repository.New(dbPool)
	useCase := usecase.New(logger, repo, repo)
	cntr := controller.NewPrServiceController(logger, useCase, useCase, useCase)
	r := http.NewServeMux()
	h := pr_service.HandlerFromMux(pr_service.NewStrictHandler(
		cntr,
		[]pr_service.StrictMiddlewareFunc{},
	), r)

	s := &http.Server{
		Handler: h,
		Addr:    "0.0.0.0:" + cfg.Server.Port,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	logger.Info("Starting server on :" + cfg.Server.Port)
	errCh := make(chan error, 1)
	go func() {
		if err = s.ListenAndServe(); !errors.Is(http.ErrServerClosed, err) && err != nil {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	<-ctx.Done()
	logger.Info("Shutting down server...")
	if err = s.Shutdown(context.Background()); err != nil {
		logger.Error("Server shutdown error", zap.Error(err))
	}
	if err = <-errCh; err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
	logger.Info("Server stopped")
}
