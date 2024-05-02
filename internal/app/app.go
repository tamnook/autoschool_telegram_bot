package app

import (
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/handler"
)

type HandlerInterface interface {
	StartHandler(bot *telego.Bot, update telego.Update)
	CatalogHandler(bot *telego.Bot, update telego.Update)
	Close() error
}

type AppStruct struct {
	handler HandlerInterface
}

func NewApp() (app *AppStruct, err error) {
	handler, err := handler.NewHandler()
	if err != nil {
		return
	}
	app = &AppStruct{
		handler: handler,
	}
	return
}

func (app *AppStruct) Start() (err error) {
	token := config.Config.Token
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		return
	}

	updates, _ := bot.UpdatesViaLongPolling(nil)
	defer bot.StopLongPolling()

	bh, _ := th.NewBotHandler(bot, updates)
	defer bh.Stop()

	bh.Handle(app.handler.StartHandler, th.CommandEqual("start"))
	// bh.Handle(KeyboardHandler, th.CommandEqual("keyboard"))
	bh.Handle(app.handler.CatalogHandler, th.CommandEqual("catalog"))
	// // bh.HandleCallbackQuery(CallbackHandler)

	bh.Start()
	defer app.handler.Close()
	return
}
