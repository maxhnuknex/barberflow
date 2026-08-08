package mybookings

import "github.com/go-telegram/bot"

func RegisterHandler(b *bot.Bot, h *Handler) {
	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"bookings:list",
		bot.MatchTypeExact,
		h.HandlerListMyBooking,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"bookings:cancel:",
		bot.MatchTypePrefix,
		h.HandlerCancelBooking,
	)

	b.RegisterHandler(
		bot.HandlerTypeCallbackQueryData,
		"bookings:booking:",
		bot.MatchTypePrefix,
		h.HandlerBookingDetail,
	)
}
