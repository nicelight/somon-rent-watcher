package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(apiBase, token string) *Client {
	return &Client{
		baseURL: strings.TrimRight(apiBase, "/") + "/bot" + token,
		httpClient: &http.Client{
			Timeout: 70 * time.Second,
		},
	}
}

type APIError struct {
	Code        int
	Description string
	RetryAfter  int
}

func (e *APIError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("Telegram API %d: %s (retry after %ds)", e.Code, e.Description, e.RetryAfter)
	}
	return fmt.Sprintf("Telegram API %d: %s", e.Code, e.Description)
}

type apiResponse[T any] struct {
	OK          bool   `json:"ok"`
	Result      T      `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Description string `json:"description"`
	Parameters  struct {
		RetryAfter int `json:"retry_after"`
	} `json:"parameters"`
}

type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Message struct {
	MessageID int64  `json:"message_id"`
	From      *User  `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

type CallbackQuery struct {
	ID      string   `json:"id"`
	From    User     `json:"from"`
	Message *Message `json:"message"`
	Data    string   `json:"data"`
}

type Update struct {
	UpdateID      int64          `json:"update_id"`
	Message       *Message       `json:"message"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
}

type WebhookInfo struct {
	URL                  string `json:"url"`
	HasCustomCertificate bool   `json:"has_custom_certificate"`
	PendingUpdateCount   int    `json:"pending_update_count"`
	LastErrorDate        int64  `json:"last_error_date"`
	LastErrorMessage     string `json:"last_error_message"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data,omitempty"`
	URL          string `json:"url,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (c *Client) GetMe(ctx context.Context) (User, error) {
	return doForm[User](ctx, c, "getMe", nil)
}

func (c *Client) DeleteWebhook(ctx context.Context, dropPendingUpdates bool) error {
	form := url.Values{}
	if dropPendingUpdates {
		form.Set("drop_pending_updates", "true")
	}
	_, err := doForm[bool](ctx, c, "deleteWebhook", form)
	return err
}

func (c *Client) GetWebhookInfo(ctx context.Context) (WebhookInfo, error) {
	return doForm[WebhookInfo](ctx, c, "getWebhookInfo", nil)
}

func (c *Client) GetChat(ctx context.Context, chatID int64) (Chat, error) {
	form := url.Values{"chat_id": {strconv.FormatInt(chatID, 10)}}
	return doForm[Chat](ctx, c, "getChat", form)
}

func (c *Client) GetUpdates(ctx context.Context, offset int64, timeoutSeconds int) ([]Update, error) {
	allowed, _ := json.Marshal([]string{"message", "callback_query"})
	form := url.Values{
		"offset":          {strconv.FormatInt(offset, 10)},
		"timeout":         {strconv.Itoa(timeoutSeconds)},
		"allowed_updates": {string(allowed)},
	}
	return doForm[[]Update](ctx, c, "getUpdates", form)
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string, keyboard *InlineKeyboardMarkup) (Message, error) {
	form := url.Values{
		"chat_id":                  {strconv.FormatInt(chatID, 10)},
		"text":                     {text},
		"parse_mode":               {"HTML"},
		"disable_web_page_preview": {"true"},
	}
	if keyboard != nil {
		data, err := json.Marshal(keyboard)
		if err != nil {
			return Message{}, err
		}
		form.Set("reply_markup", string(data))
	}
	return doForm[Message](ctx, c, "sendMessage", form)
}

func (c *Client) SendPhoto(ctx context.Context, chatID int64, photoURL, caption string, keyboard *InlineKeyboardMarkup) (Message, error) {
	form := url.Values{
		"chat_id":    {strconv.FormatInt(chatID, 10)},
		"photo":      {photoURL},
		"caption":    {caption},
		"parse_mode": {"HTML"},
	}
	if keyboard != nil {
		data, err := json.Marshal(keyboard)
		if err != nil {
			return Message{}, err
		}
		form.Set("reply_markup", string(data))
	}
	return doForm[Message](ctx, c, "sendPhoto", form)
}

func (c *Client) EditMessageText(ctx context.Context, chatID, messageID int64, text string, keyboard *InlineKeyboardMarkup) (Message, error) {
	form := url.Values{
		"chat_id":                  {strconv.FormatInt(chatID, 10)},
		"message_id":               {strconv.FormatInt(messageID, 10)},
		"text":                     {text},
		"parse_mode":               {"HTML"},
		"disable_web_page_preview": {"true"},
	}
	if keyboard != nil {
		data, err := json.Marshal(keyboard)
		if err != nil {
			return Message{}, err
		}
		form.Set("reply_markup", string(data))
	}
	return doForm[Message](ctx, c, "editMessageText", form)
}

func (c *Client) AnswerCallbackQuery(ctx context.Context, callbackID, text string, alert bool) error {
	form := url.Values{"callback_query_id": {callbackID}}
	if text != "" {
		form.Set("text", text)
	}
	if alert {
		form.Set("show_alert", "true")
	}
	_, err := doForm[bool](ctx, c, "answerCallbackQuery", form)
	return err
}

func doForm[T any](ctx context.Context, c *Client, method string, form url.Values) (T, error) {
	var zero T
	if form == nil {
		form = make(url.Values)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/"+method, strings.NewReader(form.Encode()))
	if err != nil {
		return zero, fmt.Errorf("create Telegram %s request: %s", method, c.redactError(err))
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("Telegram %s request failed: %s", method, c.redactError(err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return zero, fmt.Errorf("read Telegram %s: %w", method, err)
	}
	var decoded apiResponse[T]
	if err := json.Unmarshal(body, &decoded); err != nil {
		return zero, fmt.Errorf("decode Telegram %s response (HTTP %d): %w", method, resp.StatusCode, err)
	}
	if !decoded.OK {
		return zero, &APIError{Code: decoded.ErrorCode, Description: decoded.Description, RetryAfter: decoded.Parameters.RetryAfter}
	}
	return decoded.Result, nil
}

func (c *Client) redactError(err error) string {
	if err == nil {
		return ""
	}
	message := err.Error()
	if c.baseURL != "" {
		message = strings.ReplaceAll(message, c.baseURL, "<telegram-bot-api>")
	}
	return message
}
