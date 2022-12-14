package bot_handler

import (
	"context"
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/openlyinc/pointy"
)

const (
	cmdAddConverterWayWaitConverterWay = 1
)

func (b *BotHandler) cmdAddConverterWay(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig,
	err error) {

	if currentState == nil {
		return b.cmdAddConverterWayStart(ctx, update, currentState)
	}
	switch currentState.currentState {
	case cmdAddConverterWayWaitConverterWay:
		return b.cmdAddConverterWayWaitConverterWay(ctx, update, currentState)
	}

	b.removeUsersState(update.Message.Chat.ID)
	return b.internalErrorMassage(update.Message.Chat.ID)
}

func (b *BotHandler) cmdAddConverterWayStart(ctx context.Context, update tgbotapi.Update,
	state *userState) (msg *tgbotapi.MessageConfig, err error) {
	if update.Message == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	converterWay, err := b.converterService.GetAvailableConverterWay(ctx, update.Message.Chat.ID)
	if err != nil {
		b.removeUsersState(update.Message.Chat.ID)
		return b.internalErrorMassage(update.Message.Chat.ID)
	}

	chooseConverterWay := generateInlineKeyboardConverterWayWithCancel(converterWay)

	msg = pointy.Pointer(tgbotapi.NewMessage(update.Message.Chat.ID, "Choose converter way"))
	msg.ReplyMarkup = chooseConverterWay

	usersState := userState{
		currentState:   cmdAddConverterWayWaitConverterWay,
		commandHandler: b.cmdAddConverterWay,
	}

	b.addUserToSateMachine(update.Message.Chat.ID, &usersState)

	return msg, nil
}

func (b *BotHandler) cmdAddConverterWayWaitConverterWay(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (*tgbotapi.MessageConfig, error) {
	if update.CallbackQuery == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	converterPair := parseConverterWayStr(update.CallbackQuery.Data)

	err := b.converterService.AddUserConverterWay(ctx, update.CallbackQuery.Message.Chat.ID,
		converterPair)
	if err != nil {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}

	b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
	return b.SuccessMassage(update.CallbackQuery.Message.Chat.ID)
}
