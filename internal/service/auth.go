package service

import (
	"context"
	"github.com/binance-converter/backend-api/api/auth"
	"github.com/binance-converter/telegram-bot/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth struct {
	authServer auth.AuthClient
}

func NewAuth(grpcConn grpc.ClientConnInterface) *Auth {
	newAuth := Auth{}
	newAuth.authServer = auth.NewAuthClient(grpcConn)
	return &newAuth
}

func (a Auth) SignUp(ctx context.Context, userData core.ServiceSignUpUser) error {
	_, err := a.authServer.SignUpUserByTelegram(ctx, convertCoreServiceSignUpUserToProto(userData))
	if err != nil {
		if reqStatus, ok := status.FromError(err); ok {
			switch reqStatus.Code() {
			case codes.AlreadyExists:
				return core.ErrorAuthServiceAuthUserAlreadyExists
			default:
				return core.ErrorAuthServiceInternalError
			}
		} else {
			return err
		}
	}
	return nil
}

func convertCoreServiceSignUpUserToProto(userData core.ServiceSignUpUser) *auth.
	SignUpUserByTelegramRequest {
	return &auth.SignUpUserByTelegramRequest{
		ChatId:       userData.ChatId,
		UserName:     userData.UserName,
		FirstName:    userData.FirstName,
		LastName:     userData.LastName,
		LanguageCode: userData.LanguageCode,
	}
}
