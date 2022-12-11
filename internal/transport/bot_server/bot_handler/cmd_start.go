package bot_handler

import (
	"context"
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *BotHandler) cmdStart(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {

	signUpData := core.ServiceSignUpUser{
		ChatId:       update.Message.Chat.ID,
		UserName:     update.Message.From.UserName,
		FirstName:    update.Message.From.FirstName,
		LastName:     update.Message.From.LastName,
		LanguageCode: update.Message.From.LanguageCode,
	}
	err = b.authService.SignUp(ctx, signUpData)
	if err != nil {
		switch err {
		case core.ErrorAuthServiceAuthUserAlreadyExists:
			break
		default:
			massage := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			return &massage, err
		}
	}

	massage := tgbotapi.NewMessage(update.Message.Chat.ID, commandStartAnswer)
	return &massage, err
}
