package bot_handler

import (
	"context"
	"fmt"
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/openlyinc/pointy"
)

const (
	cmdGetCurrentExchangeWaitConverterWay = 1
)

func (b *BotHandler) cmdGetCurrentExchange(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	if currentState == nil {
		return b.cmdGetCurrentExchangeStart(ctx, update, currentState)
	}
	switch currentState.currentState {
	case cmdGetCurrentExchangeWaitConverterWay:
		return b.cmdGetCurrentExchangeWaitConverterWay(ctx, update, currentState)
	}

	b.removeUsersState(update.Message.Chat.ID)
	return b.internalErrorMassage(update.Message.Chat.ID)
}

func (b *BotHandler) cmdGetCurrentExchangeStart(ctx context.Context, update tgbotapi.Update,
	state *userState) (msg *tgbotapi.MessageConfig, err error) {
	if update.Message == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	converterWay, err := b.converterService.GetMyConverterWay(ctx, update.Message.Chat.ID)
	if err != nil {
		b.removeUsersState(update.Message.Chat.ID)
		return b.internalErrorMassage(update.Message.Chat.ID)
	}

	chooseConverterWay := generateInlineKeyboardConverterWayWithCancel(converterWay)

	msg = pointy.Pointer(tgbotapi.NewMessage(update.Message.Chat.ID, "Choose converter way"))
	msg.ReplyMarkup = chooseConverterWay

	usersState := userState{
		currentState:   cmdGetCurrentExchangeWaitConverterWay,
		commandHandler: b.cmdGetCurrentExchange,
	}

	b.addUserToSateMachine(update.Message.Chat.ID, &usersState)

	return msg, nil
}

func (b *BotHandler) cmdGetCurrentExchangeWaitConverterWay(ctx context.Context,
	update tgbotapi.Update,
	currentState *userState) (*tgbotapi.MessageConfig, error) {
	if update.CallbackQuery == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	converterPair := parseConverterWayStr(update.CallbackQuery.Data)

	exchange, err := b.converterService.GetCurrentExchange(ctx,
		update.CallbackQuery.Message.Chat.ID,
		converterPair)
	if err != nil {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}

	b.removeUsersState(update.CallbackQuery.Message.Chat.ID)

	return pointy.Pointer(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
			fmt.Sprintf("Currnt exchange for\n%s:\n%f (%f)", update.CallbackQuery.Data, exchange,
				1/exchange))),
		nil
}
