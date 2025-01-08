package bot

import (
	"context"
	"errors"

	"github.com/fasthttp/router"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/repository"
	"github.com/valyala/fasthttp"
)

const (
	pathPrefix = "/bot"
)

type Bot interface {
	Start(ctx context.Context)
}

type bot struct {
	telebot *telego.Bot
	repo    repository.Repository
	updates <-chan telego.Update
	bh      *th.BotHandler
}

func NewBot(ctx context.Context, telebot *telego.Bot, server *fasthttp.Server, repo repository.Repository) (Bot, error) {
	bot := &bot{
		telebot: telebot,
		repo:    repo,
	}
	srv := telego.FuncWebhookServer{
		Server: telego.FastHTTPWebhookServer{
			Logger: bot.telebot.Logger(),
			Server: server,
			Router: router.New(),
		},
	}
	bot.updates, _ = bot.telebot.UpdatesViaWebhook(pathPrefix+bot.telebot.Token(),
		telego.WithWebhookServer(srv),
		telego.WithWebhookSet(&telego.SetWebhookParams{
			URL: config.WebHookHost + pathPrefix + bot.telebot.Token(),
		}),
	)
	var err error
	bot.bh, err = th.NewBotHandler(bot.telebot, bot.updates)
	if err != nil {
		return nil, errors.New("error th.NewBotHandler")
	}

	bot.bh.Handle(func(b *telego.Bot, update telego.Update) {
		bot.startHandler(ctx, b, update)
	}, th.CommandEqual("start"))

	return bot, nil
}

func (bot *bot) Start(ctx context.Context) {

	go func() {
		_ = bot.telebot.StartWebhook("localhost:443")
	}()

	go func() {
		<-ctx.Done()
		bot.bh.Stop()
		bot.telebot.StopWebhook()
	}()

	bot.bh.Start()
}
