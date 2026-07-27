package booking

import (
	"fmt"

	"github.com/go-telegram/bot/models"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func serviceKeyboard(services []domain.Service) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(services))

	for _, service := range services {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"%s — %d мин — %d ₽",
				service.Name,
				service.DurationMinutes,
				service.PriceMinorUnits/100,
			),
			CallbackData: fmt.Sprintf("booking:service:%d", service.ID),
		}

		rows = append(rows, []models.InlineKeyboardButton{button})
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func barberKeyboard(barbers []domain.Barber) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(barbers))

	for _, barber := range barbers {
		button := models.InlineKeyboardButton{
			Text: barber.Name,
			CallbackData: fmt.Sprintf(
				"booking:barber:%d",
				barber.ID,
			),
		}

		rows = append(rows, []models.InlineKeyboardButton{
			button,
		})
	}

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}
