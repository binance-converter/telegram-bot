package bot_handler

import (
	"context"
	"errors"
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/openlyinc/pointy"
)

const (
	commandStart              = "/start"
	commandAddCurrency        = "/add_currency"
	commandRemoveCurrency     = "/remove_currency"
	commandAddConverterWay    = "/add_converter_way"
	commandRemoveConverterWay = "/remove_converter_way"
	commandGetCurrentExchange = "/get_exchange"
)

const (
	commandStartAnswer = "" +
		"Hi, I'm binance converter bot now I can give you current exchange, but in future we " +
		"will be add more functionality\n" +
		"First you need choose currencies interested for you by command:\n" + commandAddCurrency +
		"\n" +
		"After choose available converter way for your currencies by command:\n" +
		commandAddConverterWay + "\n" +
		"Now you can give current exchange by command:\n" + commandGetCurrentExchange + "\n" +
		"Enjoy!"

	invalidCommandMassage      = "Unknown command"
	notAvailableCommandMassage = "Sorry now this command not available"
	internalErrorMassage       = "Sorry, an internal error has occurred"
	CanceledMassage            = "Command canceled"
	SuccessMassage             = "Success"
)

type AuthService interface {
	SignUp(ctx context.Context, userData core.ServiceSignUpUser) error
}

type CurrencyService interface {
	GetAvailableCurrencies(ctx context.Context, chatId int64,
		currencyType core.CurrencyType) ([]core.CurrencyCode, error)
	GetAvailableBanks(ctx context.Context, chatId int64,
		currencyCode core.CurrencyCode) ([]core.CurrencyBank, error)
	AddUserCurrency(ctx context.Context, chatId int64, currency core.FullCurrency) error
}

type BotHandler struct {
	usersState      usersStateStorage
	authService     AuthService
	currencyService CurrencyService
}

func NewBotHandler(authService AuthService, currencyService CurrencyService) *BotHandler {
	newBotHandler := BotHandler{
		authService:     authService,
		currencyService: currencyService,
	}
	newBotHandler.usersState = make(usersStateStorage)
	return &newBotHandler
}

func (b *BotHandler) CmdHandler(ctx context.Context,
	update tgbotapi.Update) (msg *tgbotapi.MessageConfig, err error) {

	switch update.Message.Text {
	case commandStart:
		msg, err = b.cmdStart(ctx, update, nil)
		break
	case commandAddCurrency:
		msg, err = b.cmdAddCurrency(ctx, update, nil)
		break
	case commandRemoveCurrency:
		msg, err = b.cmdRemoveCurrency(ctx, update, nil)
		break
	case commandAddConverterWay:
		msg, err = b.cmdAddConverterWay(ctx, update, nil)
		break
	case commandRemoveConverterWay:
		msg, err = b.cmdRemoveConverterWay(ctx, update, nil)
		break
	case commandGetCurrentExchange:
		msg, err = b.cmdGetCurrentExchange(ctx, update, nil)
		break
	default:
		msg, err = b.invalidCommandMassage(update.Message.Chat.ID)
	}

	return msg, err
}

func (b *BotHandler) AnswerHandler(ctx context.Context, update tgbotapi.Update) (msg *tgbotapi.
	MessageConfig,
	err error) {
	currentState := b.getUserState(update.Message.Chat.ID)
	if currentState == nil {
		return b.invalidCommandMassage(update.Message.Chat.ID)
	}

	return currentState.commandHandler(ctx, update, currentState)
}

func (b *BotHandler) QueryHandler(ctx context.Context, update tgbotapi.Update) (msg *tgbotapi.
	MessageConfig,
	err error) {
	currentState := b.getUserState(update.CallbackQuery.Message.Chat.ID)
	if currentState == nil {
		return b.invalidCommandMassage(update.CallbackQuery.Message.Chat.ID)
	}

	if update.CallbackQuery.Data == cancelButton {
		b.removeUsersState(update.CallbackQuery.Message.Chat.ID)
		return b.CanceledMassage(update.CallbackQuery.Message.Chat.ID)
	}

	return currentState.commandHandler(ctx, update, currentState)
}

func (b *BotHandler) invalidCommandMassage(chatId int64) (msg *tgbotapi.MessageConfig, err error) {
	msg = pointy.Pointer(tgbotapi.NewMessage(chatId, invalidCommandMassage))
	return msg, nil
}

func (b *BotHandler) notAvailableCommandMassage(chatId int64) (msg *tgbotapi.MessageConfig,
	err error) {
	msg = pointy.Pointer(tgbotapi.NewMessage(chatId, notAvailableCommandMassage))
	return msg, nil
}

func (b *BotHandler) internalErrorMassage(chatId int64) (msg *tgbotapi.MessageConfig,
	err error) {
	msg = pointy.Pointer(tgbotapi.NewMessage(chatId, internalErrorMassage))
	return msg, errors.New("internal error")
}

func (b *BotHandler) CanceledMassage(chatId int64) (msg *tgbotapi.MessageConfig,
	err error) {
	msg = pointy.Pointer(tgbotapi.NewMessage(chatId, CanceledMassage))
	return msg, nil
}

func (b *BotHandler) SuccessMassage(chatId int64) (msg *tgbotapi.MessageConfig,
	err error) {
	msg = pointy.Pointer(tgbotapi.NewMessage(chatId, SuccessMassage))
	return msg, nil
}
