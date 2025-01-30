package bot

import (
	"context"
	"errors"

	"github.com/fasthttp/router"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/cache"
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
	cache   cache.CacheMu
}

func NewBot(ctx context.Context, telebot *telego.Bot, server *fasthttp.Server, repo repository.Repository, cache cache.CacheMu) (Bot, error) {
	bot := &bot{
		telebot: telebot,
		repo:    repo,
		cache:   cache,
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

	return bot, nil
}
func (bot *bot) Start(ctx context.Context) {

	go func() {
		//_ = bot.telebot.StartWebhook("192.168.1.17:443")
		_ = bot.telebot.StartWebhook("localhost:443")
	}()

	go func() {
		<-ctx.Done()
		bot.bh.Stop()
		bot.telebot.StopWebhook()
	}()

	for update := range bot.updates {
		if update.Message != nil {
			bot.handleMessage(ctx, bot.telebot, update)
		} else if update.CallbackQuery != nil {
			bot.handleCallback(ctx, bot.telebot, update.CallbackQuery)
		}
	}

	bot.bh.Start()
}
