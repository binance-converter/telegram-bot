package bot_handler

import (
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cancelButton = "Cancel"
)

func generateInlineKeyboardWithCancel[T string | core.CurrencyCode | core.CurrencyBank](buttons []T) tgbotapi.
	InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton
	buttons = append(buttons, cancelButton)
	for _, button := range buttons {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			string(button), string(button))))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
