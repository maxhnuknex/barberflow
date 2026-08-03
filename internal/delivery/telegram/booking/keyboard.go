package booking

import (
	"fmt"
	"time"

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

func barberKeyboard(
	serviceID int64,
	barbers []domain.Barber,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(barbers))

	for _, barber := range barbers {
		button := models.InlineKeyboardButton{
			Text: barber.Name,
			CallbackData: fmt.Sprintf(
				"booking:barber:%d:%d",
				serviceID,
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

func keyboardDate(
	serviceID int64,
	barberID int64,
	days []time.Time,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(days))

	for _, day := range days {
		button := models.InlineKeyboardButton{
			Text: day.Format("02.01.2006"),
			CallbackData: fmt.Sprintf(
				"booking:date:%d:%d:%s",
				serviceID,
				barberID,
				day.Format("2006-01-02"),
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

func freeSlotKeyboard(
	barberID int64,
	serviceID int64,
	freeSlots []domain.TimeInterval,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(freeSlots))

	for _, slot := range freeSlots {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"%s-%s",
				slot.StartsAt.Format("15:04"),
				slot.EndsAt.Format("15:04"),
			),
			CallbackData: fmt.Sprintf(
				"booking:time:%d:%d:%d:%d",
				serviceID,
				barberID,
				slot.StartsAt.Unix(),
				slot.EndsAt.Unix(),
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

func mainMenuKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "Главное меню",
					CallbackData: "/start",
				},
			},
		},
	}
}
