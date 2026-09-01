package somon

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/nicelight/somon-rent-watcher/internal/model"
)

const DefaultCategoryURL = "https://somon.tj/nedvizhimost/arenda-kvartir/dushanbe/"

var recoveryRoomPaths = map[string]string{
	"1":  "/nedvizhimost/arenda-kvartir/1-komnatnyie/dushanbe/",
	"2":  "/nedvizhimost/arenda-kvartir/2-komnatnyie/dushanbe/",
	"3":  "/nedvizhimost/arenda-kvartir/3-komnatnyie/dushanbe/",
	"4":  "/nedvizhimost/arenda-kvartir/4-komnatnyie/dushanbe/",
	"5":  "/nedvizhimost/arenda-kvartir/5-komnatnyie/dushanbe/",
	"6+": "/nedvizhimost/arenda-kvartir/6-i-bolee-komnat/dushanbe/",
}

type BlockedPageError struct {
	URL    string
	Reason string
}

func (e *BlockedPageError) Error() string {
	if e.Reason == "" {
		return fmt.Sprintf("blocked/anti-bot page from %s", e.URL)
	}
	return fmt.Sprintf("blocked/anti-bot page from %s: %s", e.URL, e.Reason)
}

type HTTPError struct {
	URL        string
	StatusCode int
	RetryAfter time.Duration
}

func (e *HTTPError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("HTTP %d from %s (retry after %s)", e.StatusCode, e.URL, e.RetryAfter)
	}
	return fmt.Sprintf("HTTP %d from %s", e.StatusCode, e.URL)
}

func IsBlocked(err error) (status int, retryAfter time.Duration, ok bool) {
	var blockedPage *BlockedPageError
	if errors.As(err, &blockedPage) {
		return http.StatusForbidden, 0, true
	}
	var httpErr *HTTPError
	if !errors.As(err, &httpErr) {
		return 0, 0, false
	}
	if httpErr.StatusCode != http.StatusForbidden && httpErr.StatusCode != http.StatusTooManyRequests {
		return httpErr.StatusCode, httpErr.RetryAfter, false
	}
	return httpErr.StatusCode, httpErr.RetryAfter, true
}

type Client struct {
	httpClient   *http.Client
	userAgent    string
	requestDelay time.Duration
	maxBodyBytes int64

	mu          sync.Mutex
	lastRequest time.Time
}

func NewClient(userAgent string, requestDelay, timeout time.Duration, maxBodyBytes int64) *Client {
	if strings.TrimSpace(userAgent) == "" {
		userAgent = "SomonRentWatcher/1.0 (private notifier)"
	}
	if timeout <= 0 {
		timeout = 25 * time.Second
	}
	if maxBodyBytes <= 0 {
		maxBodyBytes = 10 << 20
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		MaxIdleConns:          10,
		MaxIdleConnsPerHost:   4,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: timeout,
		ForceAttemptHTTP2:     true,
	}
	return &Client{
		httpClient:   &http.Client{Transport: transport, Timeout: timeout},
		userAgent:    userAgent,
		requestDelay: requestDelay,
		maxBodyBytes: maxBodyBytes,
	}
}

func (c *Client) FetchCategory(ctx context.Context, pageURL string) ([]model.Card, []byte, error) {
	body, err := c.get(ctx, pageURL, "")
	if err != nil {
		return nil, nil, err
	}
	cards, err := ParseCategory(pageURL, body)
	if err != nil {
		return nil, body, err
	}
	return cards, body, nil
}

func (c *Client) FetchDetail(ctx context.Context, card model.Card, referer string) (model.Ad, []byte, error) {
	body, err := c.get(ctx, card.URL, referer)
	if err != nil {
		return model.Ad{}, nil, err
	}
	ad, err := ParseDetail(card.URL, body, card)
	if err != nil {
		return model.Ad{}, body, err
	}
	return ad, body, nil
}

func (c *Client) GetRaw(ctx context.Context, pageURL string) ([]byte, error) {
	return c.get(ctx, pageURL, "")
}

func (c *Client) get(ctx context.Context, pageURL, referer string) ([]byte, error) {
	if err := c.wait(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ru,en;q=0.7")
	if referer != "" {
		req.Header.Set("Referer", referer)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", pageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			URL:        pageURL,
			StatusCode: resp.StatusCode,
			RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After")),
		}
	}

	limited := io.LimitReader(resp.Body, c.maxBodyBytes+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", pageURL, err)
	}
	if int64(len(body)) > c.maxBodyBytes {
		return nil, fmt.Errorf("response from %s exceeds %d bytes", pageURL, c.maxBodyBytes)
	}
	return body, nil
}

func (c *Client) wait(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.requestDelay <= 0 || c.lastRequest.IsZero() {
		c.lastRequest = time.Now()
		return nil
	}
	wait := time.Until(c.lastRequest.Add(c.requestDelay))
	if wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
		}
	}
	c.lastRequest = time.Now()
	return nil
}

func parseRetryAfter(value string) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	if seconds, err := strconv.Atoi(value); err == nil && seconds >= 0 {
		return time.Duration(seconds) * time.Second
	}
	if when, err := http.ParseTime(value); err == nil {
		if d := time.Until(when); d > 0 {
			return d
		}
	}
	return 0
}

func RecoveryURLs(categoryURL string, rooms []string) []string {
	base := strings.TrimSuffix(categoryURL, "/")
	if idx := strings.Index(base, "/nedvizhimost/"); idx >= 0 {
		base = base[:idx]
	}
	seen := make(map[string]struct{})
	urls := make([]string, 0, len(rooms))
	for _, room := range rooms {
		path, ok := recoveryRoomPaths[room]
		if !ok {
			continue
		}
		u := base + path
		if _, exists := seen[u]; exists {
			continue
		}
		seen[u] = struct{}{}
		urls = append(urls, u)
	}
	return urls
}
