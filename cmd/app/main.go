package main

import (
	"context"
	"github.com/binance-converter/telegram-bot/internal/transport/bot_server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/sirupsen/logrus"
	"log"
)

type appConfig struct {
	Grpc struct {
		Port *int
	}
	Bot struct {
		Token string `env:"BINANCE_CONVERTER_BOOT_TOKEN"`
	}
}

func main() {
	setupLogs()
	cfg, err := initConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	botServer := bot_server.NewConverterBot(bot)

	ctx := context.Background()

	botServer.Start(ctx)
}

func setupLogs() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006.01.02 15:04:05",
		FullTimestamp:   true,
		DisableSorting:  true,
	})
}

func initConfig() (appConfig, error) {
	var cfg appConfig

	yamlFeeder := feeder.Yaml{Path: "config.yaml"}
	envFeeder := feeder.Env{}
	dotEnvFeeder := feeder.DotEnv{Path: ".env"}

	err := config.New().AddFeeder(yamlFeeder, envFeeder, dotEnvFeeder).AddStruct(&cfg).Feed()

	return cfg, err
}