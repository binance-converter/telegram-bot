package core

import "errors"

type CurrencyType int32

const (
	CurrencyTypeCrypto  CurrencyType = 0
	CurrencyTypeClassic CurrencyType = 1
)

const (
	CurrencyTypeCryptoLabel  = "CRYPTO"
	CurrencyTypeClassicLabel = "CLASSIC"
)

type CurrencyCode string
type CurrencyBank string

type FullCurrency struct {
	CurrencyType CurrencyType
	CurrencyCode CurrencyCode
	BankCode     CurrencyBank
}

var (
	ErrorCurrencyEmptyInputArg       = errors.New("empty input arguments")
	ErrorCurrencyInvalidCurrencyType = errors.New("invalid currency type")
	ErrorCurrencyInvalidCurrencyCode = errors.New("invalid currency code")
	ErrorCurrencyInvalidBankCode     = errors.New("invalid bank code")
	ErrorCurrencyInternal            = errors.New("internal error")
	ErrorCurrencyNotAuthorized       = errors.New("not authorized")
	ErrorCurrencyAlreadyHas          = errors.New("currency already has")
	ErrorCurrencyNotFound            = errors.New("currency not found")
)
