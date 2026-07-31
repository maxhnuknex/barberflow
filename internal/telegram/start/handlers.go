package start

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Записаться",
					CallbackData: "booking:start",
				},
				{
					Text:         "Мои записи",
					CallbackData: "bookings:list",
				},
			},
		},
	}

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatID(update),
		Text:        "Добро пожаловать в BarberFlow",
		ReplyMarkup: keyboard,
	})
	if err != nil {
		// логирование добавим отдельно
		return
	}
}

func chatID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	return update.CallbackQuery.Message.Message.Chat.ID
}
