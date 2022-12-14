package service

import (
	"context"
	"github.com/binance-converter/backend-api/api/converter"
	"github.com/binance-converter/telegram-bot/core"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Converter struct {
	converterServer converter.ConverterClient
}

func NewConverter(grpcConn grpc.ClientConnInterface) *Converter {
	newConverter := Converter{}
	newConverter.converterServer = converter.NewConverterClient(grpcConn)
	return &newConverter
}

func (c *Converter) GetAvailableConverterWay(ctx context.Context,
	chatId int64) ([]core.ConverterPair, error) {

	ctxWithChatId := addChatIdToContext(ctx, chatId)

	protoConverterPairs, err := c.converterServer.GetAvailableConverterPairs(ctxWithChatId,
		&emptypb.Empty{})
	if err != nil {
		return nil, err
	}
	converterPairs, err := convertProtoConverterPairsToCore(protoConverterPairs)
	if err != nil {
		return nil, err
	}
	return converterPairs, nil
}

func (c *Converter) AddUserConverterWay(ctx context.Context, chatId int64,
	converterPair core.ConverterPair) error {
	protoConverterPair, err := convertCoreConverterPairToProto(converterPair)
	if err != nil {
		return err
	}
	ctxWithChatId := addChatIdToContext(ctx, chatId)
	_, err = c.converterServer.SetConvertPair(ctxWithChatId, protoConverterPair)

	return err
}

// ------------------------------------------------------------------------------------------------
// helper functions

func convertCoreConverterPairToProto(corePair core.ConverterPair) (*converter.ConverterPair,
	error) {
	pair := &converter.ConverterPair{}
	for _, currency := range corePair.Currencies {
		protoCurrency, err := convertCoreFullCurrencyToProto(currency)
		if err != nil {
			return nil, err
		}
		pair.ConverterPair = append(pair.ConverterPair, protoCurrency)
	}
	return pair, nil
}

func convertCoreConverterPairsToProto(corePairs []core.ConverterPair) (*converter.ConverterPairs,
	error) {
	pairs := &converter.ConverterPairs{}
	for _, corePair := range corePairs {
		pair, err := convertCoreConverterPairToProto(corePair)
		if err != nil {
			return nil, err
		}
		pairs.ConverterPairs = append(pairs.ConverterPairs, pair)
	}

	return pairs, nil
}

func convertProtoConverterPairToCore(protoConverterPair *converter.ConverterPair) (core.
	ConverterPair, error) {
	if protoConverterPair == nil {
		return core.ConverterPair{}, core.ErrorConverterEmptyInputArg
	}
	coreConverterPair := core.ConverterPair{}
	for _, protoCurrency := range protoConverterPair.ConverterPair {
		coreCurrency, err := convertProtoFullCurrencyToCore(protoCurrency)
		if err != nil {
			return coreConverterPair, err
		}
		coreConverterPair.Currencies = append(coreConverterPair.Currencies, coreCurrency)
	}
	return coreConverterPair, nil
}
func convertProtoConverterPairsToCore(protoConverterPairs *converter.ConverterPairs) ([]core.
	ConverterPair, error) {
	if protoConverterPairs == nil {
		return nil, core.ErrorConverterEmptyInputArg
	}
	var converterPairs []core.ConverterPair

	for _, protoConverterPair := range protoConverterPairs.ConverterPairs {
		coreConverterPair, err := convertProtoConverterPairToCore(protoConverterPair)
		if err != nil {
			return nil, err
		}
		converterPairs = append(converterPairs, coreConverterPair)
	}
	return converterPairs, nil
}

func convertCoreExchangeToProto(coreExchange core.Exchange) *converter.Exchange {
	return &converter.Exchange{
		Exchange: float32(coreExchange),
	}
}

func convertCoreThresholdConverterPairToProto(coreThreshold core.ThresholdConvertPair) (
	*converter.ThresholdConvertPair, error) {
	threshold := &converter.ThresholdConvertPair{
		Exchange: convertCoreExchangeToProto(coreThreshold.Exchange),
	}
	converterPair, err := convertCoreConverterPairToProto(coreThreshold.ConverterPair)
	if err != nil {
		return nil, err
	}
	threshold.ConverterPair = converterPair
	return threshold, nil
}

func convertCoreThresholdConverterPairsToProto(coreThreshold []core.ThresholdConvertPair) (
	*converter.ThresholdConvertPairs, error) {
	threshold := &converter.ThresholdConvertPairs{}
	for _, coreThresholdPair := range coreThreshold {
		converterPair, err := convertCoreThresholdConverterPairToProto(coreThresholdPair)
		if err != nil {
			return nil, err
		}
		threshold.ConverterPairs = append(threshold.ConverterPairs, converterPair)
	}
	return threshold, nil
}

func convertProtoExchangeToCore(protoExchange *converter.Exchange) (core.Exchange, error) {
	if protoExchange == nil {
		return core.Exchange(0), core.ErrorConverterEmptyInputArg
	}
	return core.Exchange(protoExchange.Exchange), nil
}

func convertProtoThresholdConverterPair(protoThreshold *converter.ThresholdConvertPair) (core.
	ThresholdConvertPair, error) {
	coreThreshold := core.ThresholdConvertPair{}

	var err error

	coreThreshold.Exchange, err = convertProtoExchangeToCore(protoThreshold.Exchange)
	if err != nil {
		return coreThreshold, err
	}

	coreThreshold.ConverterPair, err = convertProtoConverterPairToCore(protoThreshold.
		ConverterPair)

	return coreThreshold, err
}
