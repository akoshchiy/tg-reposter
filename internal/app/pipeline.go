package app

import (
	"github.com/sirupsen/logrus"
	"regexp"
	"tg-reposter/pkg/tgbot"
	"tg-reposter/pkg/tgclient"
)

type Pipeline struct {
	logger *logrus.Entry
	client *tgclient.Client
	bot    *tgbot.Bot
	re     *regexp.Regexp
}

func NewPipeline(regexMatch string, client *tgclient.Client, bot *tgbot.Bot) *Pipeline {
	re := regexp.MustCompile(regexMatch)
	return &Pipeline{
		client: client,
		bot:    bot,
		re:     re,
		logger: logrus.WithField("logger", "pipeline"),
	}
}

func (p *Pipeline) Start() error {
	bot, err := p.bot.GetMe()
	me, err := p.client.GetMe()

	if err != nil {
		return err
	}

	logger.Info("start listening messages")

	for msg := range p.client.ListenNewMessages() {
		text, ok, err := p.filterMessage(bot.Id, msg)
		if err != nil {
			logger.Errorf("message filter failed. msg: %s. %+v", msg, err)
		}
		if ok {
			err = p.bot.SendMessage(int64(me.Id), text)
			if err != nil {
				logger.Errorf("message repost failed. msg: %s. %+v", msg, err)
			}
			logger.Info("message repost ", msg)
		}
	}

	return nil
}

func (p *Pipeline) filterMessage(botId int32, msg tgclient.Message) (txt string, ok bool, err error) {
	if msg.SenderUserId == botId {
		return
	}
	//if msg.IsOutgoing {
	//	return
	//}
	classType, err := msg.GetContentType()
	if err != nil {
		return
	}
	if classType != tgclient.MessageTextType {
		return
	}
	msgText := tgclient.MessageText{}
	err = msg.UnmarshalContent(&msgText)
	if err != nil {
		return
	}
	txt = msgText.Text.Text
	ok = p.re.MatchString(msgText.Text.Text)
	return
}
