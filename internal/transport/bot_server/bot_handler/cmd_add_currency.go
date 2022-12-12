package bot_handler

import (
	"context"
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/openlyinc/pointy"
)

const (
	cmdAddCurrencyStateWaitCurrencyType = 1
	cmdAddCurrencyStateWaitCurrencyCode = 2
	cmdAddCurrencyStateWaitBankCode     = 3
)

func (b *BotHandler) cmdAddCurrency(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	if currentState == nil {
		return b.cmdAddCurrencyStart(ctx, update, currentState)
	}
	switch currentState.currentState {
	case cmdAddCurrencyStateWaitCurrencyType:
		return b.cmdAddCurrencyWaitCurrencyType(ctx, update, currentState)
	case cmdAddCurrencyStateWaitCurrencyCode:
		return b.cmdAddCurrencyWaitCurrencyCode(ctx, update, currentState)
	case cmdAddCurrencyStateWaitBankCode:
		return b.cmdAddCurrencyWaitWaitBankCode(ctx, update, currentState)
	}

	b.removeUsersState(update.Message.Chat.ID)
	return b.internalErrorMassage(update.Message.Chat.ID)
}

func (b *BotHandler) cmdAddCurrencyStart(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	if update.Message == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}
	chooseCurrencyType := generateInlineKeyboardWithCancel([]string{core.CurrencyTypeCryptoLabel,
		core.CurrencyTypeClassicLabel})

	msg = pointy.Pointer(tgbotapi.NewMessage(update.Message.Chat.ID, "Choose currency type"))
	msg.ReplyMarkup = chooseCurrencyType

	usersState := userState{
		currentState:   cmdAddCurrencyStateWaitCurrencyType,
		commandHandler: b.cmdAddCurrency,
		commandData:    &core.FullCurrency{},
	}

	b.addUserToSateMachine(update.Message.Chat.ID, &usersState)

	return msg, nil
}

func (b *BotHandler) cmdAddCurrencyWaitCurrencyType(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	if update.CallbackQuery == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	currency, ok := currentState.commandData.(*core.FullCurrency)
	if !ok {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}

	switch update.CallbackQuery.Data {
	case core.CurrencyTypeCryptoLabel:
		currency.CurrencyType = core.CurrencyTypeCrypto
		break
	case core.CurrencyTypeClassicLabel:
		currency.CurrencyType = core.CurrencyTypeClassic
		break
	default:
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return nil, core.ErrorCurrencyInvalidCurrencyType
	}

	currencyCodes, err := b.currencyService.GetAvailableCurrencies(ctx,
		update.CallbackQuery.Message.Chat.ID,
		currency.CurrencyType)
	if err != nil {
		return nil, err
	}

	chooseCurrencyType := generateInlineKeyboardWithCancel(currencyCodes)

	msg = pointy.Pointer(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Choose currency code"))
	msg.ReplyMarkup = chooseCurrencyType
	currentState.setCurrentState(cmdAddCurrencyStateWaitCurrencyCode)

	return msg, nil
}

func (b *BotHandler) cmdAddCurrencyWaitCurrencyCode(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	if update.CallbackQuery == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	currency, ok := currentState.commandData.(*core.FullCurrency)
	if !ok {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}

	currency.CurrencyCode = core.CurrencyCode(update.CallbackQuery.Data)

	if currency.CurrencyType == core.CurrencyTypeCrypto {
		err := b.currencyService.AddUserCurrency(ctx, update.CallbackQuery.Message.Chat.ID,
			*currency)
		if err != nil {
			b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
			return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
		}
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.SuccessMassage(update.CallbackQuery.Message.Chat.ID)
	}
	banks, err := b.currencyService.GetAvailableBanks(ctx,
		update.CallbackQuery.Message.Chat.ID, currency.CurrencyCode)
	if err != nil {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}

	chooseBank := generateInlineKeyboardWithCancel(banks)
	msg = pointy.Pointer(tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
		"Choose currency code"))
	msg.ReplyMarkup = chooseBank
	currentState.setCurrentState(cmdAddCurrencyStateWaitBankCode)
	return msg, nil
}

func (b *BotHandler) cmdAddCurrencyWaitWaitBankCode(ctx context.Context, update tgbotapi.Update,
	currentState *userState) (msg *tgbotapi.MessageConfig, err error) {
	if update.CallbackQuery == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	currency, ok := currentState.commandData.(*core.FullCurrency)
	if !ok {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}

	currency.BankCode = core.CurrencyBank(update.CallbackQuery.Data)

	err = b.currencyService.AddUserCurrency(ctx, update.CallbackQuery.Message.Chat.ID,
		*currency)
	if err != nil {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.internalErrorMassage(update.CallbackQuery.Message.Chat.ID)
	}
	b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
	return b.SuccessMassage(update.CallbackQuery.Message.Chat.ID)
}
