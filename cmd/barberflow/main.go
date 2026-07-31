package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxhnucknex/barberflow/internal/telegram/booking"
	"github.com/maxhnucknex/barberflow/internal/telegram/start"
	"github.com/maxhnucknex/barberflow/repository/postgres"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	//OS
	token := os.Getenv("TG_TOKEN")
	if token == "" {
		logger.Error("TELEGRAM_BOT_TOKEN is not set")
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		logger.Error("DATABASE_URL is not set")
		os.Exit(1)
	}

	//CTX
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
	)
	defer stop()

	// db start

	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		logger.Error(
			"failed to create poolDB",
			"error", err,
		)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		logger.Error("failed to ping postgres", "error", err)
		os.Exit(1)
	}
	logger.Info("postgres connected")

	//TG BOT
	b, err := bot.New(token)
	if err != nil {
		logger.Error("Fail to create bot", "error", err)
		os.Exit(1)
	}

	startHandler := start.NewHandler()
	start.RegisterRoutes(b, startHandler)

	//repository
	serviceRepository := postgres.NewServiceRepository(db)
	barberRepository := postgres.NewBarberRepository(db)
	bookingRepository := postgres.NewBookingRepository(db)

	bookingService := booking.NewService(serviceRepository, barberRepository, bookingRepository)
	bookingHandler := booking.NewBookingHandler(bookingService)
	booking.RegisterHandler(b, bookingHandler)

	//start

	b.Start(ctx)
}
