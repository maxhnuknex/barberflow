package booking

import (
	"fmt"
	"time"

	"github.com/go-telegram/bot/models"
	tgui "github.com/maxhnucknex/barberflow/internal/delivery/telegram/ui"
	"github.com/maxhnucknex/barberflow/internal/domain"
)

func serviceKeyboard(services []domain.Service) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(services)+1)

	for _, service := range services {
		button := models.InlineKeyboardButton{
			Text: fmt.Sprintf(
				"✂️ %s • %d мин • %s",
				service.Name,
				service.DurationMinutes,
				tgui.Price(service.PriceMinorUnits),
			),
			CallbackData: fmt.Sprintf("booking:service:%d", service.ID),
		}

		rows = append(rows, []models.InlineKeyboardButton{button})
	}

	rows = append(rows, userMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func barberKeyboard(
	serviceID int64,
	barbers []domain.Barber,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(barbers)+2)

	for _, barber := range barbers {
		button := models.InlineKeyboardButton{
			Text: "💈 " + barber.Name,
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

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "↩️ Назад к услугам",
			CallbackData: "booking:start",
		},
	})
	rows = append(rows, userMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func keyboardDate(
	serviceID int64,
	barberID int64,
	days []time.Time,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(days)+2)
	today := time.Now().In(moscowLocation())
	if len(days) > 0 {
		today = days[0]
	}

	for _, day := range days {
		button := models.InlineKeyboardButton{
			Text: tgui.DateButton(day, today),
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

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text:         "↩️ Назад к мастерам",
			CallbackData: fmt.Sprintf("booking:service:%d", serviceID),
		},
	})
	rows = append(rows, userMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func freeSlotKeyboard(
	barberID int64,
	serviceID int64,
	freeSlots []domain.TimeInterval,
) *models.InlineKeyboardMarkup {
	rows := make([][]models.InlineKeyboardButton, 0, len(freeSlots)+2)
	row := make([]models.InlineKeyboardButton, 0, 2)

	for _, slot := range freeSlots {
		button := models.InlineKeyboardButton{
			Text: tgui.TimeInterval(slot.StartsAt, slot.EndsAt),
			CallbackData: fmt.Sprintf(
				"booking:time:%d:%d:%d:%d",
				serviceID,
				barberID,
				slot.StartsAt.Unix(),
				slot.EndsAt.Unix(),
			),
		}

		row = append(row, button)
		if len(row) == 2 {
			rows = append(rows, row)
			row = make([]models.InlineKeyboardButton, 0, 2)
		}
	}
	if len(row) > 0 {
		rows = append(rows, row)
	}

	rows = append(rows, []models.InlineKeyboardButton{
		{
			Text: "↩️ Назад к датам",
			CallbackData: fmt.Sprintf(
				"booking:barber:%d:%d",
				serviceID,
				barberID,
			),
		},
	})
	rows = append(rows, userMainMenuButtonRow())

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: rows,
	}
}

func confirmBookingKeyboard(
	serviceID int64,
	barberID int64,
	slot domain.TimeInterval,
) *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text: "✅ Подтвердить запись",
					CallbackData: fmt.Sprintf(
						"booking:confirm:%d:%d:%d:%d",
						serviceID,
						barberID,
						slot.StartsAt.Unix(),
						slot.EndsAt.Unix(),
					),
				},
			},
			{
				{
					Text: "↩️ Изменить время",
					CallbackData: fmt.Sprintf(
						"booking:date:%d:%d:%s",
						serviceID,
						barberID,
						slot.StartsAt.In(moscowLocation()).Format("2006-01-02"),
					),
				},
			},
			userMainMenuButtonRow(),
		},
	}
}

func userMainMenuButtonRow() []models.InlineKeyboardButton {
	return []models.InlineKeyboardButton{
		{
			Text:         "🏠 Главное меню",
			CallbackData: "/start",
		},
	}
}

func mainMenuKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			userMainMenuButtonRow(),
		},
	}
}

func createdBookingKeyboard() *models.InlineKeyboardMarkup {
	return &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{
					Text:         "📋 Мои записи",
					CallbackData: "bookings:list",
				},
			},
			userMainMenuButtonRow(),
		},
	}
}
