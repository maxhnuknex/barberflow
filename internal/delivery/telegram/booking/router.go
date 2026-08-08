package booking

import (
	"github.com/go-telegram/bot"
)

func RegisterHandler(b *bot.Bot, h *Handler) {
	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"booking:start",
		bot.MatchTypeExact,
		h.HandlerListActive,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"booking:service",
		bot.MatchTypePrefix,
		h.HandlerListBarber,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"booking:barber",
		bot.MatchTypePrefix,
		h.HandlerListTimeBarberActive,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"booking:date",
		bot.MatchTypePrefix,
		h.HandleSelectDate,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"booking:time",
		bot.MatchTypePrefix,
		h.HandlerBookingTime,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"booking:confirm",
		bot.MatchTypePrefix,
		h.HandlerConfirmBooking,
	)

}
