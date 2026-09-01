package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/app"
	"github.com/nicelight/somon-rent-watcher/internal/config"
	"github.com/nicelight/somon-rent-watcher/internal/somon"
	"github.com/nicelight/somon-rent-watcher/internal/store"
	"github.com/nicelight/somon-rent-watcher/internal/telegram"
)

var (
	version   = "dev"
	commit    = "unknown"
	buildTime = "unknown"
)

func main() {
	command := "run"
	if len(os.Args) > 1 {
		command = strings.ToLower(strings.TrimSpace(os.Args[1]))
	}

	var err error
	switch command {
	case "run":
		err = runService()
	case "doctor":
		err = runDoctor()
	case "ids":
		err = runIDs()
	case "version", "--version", "-version", "-v":
		fmt.Printf("somonwatch %s (commit %s, built %s)\n", version, commit, buildTime)
		return
	case "help", "--help", "-h":
		printHelp()
		return
	default:
		printHelp()
		err = fmt.Errorf("unknown command %q", command)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR:", err)
		os.Exit(1)
	}
}

func runService() error {
	cfg, err := config.LoadFull()
	if err != nil {
		return err
	}
	logger := newLogger(cfg.LogLevel)
	slog.SetDefault(logger)

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		return err
	}
	defer db.Close()

	application, err := app.New(cfg, db, logger)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	logger.Info("somonwatch starting", "version", version, "db", cfg.DBPath, "category", cfg.CategoryURL)
	if err := application.Run(ctx); err != nil {
		return err
	}
	logger.Info("somonwatch stopped")
	return nil
}

func runDoctor() error {
	cfg, err := config.LoadFull()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	fmt.Println("Somon Rent Watcher — production doctor")
	fmt.Println("Version:", version)

	db, err := store.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("SQLite: %w", err)
	}
	defer db.Close()
	count, err := db.CountSeen()
	if err != nil {
		return fmt.Errorf("SQLite count: %w", err)
	}
	fmt.Printf("[OK] SQLite: %s (seen IDs: %d)\n", db.Path(), count)

	tg := telegram.NewClient(cfg.TelegramAPIBase, cfg.TelegramBotToken)
	me, err := tg.GetMe(ctx)
	if err != nil {
		return fmt.Errorf("Telegram getMe: %w", err)
	}
	fmt.Printf("[OK] Telegram bot: @%s, ID %d\n", me.Username, me.ID)
	webhook, err := tg.GetWebhookInfo(ctx)
	if err != nil {
		return fmt.Errorf("Telegram getWebhookInfo: %w", err)
	}
	if webhook.URL != "" {
		fmt.Println("[WARN] Telegram webhook is configured; run mode will remove it without dropping updates")
	} else {
		fmt.Printf("[OK] Telegram long polling available; queued updates: %d\n", webhook.PendingUpdateCount)
	}
	target, err := tg.GetChat(ctx, cfg.TelegramTargetChatID)
	if err != nil {
		return fmt.Errorf("Telegram target chat %d: %w", cfg.TelegramTargetChatID, err)
	}
	fmt.Printf("[OK] Target chat: %s (ID %d, type %s)\n", chatName(target), target.ID, target.Type)

	somonClient := somon.NewClient(cfg.UserAgent, cfg.RequestDelay, cfg.HTTPTimeout, cfg.MaxBodyBytes)
	cards, body, err := somonClient.FetchCategory(ctx, cfg.CategoryURL)
	if err != nil {
		if path, saveErr := saveDoctorHTML(cfg.DebugDir, body); saveErr == nil && path != "" {
			fmt.Println("[DEBUG] Saved Somon HTML:", path)
		}
		return fmt.Errorf("Somon category: %w", err)
	}
	ordinary, promoted := 0, 0
	for _, card := range cards {
		if card.Promoted {
			promoted++
		} else {
			ordinary++
		}
	}
	if len(cards) < cfg.MinCards {
		path, _ := saveDoctorHTML(cfg.DebugDir, body)
		return fmt.Errorf("Somon parser found only %d cards; minimum is %d; HTML saved to %s", len(cards), cfg.MinCards, nonEmptyPath(path))
	}
	if ordinary == 0 {
		path, _ := saveDoctorHTML(cfg.DebugDir, body)
		return fmt.Errorf("Somon parser found no ordinary cards; HTML saved to %s", nonEmptyPath(path))
	}
	fmt.Printf("[OK] Somon category parsed: %d cards (%d ordinary, %d promoted)\n", len(cards), ordinary, promoted)
	for i, card := range cards {
		if i >= 3 {
			break
		}
		fmt.Printf("     ID %d | %s | %s\n", card.ID, pointerInt(card.Price), card.Title)
	}
	fmt.Println("[OK] Doctor finished. It did not create or change the baseline.")
	return nil
}

func saveDoctorHTML(dir string, body []byte) (string, error) {
	if len(body) == 0 {
		return "", nil
	}
	if strings.TrimSpace(dir) == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	path := dir + string(os.PathSeparator) + "doctor-" + time.Now().UTC().Format("20060102T150405Z") + ".html"
	if err := os.WriteFile(path, body, 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func nonEmptyPath(path string) string {
	if strings.TrimSpace(path) == "" {
		return "<not saved>"
	}
	return path
}

func runIDs() error {
	token, apiBase, err := config.LoadTokenOnly()
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	client := telegram.NewClient(apiBase, token)
	me, err := client.GetMe(ctx)
	if err != nil {
		return err
	}
	if err := client.DeleteWebhook(ctx, false); err != nil {
		return err
	}
	fmt.Printf("Bot @%s, ID %d\n", me.Username, me.ID)
	fmt.Println("Send /start to the bot privately and one message in the target group. Incoming IDs will appear below. Ctrl+C to stop.")

	var offset int64
	for {
		updates, err := client.GetUpdates(ctx, offset, 50)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		for _, update := range updates {
			offset = update.UpdateID + 1
			if update.Message != nil {
				from := int64(0)
				username := ""
				if update.Message.From != nil {
					from = update.Message.From.ID
					username = update.Message.From.Username
				}
				fmt.Printf("update=%d user_id=%d username=@%s chat_id=%d chat_type=%s chat=%q text=%q\n",
					update.UpdateID, from, username, update.Message.Chat.ID, update.Message.Chat.Type,
					chatName(update.Message.Chat), truncate(update.Message.Text, 120))
			}
			if update.CallbackQuery != nil {
				chatID := int64(0)
				chatType := ""
				if update.CallbackQuery.Message != nil {
					chatID = update.CallbackQuery.Message.Chat.ID
					chatType = update.CallbackQuery.Message.Chat.Type
				}
				fmt.Printf("update=%d callback user_id=%d username=@%s chat_id=%d chat_type=%s data=%q\n",
					update.UpdateID, update.CallbackQuery.From.ID, update.CallbackQuery.From.Username,
					chatID, chatType, update.CallbackQuery.Data)
			}
		}
	}
}

func newLogger(level string) *slog.Logger {
	var slogLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn", "warning":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
}

func chatName(chat telegram.Chat) string {
	if chat.Title != "" {
		return chat.Title
	}
	name := strings.TrimSpace(chat.FirstName + " " + chat.LastName)
	if name != "" {
		return name
	}
	if chat.Username != "" {
		return "@" + chat.Username
	}
	return strconv.FormatInt(chat.ID, 10)
}

func pointerInt(value *int) string {
	if value == nil {
		return "unknown"
	}
	return strconv.Itoa(*value)
}

func truncate(value string, max int) string {
	runes := []rune(strings.TrimSpace(value))
	if len(runes) <= max {
		return string(runes)
	}
	return string(runes[:max]) + "..."
}

func printHelp() {
	fmt.Println(`Usage: somonwatch [command]

Commands:
  run      start Telegram bot and Somon polling (default)
  doctor   verify environment, SQLite, Telegram and live Somon parser; no baseline changes
  ids      print Telegram user/chat IDs from incoming messages
  version  print build version
  help     show this help`)
}
