package bot_handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) cmdAddCurrency(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	return b.notAvailableCommandMassage(update.Message.Chat.ID)
}
