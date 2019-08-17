package app

import (
	"github.com/sirupsen/logrus"
	"tg-reposter/pkg/tgbot"
	"tg-reposter/pkg/tgclient"
)

var logger = logrus.WithField("logger", "app")

func Start() {
	conf, err := LoadConfigFile("config.yaml")
	if err != nil {
		logger.Fatalf("config load failed %+v", err)
	}

	client := prepareClient(conf)
	bot := prepareBot(conf)

	pipeline := NewPipeline(conf.FilterRegex, client, bot)
	err = pipeline.Start()

	if err != nil {
		logger.Fatalf("%+v", err)
	}
}

func prepareClient(conf *Config) *tgclient.Client {
	client := tgclient.NewBuilder().
		DeviceModel(conf.Client.DeviceModel).
		SystemVersion(conf.Client.SystemVersion).
		ApplicationVersion(conf.Client.ApplicationVersion).
		SystemLanguageCode(conf.Client.SystemLanguageCode).
		AuthPhone(conf.Client.Phone).
		ApiId(conf.Client.ApiId).
		ApiHash(conf.Client.ApiHash).
		Socks5Proxy(
			conf.Client.Proxy.Host,
			conf.Client.Proxy.Port,
			conf.Client.Proxy.Login,
			conf.Client.Proxy.Password,
		).
		DatabaseDirectory(conf.Client.DatabaseDirectory).
		FilesDirectory(conf.Client.FilesDirectory).
		CheckCode(conf.Client.CheckCode).
		Build()

	client.SetLogVerbosity(1)

	err := client.Authorize()
	if err != nil {
		logger.Fatalf("auth failed. %+v", err)
	}

	return client
}

func prepareBot(conf *Config) *tgbot.Bot {
	proxy := conf.Client.Proxy

	bot, err := tgbot.NewBuilder().
		Token(conf.Bot.Token).
		TimeoutSec(conf.Bot.Timeout).
		Socks5Proxy(proxy.Host, proxy.Port, proxy.Login, proxy.Password).
		Build()

	if err != nil {
		logger.Fatalf("bot build failed. %+v", err)
	}

	return bot
}
