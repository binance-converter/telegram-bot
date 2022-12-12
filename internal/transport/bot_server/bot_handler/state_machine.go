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
	currentState   int
	timeout        time.Time
	commandData    interface{}
}

func (u *userState) setHandler(commandHandler commandHandler) {
	u.commandHandler = commandHandler
}

func (u *userState) setCurrentState(currentState int) {
	u.currentState = currentState
}

func (u *userState) setTimeout(timeout time.Time) {
	u.timeout = timeout
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

func (b *BotHandler) removeUsersState(chatId int64) {
	delete(b.usersState, chatId)
}
