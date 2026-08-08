package admin

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
	bookingService BookingService
	catalogService CatalogService
	barberService  BarberService
}

func NewHandler(
	bookingService BookingService,
	catalogService CatalogService,
	barberService BarberService,
) *Handler {
	return &Handler{
		bookingService: bookingService,
		catalogService: catalogService,
		barberService:  barberService,
	}
}

type BookingService interface {
	Cancel(ctx context.Context, bookingID int64) error
	GetByID(ctx context.Context, bookingID int64) (domain.Booking, error)
	ListByDate(ctx context.Context, date time.Time) ([]domain.Booking, error)
	NextSevenDays() []time.Time
}

type CatalogService interface {
	GetByID(ctx context.Context, id int64) (domain.Service, error)
	ListAll(ctx context.Context) ([]domain.Service, error)
	SetActive(ctx context.Context, id int64, active bool) error
}

type BarberService interface {
	GetByID(ctx context.Context, id int64) (domain.Barber, error)
	ListAll(ctx context.Context) ([]domain.Barber, error)
	SetActive(ctx context.Context, id int64, active bool) error
}

func (h *Handler) adminStart(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery != nil {
		tgui.AnswerCallbackQuery(ctx, b, update)
	}

	keyboard := keyboardAdminStart()

	tgui.Respond(
		ctx,
		b,
		update,
		"🛠 Админ-панель\n\nУправление записями, услугами и мастерами.",
		keyboard,
	)
}

func (h *Handler) handleBookingsMenu(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	tgui.Respond(
		ctx,
		b,
		update,
		"📋 Управление записями\n\nВыберите период просмотра.",
		keyboardBookingsMenu(),
	)
}

func (h *Handler) handleBookingsToday(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	bookings, err := h.bookingService.ListByDate(ctx, time.Now().In(moscowLocation()))
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBookingsMenu())
		return
	}

	selectedDate := time.Now().In(moscowLocation())
	tgui.Respond(
		ctx,
		b,
		update,
		bookingsListText(bookings, selectedDate, true),
		keyboardBookingsList(bookings),
	)
}

func (h *Handler) handleBookingsDateMenu(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	tgui.Respond(
		ctx,
		b,
		update,
		"📅 Выберите дату\n\nПоказаны ближайшие 7 дней.",
		keyboardBookingsDate(h.bookingService.NextSevenDays()),
	)
}

func (h *Handler) handleBookingsDate(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:bookings:date:"

	rawDate := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	selectedDate, err := time.ParseInLocation("2006-01-02", rawDate, moscowLocation())
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBookingsMenu())
		return
	}

	bookings, err := h.bookingService.ListByDate(ctx, selectedDate)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBookingsMenu())
		return
	}

	tgui.Respond(
		ctx,
		b,
		update,
		bookingsListText(bookings, selectedDate, false),
		keyboardBookingsList(bookings),
	)
}

func (h *Handler) handleBookingDetail(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:booking:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	bookingID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBookingsMenu())
		return
	}

	booking, err := h.bookingService.GetByID(ctx, bookingID)
	if err != nil {
		tgui.Respond(ctx, b, update, "⚠️ Запись не найдена\n\nВозможно, она уже была отменена.", keyboardBookingsMenu())
		return
	}

	tgui.Respond(ctx, b, update, bookingDetailText(booking), keyboardBookingDetail(booking.ID))
}

func (h *Handler) handleBookingCancel(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:booking:cancel:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	bookingID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBookingsMenu())
		return
	}

	if err := h.bookingService.Cancel(ctx, bookingID); err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBookingsMenu())
		return
	}

	tgui.Respond(ctx, b, update, "✅ Запись отменена\n\nЗапись больше не отображается среди активных.", keyboardBookingCanceled())
}

func (h *Handler) handleServices(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	services, err := h.catalogService.ListAll(ctx)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardAdminStart())
		return
	}

	if len(services) == 0 {
		tgui.Respond(ctx, b, update, "📭 Услуг нет\n\nВ каталоге пока нет услуг.", keyboardAdminStart())
		return
	}

	tgui.Respond(ctx, b, update, "✂️ Услуги\n\nВыберите услугу для просмотра.", keyboardServicesList(services))
}

func (h *Handler) handleServiceDetail(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:service:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	serviceID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToServices())
		return
	}

	service, err := h.catalogService.GetByID(ctx, serviceID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToServices())
		return
	}

	tgui.Respond(ctx, b, update, serviceDetailText(service), keyboardServiceDetail(service))
}

func (h *Handler) handleServiceDisable(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:service:disable:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	serviceID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToServices())
		return
	}

	if err := h.catalogService.SetActive(ctx, serviceID, false); err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToServices())
		return
	}

	service, err := h.catalogService.GetByID(ctx, serviceID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToServices())
		return
	}

	tgui.Respond(ctx, b, update, serviceDetailText(service), keyboardServiceDetail(service))
}

func (h *Handler) handleBarbers(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	barbers, err := h.barberService.ListAll(ctx)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardAdminStart())
		return
	}

	if len(barbers) == 0 {
		tgui.Respond(ctx, b, update, "📭 Мастеров нет\n\nВ системе пока нет мастеров.", keyboardAdminStart())
		return
	}

	tgui.Respond(ctx, b, update, "💈 Мастера\n\nВыберите мастера для просмотра.", keyboardBarbersList(barbers))
}

func (h *Handler) handleBarberDetail(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:barber:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	barberID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToBarbers())
		return
	}

	barber, err := h.barberService.GetByID(ctx, barberID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToBarbers())
		return
	}

	tgui.Respond(ctx, b, update, barberDetailText(barber), keyboardBarberDetail(barber))
}

func (h *Handler) handleBarberDisable(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}
	tgui.AnswerCallbackQuery(ctx, b, update)

	const prefix = "admin:barber:disable:"

	rawID := strings.TrimPrefix(update.CallbackQuery.Data, prefix)
	barberID, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToBarbers())
		return
	}

	if err := h.barberService.SetActive(ctx, barberID, false); err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToBarbers())
		return
	}

	barber, err := h.barberService.GetByID(ctx, barberID)
	if err != nil {
		tgui.Respond(ctx, b, update, actionErrorText, keyboardBackToBarbers())
		return
	}

	tgui.Respond(ctx, b, update, barberDetailText(barber), keyboardBarberDetail(barber))
}

func bookingsListText(bookings []domain.Booking, selectedDate time.Time, today bool) string {
	if len(bookings) == 0 {
		return "📭 Записей нет\n\nНа выбранную дату активных записей нет."
	}

	title := fmt.Sprintf("📋 Записи на %s", tgui.DayMonth(selectedDate))
	if today {
		title = "📋 Записи на сегодня"
	}

	return fmt.Sprintf("%s\n\nНайдено записей: %d.\n\nВыберите запись для просмотра.", title, len(bookings))
}

const actionErrorText = "⚠️ Не удалось выполнить действие\n\nПопробуйте ещё раз или вернитесь в админ-панель."

func bookingDetailText(booking domain.Booking) string {
	return fmt.Sprintf(
		"📋 Запись #%d\n\n🙍 Клиент: %s\n✂️ Услуга: %s\n💈 Мастер: %s\n📅 Дата: %s\n🕒 Время: %s\n💳 Стоимость: %s",
		booking.ID,
		customerText(booking),
		booking.ServiceName,
		booking.BarberName,
		tgui.FullDate(booking.StartsAt),
		tgui.TimeInterval(booking.StartsAt, booking.EndsAt),
		tgui.Price(booking.PriceMinorUnits),
	)
}

func serviceDetailText(service domain.Service) string {
	return fmt.Sprintf(
		"✂️ Услуга #%d\n\nНазвание: %s\n⏱ Длительность: %d мин\n💳 Стоимость: %s\nСтатус: %s",
		service.ID,
		service.Name,
		service.DurationMinutes,
		tgui.Price(service.PriceMinorUnits),
		serviceStatusText(service.Active),
	)
}

func barberDetailText(barber domain.Barber) string {
	return fmt.Sprintf(
		"💈 Мастер #%d\n\nИмя: %s\nСтатус: %s",
		barber.ID,
		barber.Name,
		barberStatusText(barber.Active),
	)
}

func serviceStatusText(active bool) string {
	if active {
		return "✅ Активна"
	}

	return "⛔ Отключена"
}

func barberStatusText(active bool) string {
	if active {
		return "✅ Работает"
	}

	return "⛔ Отключён"
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
