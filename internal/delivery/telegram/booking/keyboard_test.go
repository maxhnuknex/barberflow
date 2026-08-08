package booking

import (
	"testing"
	"time"

	"github.com/maxhnucknex/barberflow/internal/domain"
)

func TestConfirmBookingKeyboard(t *testing.T) {
	slot := domain.TimeInterval{
		StartsAt: time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC),
		EndsAt:   time.Date(2026, 8, 4, 16, 0, 0, 0, time.UTC),
	}

	keyboard := confirmBookingKeyboard(2, 5, slot)

	confirmButton := keyboard.InlineKeyboard[0][0]
	if confirmButton.Text != "✅ Подтвердить запись" {
		t.Fatalf("unexpected confirm text: %q", confirmButton.Text)
	}
	if confirmButton.CallbackData != "booking:confirm:2:5:1785855600:1785859200" {
		t.Fatalf("unexpected confirm callback: %q", confirmButton.CallbackData)
	}

	backButton := keyboard.InlineKeyboard[1][0]
	if backButton.Text != "↩️ Изменить время" {
		t.Fatalf("unexpected back text: %q", backButton.Text)
	}
	if backButton.CallbackData != "booking:date:2:5:2026-08-04" {
		t.Fatalf("unexpected back callback: %q", backButton.CallbackData)
	}

	mainMenuButton := keyboard.InlineKeyboard[2][0]
	if mainMenuButton.Text != "🏠 Главное меню" {
		t.Fatalf("unexpected main menu text: %q", mainMenuButton.Text)
	}
	if mainMenuButton.CallbackData != "/start" {
		t.Fatalf("unexpected main menu callback: %q", mainMenuButton.CallbackData)
	}
}

func TestBookingNavigationButtons(t *testing.T) {
	serviceKeyboard := serviceKeyboard(nil)
	if serviceKeyboard.InlineKeyboard[0][0].CallbackData != "/start" {
		t.Fatalf("unexpected service main menu callback: %q", serviceKeyboard.InlineKeyboard[0][0].CallbackData)
	}

	barberKeyboard := barberKeyboard(2, nil)
	if barberKeyboard.InlineKeyboard[0][0].CallbackData != "booking:start" {
		t.Fatalf("unexpected barber back callback: %q", barberKeyboard.InlineKeyboard[0][0].CallbackData)
	}
	if barberKeyboard.InlineKeyboard[1][0].CallbackData != "/start" {
		t.Fatalf("unexpected barber main menu callback: %q", barberKeyboard.InlineKeyboard[1][0].CallbackData)
	}
}
