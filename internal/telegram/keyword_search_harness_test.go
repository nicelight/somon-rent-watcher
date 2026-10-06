package telegram

import "context"

// ProcessKeywordUpdateForTest lets external integration tests use the real
// update entrypoint and App without introducing a production import cycle.
func ProcessKeywordUpdateForTest(b *Bot, ctx context.Context, update Update) error {
	return b.processUpdate(ctx, update)
}
