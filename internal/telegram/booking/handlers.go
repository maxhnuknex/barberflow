package booking

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

type Handler struct {
	service Service
}

func NewBookingHandler(
	service Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

type Service interface {
	ListActive(ctx context.Context) ([]domain.Service, error)
	ListBarberByService(ctx context.Context, id int64) ([]domain.Barber, error)
	ListTimeBarberActive(
		ctx context.Context,
		serviceID int64,
		barberID int64,
	) ([]domain.TimeInterval, error)
	PostBooking(
		ctx context.Context,
		telegramUserID int64,
		username string,
		firstName string,
		serviceID int64,
		barberID int64,
		time domain.TimeInterval,
	) (domain.Booking, error)
	ListMyBooking(
		ctx context.Context,
		telegramUserID int64,
	) ([]domain.Booking, error)
}

func (h *Handler) HandlerListActive(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	services, err := h.service.ListActive(ctx)

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

func (h *Handler) HandlerBookingTime(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	const prefix = "booking:time:"

	data := update.CallbackQuery.Data
	rawValues := strings.TrimPrefix(data, prefix)

	parts := strings.Split(rawValues, ":")
	if len(parts) != 4 {
		return
	}

	serviceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}

	barberID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}

	startUnix, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return
	}

	endUnix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		return
	}

	startsAt := time.Unix(startUnix, 0)
	endsAt := time.Unix(endUnix, 0)

	firstName := update.CallbackQuery.From.FirstName
	if firstName == "" {
		firstName = "Telegram user"
	}

	booking, err := h.service.PostBooking(
		ctx,
		update.CallbackQuery.From.ID,
		update.CallbackQuery.From.Username,
		firstName,
		serviceID,
		barberID,
		domain.TimeInterval{
			StartsAt: startsAt,
			EndsAt:   endsAt,
		},
	)
	if err != nil {
		return
	}

	text := fmt.Sprintf(
		"Запись создана:\n\nУслуга: %s\nМастер: %s\nВремя: %s-%s",
		booking.ServiceName,
		booking.BarberName,
		booking.StartsAt.Format("15:04"),
		booking.EndsAt.Format("15:04"),
	)

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        text,
			ReplyMarkup: mainMenuKeyboard(),
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

	barbers, err := h.service.ListBarberByService(ctx, serviceID)
	if err != nil {
		return
	}

	keyboard := barberKeyboard(serviceID, barbers)

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        "Выберите мастера:",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return
	}
}

func (h *Handler) HandlerListTimeBarberActive(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	const prefix = "booking:barber:"

	data := update.CallbackQuery.Data
	rawIDs := strings.TrimPrefix(data, prefix)

	parts := strings.Split(rawIDs, ":")
	if len(parts) != 2 {
		return
	}

	serviceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return
	}

	barberID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return
	}

	freeSlot, err := h.service.ListTimeBarberActive(ctx, serviceID, barberID)
	if err != nil {
		return
	}

	keyboard := freeSlotKeyboard(
		barberID,
		serviceID,
		freeSlot,
	)

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "выберите время записи",
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		return
	}
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

func myBookingsText(bookings []domain.Booking) string {
	if len(bookings) == 0 {
		return "У вас пока нет активных записей."
	}

	var builder strings.Builder
	builder.WriteString("Ваши записи:")

	for i, booking := range bookings {
		fmt.Fprintf(
			&builder,
			"\n\n%d. Услуга: %s\nМастер: %s\nВремя: %s-%s",
			i+1,
			booking.ServiceName,
			booking.BarberName,
			booking.StartsAt.Format("02.01.2006 15:04"),
			booking.EndsAt.Format("15:04"),
		)
	}

	return builder.String()
}
