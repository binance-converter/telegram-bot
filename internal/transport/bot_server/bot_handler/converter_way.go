package bot_handler

import (
	"fmt"
	"github.com/binance-converter/telegram-bot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"regexp"
	"strings"
)

const (
	converterWayStrStep = "->"
)

func generateInlineKeyboardConverterWayWithCancel(converterWays []core.ConverterPair) tgbotapi.
	InlineKeyboardMarkup {
	var converterWaysStr []string

	for _, converterWay := range converterWays {
		var way []string
		for _, currency := range converterWay.Currencies {
			if currency.CurrencyType == core.CurrencyTypeCrypto {
				way = append(way, fmt.Sprintf("%s", currency.CurrencyCode))
			} else {
				way = append(way, fmt.Sprintf("%s(%s)", currency.CurrencyCode, currency.BankCode))
			}
		}
		converterWaysStr = append(converterWaysStr, strings.Join(way, converterWayStrStep))
	}

	return generateInlineKeyboardWithCancel(converterWaysStr)
}

func parseConverterWayStr(converterWayStr string) core.ConverterPair {
	currenciesStr := strings.Split(converterWayStr, converterWayStrStep)

	var converterWay core.ConverterPair

	for _, currencyStr := range currenciesStr {
		currency, err := convertStrToFullCurrency(currencyStr)
		if err != nil {
			continue
		}
		converterWay.Currencies = append(converterWay.Currencies, currency)
	}

	return converterWay
}

func convertStrToFullCurrency(str string) (core.FullCurrency, error) {
	r, err := regexp.Compile(
		`^(?P<currency_code>\w+)(\((?P<bank_name>\w+)\)$|$)`,
	)

	currency := core.FullCurrency{}

	if err != nil {
		return currency, err
	}

	m := r.FindStringSubmatch(str)
	for i, name := range r.SubexpNames() {
		switch name {
		case "currency_code":
			currency.CurrencyCode = core.CurrencyCode(m[i])
			break
		case "bank_name":
			currency.BankCode = core.CurrencyBank(m[i])
			break
		}
	}

	if currency.BankCode == "" {
		currency.CurrencyType = core.CurrencyTypeCrypto
	} else {
		currency.CurrencyType = core.CurrencyTypeClassic
	}

	return currency, nil
}
