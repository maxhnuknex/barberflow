package start

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
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
	if update.CallbackQuery != nil {
		tgui.AnswerCallbackQuery(ctx, b, update)
	}

	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "✂️ Записаться",
					CallbackData: "booking:start",
				},
				{
					Text:         "📋 Мои записи",
					CallbackData: "bookings:list",
				},
			},
		},
	}

	tgui.Respond(
		ctx,
		b,
		update,
		"✂️ BarberFlow\n\nОнлайн-запись в барбершоп.\n\nВыберите действие:",
		keyboard,
	)
}
