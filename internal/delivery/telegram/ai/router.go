package ai

import "github.com/go-telegram/bot"

func RegisterHandler(b *bot.Bot, h *Handler) {
	// An empty prefix matches every text message. Register this handler after
	// command handlers so /start and /admin keep their dedicated routes.
	b.RegisterHandler(
		bot.HandlerTypeMessageText,
		"",
		bot.MatchTypePrefix,
		h.Handle,
	)
}
