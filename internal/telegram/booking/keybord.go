package booking

import (
	"github.com/go-telegram/bot/models"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func serviceKeyboard(services []domain.Service) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(services))

	for _, service := range services {
		button := models.InlineKeyboardButton{
			Text:         service.Name,
			CallbackData: "booking:services",
		}

		rows = append(rows, []models.InlineKeyboardButton{button})
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
