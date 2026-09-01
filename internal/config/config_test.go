package config

import "testing"

func TestLoadFullDefaultsAndValidation(t *testing.T) {
	t.Setenv("TELEGRAM_BOT_TOKEN", "token")
	t.Setenv("TELEGRAM_ADMIN_USER_ID", "123")
	t.Setenv("TELEGRAM_TARGET_CHAT_ID", "-456")
	cfg, err := LoadFull()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.PollMin.String() != "10m0s" || cfg.PollMax.String() != "30m0s" || cfg.MinCards != 20 || cfg.MaxDetailsPerPoll != 20 {
		t.Fatalf("defaults=%+v", cfg)
	}
	t.Setenv("SOMON_POLL_MIN", "31m")
	if _, err := LoadFull(); err == nil {
		t.Fatal("expected invalid poll range")
	}
}
