package reminder

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

type bookingService interface {
	GetBookingsForReminder(ctx context.Context) ([]domain.Booking, error)
	MarkReminderSent(ctx context.Context, bookingID int64) error
}

type Worker struct {
	bookingService bookingService
	bot            *bot.Bot
	interval       time.Duration
	logger         *slog.Logger
}

func NewWorker(
	bookingService bookingService,
	b *bot.Bot,
	interval time.Duration,
	logger *slog.Logger,
) *Worker {
	return &Worker{
		bookingService: bookingService,
		bot:            b,
		interval:       interval,
		logger:         logger,
	}
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.sendReminders(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.sendReminders(ctx)
		}
	}
}

func (w *Worker) sendReminders(ctx context.Context) {
	bookings, err := w.bookingService.GetBookingsForReminder(ctx)
	if err != nil {
		w.logger.Error("failed to get bookings for reminder", "error", err)
		return
	}

	for _, booking := range bookings {
		if err := w.sendReminder(ctx, booking); err != nil {
			w.logger.Error(
				"failed to send booking reminder",
				"booking_id", booking.ID,
				"error", err,
			)
			continue
		}

		if err := w.bookingService.MarkReminderSent(ctx, booking.ID); err != nil {
			w.logger.Error(
				"failed to mark reminder sent",
				"booking_id", booking.ID,
				"error", err,
			)
		}
	}
}

func (w *Worker) sendReminder(ctx context.Context, booking domain.Booking) error {
	_, err := w.bot.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID: booking.TelegramUserId,
			Text:   reminderText(booking),
		},
	)
	return err
}

func reminderText(booking domain.Booking) string {
	return fmt.Sprintf(
		"🔔 Напоминание о записи\n\n✂️ Услуга: %s\n💈 Мастер: %s\n📅 Дата: %s\n🕒 Время: %s\n\nЖдём вас!",
		booking.ServiceName,
		booking.BarberName,
		tgui.FullDate(booking.StartsAt),
		tgui.TimeInterval(booking.StartsAt, booking.EndsAt),
	)
}
