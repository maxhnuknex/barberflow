package mybookings

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
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
	tgui.AnswerCallbackQuery(ctx, b, update)

	telegramId := update.CallbackQuery.From.ID

	bookings, err := h.service.ListMyBooking(ctx, telegramId)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	tgui.Respond(ctx, b, update, myBookingsListText(bookings), keyboardListMyBooking(bookings))
}

func (h *Handler) HandlerBookingDetail(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "bookings:booking:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	bookingID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	booking, err := h.service.GetByID(ctx, bookingID)
	if err != nil {
		tgui.Respond(ctx, b, update, "⚠️ Запись не найдена\n\nВозможно, она уже была отменена.", mainMenuKeyboard())
		return
	}

	tgui.Respond(ctx, b, update, bookingDetailText(booking), keyboardBookingDetail(booking.ID))
}

func (h *Handler) HandlerCancelBooking(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "bookings:cancel:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	bookingID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	if err := h.service.Cancel(ctx, bookingID); err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	tgui.Respond(ctx, b, update, "✅ Запись отменена\n\nВы можете выбрать другое удобное время.", cancelBookingKeyboard())
}

const actionErrorText = "⚠️ Не удалось выполнить действие\n\nПопробуйте ещё раз или вернитесь в главное меню."
