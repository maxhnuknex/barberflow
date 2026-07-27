package start

import "github.com/go-telegram/bot"

func RegisterRoutes(b *bot.Bot, h *Handler) {
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"/start",
		bot.MatchTypeExact,
		h.Handle,
	)
}
