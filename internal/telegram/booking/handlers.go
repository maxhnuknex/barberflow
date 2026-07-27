package booking

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type ServiceRepository interface {
	ListActive(ctx context.Context) ([]domain.Service, error)
}

type Handler struct {
	serviceRepo ServiceRepository
}

func NewHandler(serviceRepo ServiceRepository) *Handler {
	return &Handler{
		serviceRepo: serviceRepo,
	}
}

func (h *Handler) HandlerListActive(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	services, err := h.serviceRepo.ListActive(ctx)
	if err != nil {
		return
	}

	keyboard := serviceKeyboard(services)

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "Выберите услугу:",
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		return
	}
}
