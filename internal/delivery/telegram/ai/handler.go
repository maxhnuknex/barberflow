package ai

import (
	"context"
	"log/slog"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const unavailableText = "Извините, сейчас не могу подготовить ответ. Попробуйте ещё раз позже."
const maxTelegramMessageRunes = 4096

// Client is the part of an LLM client the Telegram delivery layer needs.
type Client interface {
	Complete(ctx context.Context, text string) (string, error)
}

type Handler struct {
	client Client
	logger *slog.Logger
}

func NewHandler(client Client, logger *slog.Logger) *Handler {
	return &Handler{
		client: client,
		logger: logger,
	}
}

// Handle sends a regular text message to the LLM and returns its reply to the same chat.
func (h *Handler) Handle(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil || strings.TrimSpace(update.Message.Text) == "" {
		return
	}

	h.logger.Info("processing AI message", "chat_id", update.Message.Chat.ID)

	response, err := h.client.Complete(ctx, update.Message.Text)
	if err != nil {
		h.logger.Error("failed to complete AI request", "error", err)
		h.sendResponse(ctx, b, update.Message.Chat.ID, unavailableText)
		return
	}

	response = strings.TrimSpace(response)
	if response == "" {
		response = unavailableText
	}

	h.logger.Info("AI completion received", "chat_id", update.Message.Chat.ID, "response_length", len([]rune(response)))
	h.sendResponse(ctx, b, update.Message.Chat.ID, response)
}

func (h *Handler) sendResponse(ctx context.Context, b *bot.Bot, chatID int64, text string) {
	for _, part := range splitMessage(text) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   part,
		}); err != nil {
			h.logger.Error("failed to send AI response", "chat_id", chatID, "error", err)
			return
		}
	}
}

func splitMessage(text string) []string {
	runes := []rune(text)
	if len(runes) <= maxTelegramMessageRunes {
		return []string{text}
	}

	parts := make([]string, 0, (len(runes)+maxTelegramMessageRunes-1)/maxTelegramMessageRunes)
	for len(runes) > 0 {
		end := min(len(runes), maxTelegramMessageRunes)
		parts = append(parts, string(runes[:end]))
		runes = runes[end:]
	}

	return parts
}
