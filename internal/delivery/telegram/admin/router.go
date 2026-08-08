package admin

import "github.com/go-telegram/bot"

func RegisterHandler(b *bot.Bot, h *Handler) {
	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"/admin",
		bot.MatchTypeExact,
		h.adminStart,
	)

	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/admin",
		bot.MatchTypeExact,
		h.adminStart,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:bookings",
		bot.MatchTypeExact,
		h.handleBookingsMenu,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:bookings:today",
		bot.MatchTypeExact,
		h.handleBookingsToday,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:bookings:date",
		bot.MatchTypeExact,
		h.handleBookingsDateMenu,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:bookings:date:",
		bot.MatchTypePrefix,
		h.handleBookingsDate,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:services",
		bot.MatchTypeExact,
		h.handleServices,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:service:disable:",
		bot.MatchTypePrefix,
		h.handleServiceDisable,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:service:",
		bot.MatchTypePrefix,
		h.handleServiceDetail,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:barbers",
		bot.MatchTypeExact,
		h.handleBarbers,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:barber:disable:",
		bot.MatchTypePrefix,
		h.handleBarberDisable,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:barber:",
		bot.MatchTypePrefix,
		h.handleBarberDetail,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:booking:cancel:",
		bot.MatchTypePrefix,
		h.handleBookingCancel,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"admin:booking:",
		bot.MatchTypePrefix,
		h.handleBookingDetail,
	)
}
