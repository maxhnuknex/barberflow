package ui

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func Respond(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
	text string,
	keyboard *models.InlineKeyboardMarkup,
) {
	if update.CallbackQuery == nil {
		send(ctx, b, chatID(update), text, keyboard)
		return
	}

	message := update.CallbackQuery.Message.Message
	if message == nil {
		send(ctx, b, chatID(update), text, keyboard)
		return
	}
	_, err := b.EditMessageText(
		ctx,
		&bot.EditMessageTextParams{
			ChatID:      message.Chat.ID,
			MessageID:   message.ID,
			Text:        text,
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		send(ctx, b, message.Chat.ID, text, keyboard)
	}
}

func AnswerCallbackQuery(
	ctx context.Context,
	b *bot.Bot,
	update *models.Update,
) {
	if update.CallbackQuery == nil {
		return
	}

	_, _ = b.AnswerCallbackQuery(
		ctx,
		&bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		},
	)
}

func send(
	ctx context.Context,
	b *bot.Bot,
	chatID int64,
	text string,
	keyboard *models.InlineKeyboardMarkup,
) {
	_, err := b.SendMessage(
		ctx,
		&bot.SendMessageParams{
			ChatID:      chatID,
			Text:        text,
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		return
	}
}

func chatID(update *models.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	}

	if update.CallbackQuery != nil && update.CallbackQuery.Message.Message != nil {
		return update.CallbackQuery.Message.Message.Chat.ID
	}

	if update.CallbackQuery != nil && update.CallbackQuery.Message.InaccessibleMessage != nil {
		return update.CallbackQuery.Message.InaccessibleMessage.Chat.ID
	}

	return 0
}
