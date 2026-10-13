// cmd/app/main.go (обновленный фрагмент)
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	deliveryHTTP "calories-checker/internal/delivery/http"
	"calories-checker/internal/delivery/http/middleware"
	"calories-checker/internal/repository"
	"calories-checker/internal/usecase"
	"calories-checker/pkg/config"
	"calories-checker/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	// 1. Инициализация структурированного логгера
	appLogger := logger.Init(os.Getenv("APP_ENV"))
	slog.SetDefault(appLogger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. БД
	dbPool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Error("unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		appLogger.Error("database ping failed", "error", err)
		os.Exit(1)
	}

	// 3. DI (Repo -> Usecase -> Handler)
	foodRepo := repository.NewPostgresFoodRepository(dbPool)
	waterRepo := repository.NewPostgresWaterRepository(dbPool)
	bodyMetricRepo := repository.NewPostgresBodyMetricRepository(dbPool)
	geminiClient := repository.NewGeminiClient(cfg.GeminiAPIKey)

	foodUC := usecase.NewFoodUsecase(foodRepo)
	waterUC := usecase.NewWaterUsecase(waterRepo)
	bodyMetricUC := usecase.NewBodyMetricUsecase(bodyMetricRepo)

	foodHandler := deliveryHTTP.NewFoodHandler(foodUC, geminiClient)
	healthHandler := deliveryHTTP.NewHealthHandler(waterUC, bodyMetricUC)

	// 4. Роутер + Middleware
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(appLogger)) // Подключаем логгер запросов

	apiGroup := router.Group("/api/v1")
	foodHandler.RegisterRoutes(apiGroup)
	healthHandler.RegisterRoutes(apiGroup)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		appLogger.Info("Server starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("listen error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	appLogger.Info("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(
    context.Background(),
    5*time.Second,
)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		appLogger.Error("Server forced to shutdown", "error", err)
	}

	appLogger.Info("Server exiting")
}
