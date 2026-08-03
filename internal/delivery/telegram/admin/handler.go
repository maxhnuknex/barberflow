package admin

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
	service BookingService
}

func NewHandler(service BookingService) *Handler {
	return &Handler{
		service: service,
	}
}

type BookingService interface {
	Cancel(ctx context.Context, bookingID int64) error
	GetByID(ctx context.Context, bookingID int64) (domain.Booking, error)
	ListByDate(ctx context.Context, date time.Time) ([]domain.Booking, error)
	NextSevenDays() []time.Time
}

func (h *Handler) adminStart(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	keyboard := keyboardAdminStart()

	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      chatID(update),
			Text:        "Админ-панель",
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingsMenu(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "Записи",
			ReplyMarkup: keyboardBookingsMenu(),
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingsToday(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	bookings, err := h.service.ListByDate(ctx, time.Now().In(moscowLocation()))
	if err != nil {
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        bookingsListText(bookings),
			ReplyMarkup: keyboardBookingsList(bookings),
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingsDateMenu(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "Записи на дату",
			ReplyMarkup: keyboardBookingsDate(h.service.NextSevenDays()),
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingsDate(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	const prefix = "admin:bookings:date:"

	rawDate := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	selectedDate, err := time.ParseInLocation("2006-01-02", rawDate, moscowLocation())
	if err != nil {
		return
	}

	bookings, err := h.service.ListByDate(ctx, selectedDate)
	if err != nil {
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        bookingsListText(bookings),
			ReplyMarkup: keyboardBookingsList(bookings),
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingDetail(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	const prefix = "admin:booking:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	bookingID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return
	}

	booking, err := h.service.GetByID(ctx, bookingID)
	if err != nil {
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        bookingDetailText(booking),
			ReplyMarkup: keyboardBookingDetail(booking.ID),
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingCancel(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	const prefix = "admin:booking:cancel:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	bookingID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return
	}

	if err := h.service.Cancel(ctx, bookingID); err != nil {
		return
	}

	_, err = b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "Запись отменена",
			ReplyMarkup: keyboardBookingsMenu(),
		},
	)
	if err != nil {
		return
	}
}

func (h *Handler) handleBookingsFind(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
			Text:        "В разработке",
			ReplyMarkup: keyboardBookingsMenu(),
		},
	)
	if err != nil {
		return
	}
}

func bookingsListText(bookings []domain.Booking) string {
	if len(bookings) == 0 {
		return "Записей нет"
	}

	return "Записи"
}

func bookingDetailText(booking domain.Booking) string {
	return fmt.Sprintf(
		"ID: %d\nКлиент: %s\nМастер: %s\nУслуга: %s\nДата: %s\nВремя: %s-%s\nЦена: %d ₽",
		booking.ID,
		customerText(booking),
		booking.BarberName,
		booking.ServiceName,
		booking.StartsAt.Format("02.01.2006"),
		booking.StartsAt.Format("15:04"),
		booking.EndsAt.Format("15:04"),
		booking.PriceMinorUnits/100,
	)
}

func customerText(booking domain.Booking) string {
	if booking.CustomerUsername != "" {
		return fmt.Sprintf("%s (@%s)", booking.CustomerName, booking.CustomerUsername)
	}

	return booking.CustomerName
}

func moscowLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.FixedZone("Europe/Moscow", 3*60*60)
	}

	return loc
}

func chatID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	return update.CallbackQuery.Message.Message.Chat.ID
}
