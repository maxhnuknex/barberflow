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
}
