package bot_handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) cmdStart(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {

	massage := tgbotapi.NewMessage(update.Message.Chat.ID, commandStartAnswer)
	return &massage, err
}
