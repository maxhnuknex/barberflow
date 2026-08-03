package mybookings

import "github.com/go-telegram/bot"

func RegisterHandler(b *bot.Bot, h *Handler) {
	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"bookings:list",
		bot.MatchTypeExact,
		h.HandlerListMyBooking,
	)
}
