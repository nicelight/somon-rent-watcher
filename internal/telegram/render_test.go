package telegram

import (
	"fmt"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/nicelight/somon-rent-watcher/internal/filter"
	"github.com/nicelight/somon-rent-watcher/internal/model"
)

func intp(v int) *int { return &v }

func TestSettingsTextAndKeyboards(t *testing.T) {
	s := filter.DefaultSettings()
	s.PriceMin, s.PriceMax = intp(3500), intp(6000)
	s.SellerAdsLimit = 5
	s.NegativeWords = []string{"посуточно"}
	text := SettingsText(s)
	for _, want := range []string{"на паузе", "3 500 — 6 000", "меньше 5", "посуточно", "Обычные + VIP", "10–30 мин"} {
		if !strings.Contains(text, want) {
			t.Fatalf("settings text lacks %q: %s", want, text)
		}
	}
	if rows := floorsKeyboard(s).InlineKeyboard; len(rows) != 5 {
		t.Fatalf("floor keyboard rows=%d", len(rows))
	}
	if got := mainKeyboard(s).InlineKeyboard[0][0].Text; !strings.Contains(got, "Включить") {
		t.Fatalf("paused toggle=%q", got)
	}
	if got := mainKeyboard(s).InlineKeyboard[1][1].Text; got != "Сканировать сейчас" {
		t.Fatalf("scan button=%q", got)
	}
	busyButton := mainKeyboardWithScanState(s, true).InlineKeyboard[1][1]
	if busyButton.Text != "⏳ Сканирую..." || busyButton.CallbackData != "e:scan_busy" {
		t.Fatalf("busy scan button=%+v", busyButton)
	}
	s.Enabled = true
	if got := mainKeyboard(s).InlineKeyboard[0][0].Text; !strings.Contains(got, "паузу") {
		t.Fatalf("enabled toggle=%q", got)
	}
}

func TestAdCaptionEscapesAndFits(t *testing.T) {
	ad := model.Ad{
		Card: model.Card{
			ID:       123,
			URL:      "https://somon.tj/adv/123_x/",
			Title:    `2-комн. <квартира> & тест`,
			Price:    intp(5999),
			Floor:    intp(14),
			Promoted: true,
			AgeText:  "3 минуты назад",
		},
		Description: strings.Repeat("описание ", 100),
		SellerName:  `Rofi & Co`,
		SellerAds:   intp(17),
	}
	caption := AdCaption(ad)
	if strings.Contains(caption, "<квартира>") || !strings.Contains(caption, "&lt;квартира&gt;") || !strings.Contains(caption, "5 999 c.") {
		t.Fatalf("caption escaping/price: %s", caption)
	}
	if utf8.RuneCountInString(caption) > 900 {
		t.Fatalf("caption too long: %d", utf8.RuneCountInString(caption))
	}
}

func TestAdCaptionStartsEveryEmojiFieldOnNewLine(t *testing.T) {
	ad := model.Ad{
		Card:        model.Card{ID: 17041834},
		Description: "Сдаётся квартира: Садбарг Количество комнат: 2 🔺 Площадь: 50 🏠 Тип: новостройка 🏢 Этаж: 11 🖼️ Ремонт: новый ремонт 💵 Цена: 5000",
	}

	caption := AdCaption(ad)
	for _, want := range []string{
		"Количество комнат: 2\n🔺 Площадь: 50",
		"\n🏠 Тип: новостройка",
		"\n🏢 Этаж: 11",
		"\n🖼️ Ремонт: новый ремонт",
		"\n💵 Цена: 5000",
	} {
		if !strings.Contains(caption, want) {
			t.Fatalf("caption lacks emoji line %q:\n%s", want, caption)
		}
	}
}

func TestTelegramErrorsRedactTokenURL(t *testing.T) {
	client := NewClient("https://api.telegram.org", "SECRET_TOKEN")
	got := client.redactError(fmt.Errorf(`Post "https://api.telegram.org/botSECRET_TOKEN/getUpdates": timeout`))
	if strings.Contains(got, "SECRET_TOKEN") || !strings.Contains(got, "<telegram-bot-api>") {
		t.Fatalf("redaction failed: %s", got)
	}
}

func TestTelegramRetryDelayHonorsRetryAfter(t *testing.T) {
	if got := telegramRetryDelay(&APIError{Code: 429, RetryAfter: 17}); got != 17*time.Second {
		t.Fatalf("retry delay=%s", got)
	}
	if got := telegramRetryDelay(fmt.Errorf("network")); got != 5*time.Second {
		t.Fatalf("default retry delay=%s", got)
	}
}
