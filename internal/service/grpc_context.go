package service

import (
	"context"
	"google.golang.org/grpc/metadata"
	"strconv"
)

const (
	chatIdKey = "chat_id"
)

func addChatIdToContext(ctx context.Context, chatId int64) context.Context {
	return metadata.NewIncomingContext(ctx, metadata.New(map[string]string{
		chatIdKey: strconv.FormatInt(chatId, 10),
	}))
}
