package booking

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type ServiceRepository interface {
	ListActive(ctx context.Context) ([]domain.Service, error)
}

type BarberRepository interface {
	ListBarberByService(ctx context.Context, id int64) ([]domain.Barber, error)
}

type Handler struct {
	serviceRepo ServiceRepository
	barberRepo  BarberRepository
}

func NewHandler(
	serviceRepo ServiceRepository,
	barberRepo BarberRepository,
) *Handler {
	return &Handler{
		serviceRepo: serviceRepo,
		barberRepo:  barberRepo,
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

func (h *Handler) HandlerListBarber(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	data := update.CallbackQuery.Data
	const prefix = "booking:service:"
	idStr := strings.TrimPrefix(data, prefix)
	serviceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return
	}

	barbers, err := h.barberRepo.ListBarberByService(ctx, serviceID)
	if err != nil {
		return
	}

	keyboard := barberKeyboard(barbers)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Выберите мастера:",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return
	}
}
