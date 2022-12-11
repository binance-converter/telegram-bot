package main

import (
	"context"
	"fmt"
	"github.com/binance-converter/telegram-bot/internal/service"
	"github.com/binance-converter/telegram-bot/internal/transport/bot_server"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type appConfig struct {
	Grpc struct {
		Host *string
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

	grpcTarget := fmt.Sprintf("%s:%d", *cfg.Grpc.Host, *cfg.Grpc.Port)

	conn, err := grpc.Dial(grpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"grpc_target": grpcTarget,
		}).Fatal("error create grpc connection")
	} else {
		logrus.WithFields(logrus.Fields{
			"grpc_target": grpcTarget,
		}).Info("grpc connection created")
	}

	authService := service.NewAuth(conn)

	botServer := bot_server.NewConverterBot(bot, authService)

	ctx := context.Background()

	if err := botServer.Start(ctx); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error start bot server")
	}
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
