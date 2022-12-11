package bot_handler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/net/context"
	"time"
)

type commandHandler func(context.Context, tgbotapi.Update,
	*userState) (*tgbotapi.MessageConfig,
	error)

type userState struct {
	commandHandler commandHandler
	currentState   string
	timeout        time.Time
}

type usersStateStorage map[int64]*userState

func (b *BotHandler) addUserToSateMachine(chatId int64, userState *userState) {
	b.usersState[chatId] = userState
}

func (b *BotHandler) getUserState(chatId int64) *userState {
	if userState, ok := b.usersState[chatId]; ok {
		return userState
	}
	return nil
}
