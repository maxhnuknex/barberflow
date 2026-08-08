package booking

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
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
	GetServiceByID(ctx context.Context, id int64) (domain.Service, error)
	GetBarberByID(ctx context.Context, id int64) (domain.Barber, error)
	ListBarberByService(ctx context.Context, id int64) ([]domain.Barber, error)
	NextSevenDays() []time.Time
	ListAvailableSlots(
		ctx context.Context,
		serviceID int64,
		barberID int64,
		selectedDate time.Time,
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
}

func (h *Handler) HandlerListActive(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	services, err := h.service.ListActive(ctx)

	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	if len(services) == 0 {
		tgui.Respond(ctx, b, update, "📭 Нет доступных услуг\n\nСейчас запись временно недоступна. Попробуйте позже.", mainMenuKeyboard())
		return
	}

	keyboard := serviceKeyboard(services)

	tgui.Respond(
		ctx,
		b,
		update,
		"✂️ Выберите услугу\n\nПосле выбора покажем подходящих мастеров и свободное время.",
		keyboard,
	)
}

func (h *Handler) HandlerBookingTime(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "booking:time:"

	data := update.CallbackQuery.Data
	rawValues := strings.TrimPrefix(data, prefix)

	parts := strings.Split(rawValues, ":")
	if len(parts) != 4 {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	serviceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	barberID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	startUnix, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	endUnix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	startsAt := time.Unix(startUnix, 0)
	endsAt := time.Unix(endUnix, 0)

	service, err := h.service.GetServiceByID(ctx, serviceID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	barber, err := h.service.GetBarberByID(ctx, barberID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	slot := domain.TimeInterval{
		StartsAt: startsAt,
		EndsAt:   endsAt,
	}

	text := bookingPreviewText(
		service,
		barber,
		slot,
	)

	tgui.Respond(
		ctx,
		b,
		update,
		text,
		confirmBookingKeyboard(serviceID, barberID, slot),
	)
}

func (h *Handler) HandlerConfirmBooking(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "booking:confirm:"

	data := update.CallbackQuery.Data
	rawValues := strings.TrimPrefix(data, prefix)

	parts := strings.Split(rawValues, ":")
	if len(parts) != 4 {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	serviceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	barberID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	startUnix, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	endUnix, err := strconv.ParseInt(parts[3], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
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
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	text := createdBookingText(booking)

	tgui.Respond(
		ctx,
		b,
		update,
		text,
		createdBookingKeyboard(),
	)
}

func (h *Handler) HandleSelectDate(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "booking:date:"

	data := update.CallbackQuery.Data
	rawValues := strings.TrimPrefix(data, prefix)

	parts := strings.Split(rawValues, ":")
	if len(parts) != 3 {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	serviceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	barberID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	selectedDate, err := time.ParseInLocation("2006-01-02", parts[2], moscowLocation())
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	freeSlot, err := h.service.ListAvailableSlots(ctx, serviceID, barberID, selectedDate)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	if len(freeSlot) == 0 {
		tgui.Respond(
			ctx,
			b,
			update,
			"📭 Нет свободного времени\n\nНа выбранную дату все интервалы заняты. Выберите другой день.",
			keyboardDate(serviceID, barberID, h.service.NextSevenDays()),
		)
		return
	}

	keyboard := freeSlotKeyboard(
		barberID,
		serviceID,
		freeSlot,
	)

	tgui.Respond(
		ctx,
		b,
		update,
		fmt.Sprintf("🕒 Выберите время\n\nСвободные интервалы на %s:", tgui.DayMonth(selectedDate)),
		keyboard,
	)
}

func (h *Handler) HandlerListBarber(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	data := update.CallbackQuery.Data
	const prefix = "booking:service:"
	idStr := strings.TrimPrefix(data, prefix)
	serviceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	barbers, err := h.service.ListBarberByService(ctx, serviceID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	if len(barbers) == 0 {
		tgui.Respond(
			ctx,
			b,
			update,
			"📭 Нет доступных мастеров\n\nДля выбранной услуги сейчас нет доступных специалистов. Выберите другую услугу.",
			barberKeyboard(serviceID, nil),
		)
		return
	}

	keyboard := barberKeyboard(serviceID, barbers)

	tgui.Respond(ctx, b, update, "💈 Выберите мастера\n\nВыберите специалиста, к которому хотите записаться.", keyboard)
}

func (h *Handler) HandlerListTimeBarberActive(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "booking:barber:"

	data := update.CallbackQuery.Data
	rawIDs := strings.TrimPrefix(data, prefix)

	parts := strings.Split(rawIDs, ":")
	if len(parts) != 2 {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	serviceID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	barberID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, mainMenuKeyboard())
		return
	}

	keyboard := keyboardDate(serviceID, barberID, h.service.NextSevenDays())

	tgui.Respond(
		ctx,
		b,
		update,
		"📅 Выберите дату\n\nДоступна запись на ближайшие 7 дней.",
		keyboard,
	)
}

func moscowLocation() *time.Location {
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return time.FixedZone("Europe/Moscow", 3*60*60)
	}

	return loc
}

const actionErrorText = "⚠️ Не удалось выполнить действие\n\nПопробуйте ещё раз или вернитесь в главное меню."

func bookingPreviewText(
	service domain.Service,
	barber domain.Barber,
	slot domain.TimeInterval,
) string {
	return fmt.Sprintf(
		"✅ Проверьте запись\n\n✂️ Услуга: %s\n💈 Мастер: %s\n📅 Дата: %s\n🕒 Время: %s\n💳 Стоимость: %s\n\nВсё верно?",
		service.Name,
		barber.Name,
		tgui.FullDate(slot.StartsAt),
		tgui.TimeInterval(slot.StartsAt, slot.EndsAt),
		tgui.Price(service.PriceMinorUnits),
	)
}

func createdBookingText(booking domain.Booking) string {
	return fmt.Sprintf(
		"🎉 Запись подтверждена\n\n✂️ Услуга: %s\n💈 Мастер: %s\n📅 Дата: %s\n🕒 Время: %s\n💳 Стоимость: %s\n\nЗапись сохранена в разделе «Мои записи».",
		booking.ServiceName,
		booking.BarberName,
		tgui.FullDate(booking.StartsAt),
		tgui.TimeInterval(booking.StartsAt, booking.EndsAt),
		tgui.Price(booking.PriceMinorUnits),
	)
}
