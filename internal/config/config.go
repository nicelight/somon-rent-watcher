package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/somon"
)

type Config struct {
	TelegramBotToken     string
	TelegramAdminUserID  int64
	TelegramTargetChatID int64
	TelegramAPIBase      string

	DBPath      string
	DebugDir    string
	CategoryURL string
	UserAgent   string

	PollMin           time.Duration
	PollMax           time.Duration
	GapAfter          time.Duration
	BlockBackoff      time.Duration
	RequestDelay      time.Duration
	HTTPTimeout       time.Duration
	MinCards          int
	MaxDetailsPerPoll int
	MaxBodyBytes      int64
	LogLevel          string
}

func LoadFull() (Config, error) {
	cfg := defaults()
	cfg.TelegramBotToken = strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if cfg.TelegramBotToken == "" {
		return Config{}, errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	var err error
	cfg.TelegramAdminUserID, err = requiredInt64("TELEGRAM_ADMIN_USER_ID")
	if err != nil {
		return Config{}, err
	}
	cfg.TelegramTargetChatID, err = requiredInt64("TELEGRAM_TARGET_CHAT_ID")
	if err != nil {
		return Config{}, err
	}

	cfg.DBPath = envOr("DB_PATH", cfg.DBPath)
	cfg.DebugDir = envOr("DEBUG_DIR", filepath.Join(filepath.Dir(cfg.DBPath), "debug"))
	cfg.CategoryURL = envOr("SOMON_CATEGORY_URL", cfg.CategoryURL)
	cfg.UserAgent = envOr("SOMON_USER_AGENT", cfg.UserAgent)
	cfg.TelegramAPIBase = strings.TrimRight(envOr("TELEGRAM_API_BASE", cfg.TelegramAPIBase), "/")
	cfg.LogLevel = strings.ToLower(envOr("LOG_LEVEL", cfg.LogLevel))

	if cfg.PollMin, err = durationEnv("SOMON_POLL_MIN", cfg.PollMin); err != nil {
		return Config{}, err
	}
	if cfg.PollMax, err = durationEnv("SOMON_POLL_MAX", cfg.PollMax); err != nil {
		return Config{}, err
	}
	if cfg.GapAfter, err = durationEnv("SOMON_GAP_AFTER", cfg.GapAfter); err != nil {
		return Config{}, err
	}
	if cfg.BlockBackoff, err = durationEnv("SOMON_BLOCK_BACKOFF", cfg.BlockBackoff); err != nil {
		return Config{}, err
	}
	if cfg.RequestDelay, err = durationEnv("SOMON_REQUEST_DELAY", cfg.RequestDelay); err != nil {
		return Config{}, err
	}
	if cfg.HTTPTimeout, err = durationEnv("HTTP_TIMEOUT", cfg.HTTPTimeout); err != nil {
		return Config{}, err
	}
	if cfg.MinCards, err = intEnv("SOMON_MIN_CARDS", cfg.MinCards); err != nil {
		return Config{}, err
	}
	if cfg.MaxDetailsPerPoll, err = intEnv("SOMON_MAX_DETAILS_PER_POLL", cfg.MaxDetailsPerPoll); err != nil {
		return Config{}, err
	}
	if cfg.MaxBodyBytes, err = int64Env("SOMON_MAX_BODY_BYTES", cfg.MaxBodyBytes); err != nil {
		return Config{}, err
	}

	if cfg.TelegramAdminUserID <= 0 {
		return Config{}, errors.New("TELEGRAM_ADMIN_USER_ID must be positive")
	}
	if cfg.TelegramTargetChatID == 0 {
		return Config{}, errors.New("TELEGRAM_TARGET_CHAT_ID must be non-zero")
	}
	if cfg.PollMin <= 0 || cfg.PollMax <= 0 || cfg.PollMin > cfg.PollMax {
		return Config{}, errors.New("invalid SOMON_POLL_MIN/SOMON_POLL_MAX")
	}
	if cfg.GapAfter <= 0 || cfg.BlockBackoff <= 0 || cfg.HTTPTimeout <= 0 {
		return Config{}, errors.New("durations must be positive")
	}
	if cfg.RequestDelay < 0 {
		return Config{}, errors.New("SOMON_REQUEST_DELAY must be non-negative")
	}
	if cfg.MinCards < 1 {
		return Config{}, errors.New("SOMON_MIN_CARDS must be >= 1")
	}
	if cfg.MaxDetailsPerPoll < 1 || cfg.MaxDetailsPerPoll > 100 {
		return Config{}, errors.New("SOMON_MAX_DETAILS_PER_POLL must be between 1 and 100")
	}
	if cfg.MaxBodyBytes < 1<<20 {
		return Config{}, errors.New("SOMON_MAX_BODY_BYTES must be at least 1 MiB")
	}
	return cfg, nil
}

func LoadTokenOnly() (token, apiBase string, err error) {
	token = strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		return "", "", errors.New("TELEGRAM_BOT_TOKEN is required")
	}
	return token, strings.TrimRight(envOr("TELEGRAM_API_BASE", "https://api.telegram.org"), "/"), nil
}

func defaults() Config {
	return Config{
		TelegramAPIBase:   "https://api.telegram.org",
		DBPath:            "./somonwatch.db",
		CategoryURL:       somon.DefaultCategoryURL,
		UserAgent:         "SomonRentWatcher/1.0 (private Telegram notifier; low-rate polling)",
		PollMin:           10 * time.Minute,
		PollMax:           30 * time.Minute,
		GapAfter:          45 * time.Minute,
		BlockBackoff:      2 * time.Hour,
		RequestDelay:      900 * time.Millisecond,
		HTTPTimeout:       25 * time.Second,
		MinCards:          20,
		MaxDetailsPerPoll: 20,
		MaxBodyBytes:      10 << 20,
		LogLevel:          "info",
	}
}

func requiredInt64(name string) (int64, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return n, nil
}

func envOr(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func durationEnv(name string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return d, nil
}

func intEnv(name string, fallback int) (int, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return n, nil
}

func int64Env(name string, fallback int64) (int64, error) {
	value := strings.TrimSpace(os.Getenv(name))
	if value == "" {
		return fallback, nil
	}
	n, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", name, err)
	}
	return n, nil
}
