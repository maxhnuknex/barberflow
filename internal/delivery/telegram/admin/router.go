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
		"admin:bookings:find",
		bot.MatchTypeExact,
		h.handleBookingsFind,
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
