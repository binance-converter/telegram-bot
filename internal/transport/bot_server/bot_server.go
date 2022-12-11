package bot_server

import (
	"context"
	"github.com/binance-converter/telegram-bot/internal/transport/bot_server/bot_handler"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type ConverterBot struct {
	bot     *tgbotapi.BotAPI
	handler bot_handler.BotHandler
}

func NewConverterBot(bot *tgbotapi.BotAPI) *ConverterBot {
	return &ConverterBot{bot: bot}
}

func (c *ConverterBot) Start(ctx context.Context) error {
	logrus.WithFields(logrus.Fields{
		"user_name": c.bot.Self.UserName,
	}).Info("boot started")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := c.bot.GetUpdatesChan(u)

	c.updateHandler(ctx, updates)

	return nil
}

func (c *ConverterBot) updateHandler(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // If we got a message
			c.msgHandler(ctx, update)
		}
	}
}

func (c *ConverterBot) msgHandler(ctx context.Context, update tgbotapi.Update) {
	logrus.WithFields(logrus.Fields{
		"update_id": update.UpdateID,
		"chat_id":   update.Message.Chat.ID,
		"username":  update.Message.Chat.UserName,
		"text":      update.Message.Text,
	}).Info("receive massage")

	var massage *tgbotapi.MessageConfig

	if update.Message.IsCommand() {
		var err error
		massage, err = c.handler.CmdHandler(ctx, update)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err.Error(),
				"update_id": update.UpdateID,
			}).Error("Error handle command")
		}
	} else {
		var err error
		massage, err = c.handler.AnswerHandler(ctx, update)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err.Error(),
				"update_id": update.UpdateID,
			}).Error("Error handle answer")
		}
	}

	if massage != nil {
		c.bot.Send(*massage)
	}
}
