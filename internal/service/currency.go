package service

import (
	"context"
	"errors"
	"github.com/binance-converter/backend-api/api/currencies"
	"github.com/binance-converter/telegram-bot/core"
	"google.golang.org/grpc"
)

type Currency struct {
	currencyServer currencies.CurrenciesClient
}

func NewCurrency(grpcConn grpc.ClientConnInterface) *Currency {
	newCurrency := Currency{}
	newCurrency.currencyServer = currencies.NewCurrenciesClient(grpcConn)
	return &newCurrency
}

func (c Currency) GetAvailableCurrencies(ctx context.Context, chatId int64,
	currencyType core.CurrencyType) ([]core.CurrencyCode, error) {
	currencyTypeProto, err := convertCoreCurrencyTypeToProto(currencyType)
	if err != nil {
		return nil, err
	}

	ctxWithChatId := addChatIdToContext(ctx, chatId)

	currencyCodes, err := c.currencyServer.GetAvailableCurrencies(ctxWithChatId, currencyTypeProto)
	if err != nil {
		return nil, err
	}
	return convertProtoCurrencyCodesToCore(currencyCodes)
}

func (c Currency) GetAvailableBanks(ctx context.Context, chatId int64,
	currencyCode core.CurrencyCode) ([]core.CurrencyBank, error) {

	ctxWithChatId := addChatIdToContext(ctx, chatId)

	banks, err := c.currencyServer.GetAvailableBankByCurrency(ctxWithChatId,
		convertCoreCurrencyCodeToProto(currencyCode))
	if err != nil {
		return nil, err
	}
	return convertProtoCurrencyBanksToCore(banks)
}

func (c Currency) AddUserCurrency(ctx context.Context, chatId int64,
	currency core.FullCurrency) error {

	protoCurrency, err := convertCoreFullCurrencyToProto(currency)
	if err != nil {
		return err
	}

	ctxWithChatId := addChatIdToContext(ctx, chatId)

	_, err = c.currencyServer.SetCurrency(ctxWithChatId, protoCurrency)

	return err
}

// ------------------------------------------------------------------------------------------------
// helper functions

func convertProtoCurrencyTypeToCore(currencyType *currencies.CurrencyType) (core.CurrencyType,
	error) {
	if currencyType == nil {
		return core.CurrencyType(0), core.ErrorCurrencyEmptyInputArg
	}
	switch currencyType.Type {
	case currencies.ECurrencyType_CRYPTO:
		return core.CurrencyTypeCrypto, nil
	case currencies.ECurrencyType_CLASSIC:
		return core.CurrencyTypeClassic, nil
	}
	return 0, core.ErrorCurrencyInvalidCurrencyType
}

func convertCoreCurrencyTypeToProto(currencyType core.CurrencyType) (*currencies.CurrencyType,
	error) {
	switch currencyType {
	case core.CurrencyTypeCrypto:
		return &currencies.CurrencyType{
			Type: currencies.ECurrencyType_CRYPTO,
		}, nil
	case core.CurrencyTypeClassic:
		return &currencies.CurrencyType{
			Type: currencies.ECurrencyType_CLASSIC,
		}, nil
	}
	return nil, errors.New("error parsing currency type")
}

func convertProtoCurrencyCodeToCore(protoCurrencyCode *currencies.CurrencyCode) (
	core.CurrencyCode, error) {
	if protoCurrencyCode == nil {
		return "", core.ErrorCurrencyEmptyInputArg
	}
	currencyCode := core.CurrencyCode(protoCurrencyCode.CurrencyCode)
	return currencyCode, nil
}

func convertProtoCurrencyCodesToCore(protoCurrencyCodes *currencies.CurrencyCodes) (
	[]core.CurrencyCode, error) {
	if protoCurrencyCodes == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}

	var coreCurrencyCodes []core.CurrencyCode

	for _, currencyCode := range protoCurrencyCodes.CurrencyCodes {
		coreCurrencyCode, err := convertProtoCurrencyCodeToCore(currencyCode)
		if err != nil {
			return nil, err
		}
		coreCurrencyCodes = append(coreCurrencyCodes, coreCurrencyCode)
	}

	return coreCurrencyCodes, nil
}

func convertCoreCurrencyCodeToProto(coreCurrencyCode core.CurrencyCode) (
	currencyCode *currencies.CurrencyCode) {
	currencyCode = &currencies.CurrencyCode{}
	currencyCode.CurrencyCode = string(coreCurrencyCode)
	return currencyCode
}

func convertCoreCurrencyCodesToProto(coreCurrencies []core.CurrencyCode) (currencyCodes *currencies.
	CurrencyCodes) {
	currencyCodes = &currencies.CurrencyCodes{}
	for _, currency := range coreCurrencies {
		currencyCodes.CurrencyCodes = append(currencyCodes.CurrencyCodes,
			convertCoreCurrencyCodeToProto(currency))
	}
	return currencyCodes
}

func convertCoreCurrencyBankToProto(coreCurrencyBank core.CurrencyBank) (CurrencyBank *currencies.
	BankName) {
	CurrencyBank = &currencies.BankName{}
	CurrencyBank.BankName = string(coreCurrencyBank)
	return CurrencyBank
}

func convertProtoCurrencyBankToCore(protoCurrencyBank *currencies.
	BankName) (core.CurrencyBank, error) {
	if protoCurrencyBank == nil {
		return "", core.ErrorCurrencyEmptyInputArg
	}
	return core.CurrencyBank(protoCurrencyBank.BankName), nil
}

func convertProtoCurrencyBanksToCore(protoCurrencyBank *currencies.
	BankNames) ([]core.CurrencyBank, error) {
	if protoCurrencyBank == nil {
		return nil, core.ErrorCurrencyEmptyInputArg
	}
	var coreBanks []core.CurrencyBank

	for _, bankName := range protoCurrencyBank.BankNames {
		coreBank, err := convertProtoCurrencyBankToCore(bankName)
		if err != nil {
			return nil, err
		}
		coreBanks = append(coreBanks, coreBank)
	}
	return coreBanks, nil
}

func convertCoreCurrencyBanksToProto(coreCurrencyBanks []core.CurrencyBank) (
	CurrencyBank *currencies.BankNames) {
	CurrencyBank = &currencies.BankNames{}
	for _, coreCurrencyBank := range coreCurrencyBanks {
		bankName := convertCoreCurrencyBankToProto(coreCurrencyBank)
		CurrencyBank.BankNames = append(CurrencyBank.BankNames,
			bankName)
	}
	return CurrencyBank
}

func convertProtoFullCurrencyToCore(protoCurrency *currencies.FullCurrency) (core.
	FullCurrency, error) {

	if protoCurrency == nil {
		return core.FullCurrency{}, core.ErrorCurrencyEmptyInputArg
	}

	coreCurrencyType, err := convertProtoCurrencyTypeToCore(protoCurrency.Type)
	if err != nil {
		return core.FullCurrency{}, err
	}

	currencyCode, err := convertProtoCurrencyCodeToCore(protoCurrency.CurrencyCode)
	if err != nil {
		return core.FullCurrency{}, err
	}

	bankCode, err := convertProtoCurrencyBankToCore(protoCurrency.BankName)
	if err != nil {
		return core.FullCurrency{}, err
	}

	return core.FullCurrency{
		CurrencyType: coreCurrencyType,
		CurrencyCode: currencyCode,
		BankCode:     bankCode,
	}, nil
}

func convertCoreFullCurrencyToProto(coreCurrency core.FullCurrency) (*currencies.
	FullCurrency, error) {
	currencyType, err := convertCoreCurrencyTypeToProto(coreCurrency.CurrencyType)
	if err != nil {
		return nil, err
	}
	return &currencies.FullCurrency{
		Type:         currencyType,
		CurrencyCode: convertCoreCurrencyCodeToProto(coreCurrency.CurrencyCode),
		BankName:     convertCoreCurrencyBankToProto(coreCurrency.BankCode),
	}, nil
}

func convertCoreFullCurrenciesToProto(coreCurrencies []core.FullCurrency) (protoCurrencies *currencies.
	FullCurrencies, err error) {
	protoCurrencies = &currencies.FullCurrencies{}
	for _, currency := range coreCurrencies {
		protoCurrency, err := convertCoreFullCurrencyToProto(currency)
		if err != nil {
			return nil, err
		}
		protoCurrencies.FullCurrencies = append(protoCurrencies.FullCurrencies, protoCurrency)
	}
	return protoCurrencies, nil
}
