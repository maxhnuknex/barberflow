package mybookings

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

type Handler struct {
	service BookingService
}

func NewHandler(service BookingService) *Handler {
	return &Handler{
		service: service,
	}
}

type BookingService interface {
	ListMyBooking(
		ctx context.Context,
		telegramUserID int64,
	) ([]domain.Booking, error)
}

func (h *Handler) HandlerListMyBooking(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	telegramId := update.CallbackQuery.From.ID

	bookings, err := h.service.ListMyBooking(ctx, telegramId)
	if err != nil {
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        myBookingsText(bookings),
			ReplyMarkup: keyboardListMyBooking(bookings),
		},
	)
	if err != nil {
		return
	}
}
