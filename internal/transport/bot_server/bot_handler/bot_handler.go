package bot_handler

import (
	"context"
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

	invalidCommandMassage = ""
)

type BotHandler struct {
	usersState usersStateStorage
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

func (b *BotHandler) invalidCommandMassage(chatId int64) (msg *tgbotapi.MessageConfig, err error) {
	msg = pointy.Pointer(tgbotapi.NewMessage(chatId, invalidCommandMassage))
	return msg, nil
}
