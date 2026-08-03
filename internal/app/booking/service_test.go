package booking

import "testing"

func TestNextSevenDays(t *testing.T) {
	service := NewService(nil, nil, nil)

	days := service.NextSevenDays()
	if len(days) != 7 {
		t.Fatalf("expected 7 days, got %d", len(days))
	}

	for i, day := range days {
		if day.Location().String() != "Europe/Moscow" {
			t.Fatalf("expected Europe/Moscow location, got %s", day.Location())
		}
		if day.Hour() != 0 || day.Minute() != 0 || day.Second() != 0 || day.Nanosecond() != 0 {
			t.Fatalf("expected date without time, got %s", day)
		}
		if i == 0 {
			continue
		}

		expected := days[i-1].AddDate(0, 0, 1)
		if !day.Equal(expected) {
			t.Fatalf("expected %s, got %s", expected, day)
		}
	}
}
