package ui

import (
	"testing"
	"time"
)

func TestPrice(t *testing.T) {
	got := Price(250000)
	want := "2 500 ₽"

	if got != want {
		t.Fatalf("unexpected price:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestFullDate(t *testing.T) {
	got := FullDate(time.Date(2026, 8, 4, 15, 0, 0, 0, time.UTC))
	want := "4 августа 2026"

	if got != want {
		t.Fatalf("unexpected date:\nwant: %q\ngot:  %q", want, got)
	}
}

func TestDateButton(t *testing.T) {
	today := time.Date(2026, 8, 4, 0, 0, 0, 0, time.UTC)

	cases := []struct {
		name string
		date time.Time
		want string
	}{
		{
			name: "today",
			date: today,
			want: "Сегодня • 04.08",
		},
		{
			name: "tomorrow",
			date: today.AddDate(0, 0, 1),
			want: "Завтра • 05.08",
		},
		{
			name: "weekday",
			date: today.AddDate(0, 0, 2),
			want: "Чт • 06.08",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got := DateButton(tt.date, today)
			if got != tt.want {
				t.Fatalf("unexpected date button:\nwant: %q\ngot:  %q", tt.want, got)
			}
		})
	}
}
