package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/go-telegram/bot"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	appbarber "github.com/maxhnucknex/barberflow/internal/app/barber"
	appbooking "github.com/maxhnucknex/barberflow/internal/app/booking"
	appcatalog "github.com/maxhnucknex/barberflow/internal/app/catalog"
	"github.com/maxhnucknex/barberflow/internal/delivery/telegram/admin"
	telegrambooking "github.com/maxhnucknex/barberflow/internal/delivery/telegram/booking"
	"github.com/maxhnucknex/barberflow/internal/delivery/telegram/mybookings"
	"github.com/maxhnucknex/barberflow/internal/delivery/telegram/start"
	"github.com/maxhnucknex/barberflow/internal/repository/postgres"
	"github.com/maxhnucknex/barberflow/internal/worker/reminder"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	_ = godotenv.Load()

	//OS
	token := os.Getenv("TG_TOKEN")
	if token == "" {
		token = os.Getenv("TELEGRAM_BOT_TOKEN")
	}
	if token == "" {
		logger.Error("TG_TOKEN or TELEGRAM_BOT_TOKEN is not set")
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

	bookingService := appbooking.NewService(serviceRepository, barberRepository, bookingRepository)
	catalogService := appcatalog.NewService(serviceRepository)
	barberService := appbarber.NewService(barberRepository)

	bookingHandler := telegrambooking.NewBookingHandler(bookingService)
	telegrambooking.RegisterHandler(b, bookingHandler)

	myBookingsHandler := mybookings.NewHandler(bookingService)
	mybookings.RegisterHandler(b, myBookingsHandler)

	adminHandler := admin.NewHandler(bookingService, catalogService, barberService)
	admin.RegisterHandler(b, adminHandler)

	reminderWorker := reminder.NewWorker(bookingService, b, time.Minute, logger)
	go reminderWorker.Run(ctx)

	//start

	b.Start(ctx)
}
