package handler

import (
	"fmt"
	"strconv"

	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/tamnook/autoschool_telegram_bot/internal/entity"
	"github.com/tamnook/autoschool_telegram_bot/internal/store"
)

type StoreInterface interface {
	GetCatalog() ([]entity.Catalog, error)
	// InsertChat() error
	GetCommands() ([]entity.Command, error)
	Close() error
}

type HandlerStruct struct {
	store StoreInterface
}

func NewHandler() (handler *HandlerStruct, err error) {
	store, err := store.NewStore()
	if err != nil {
		return
	}
	handler = &HandlerStruct{
		store: store,
	}
	return
}

func (handler *HandlerStruct) StartHandler(bot *telego.Bot, update telego.Update) {
	botCommands := make([]telego.BotCommand, 0)
	commands, err := handler.store.GetCommands()
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range commands {
		botCommands = append(botCommands, telego.BotCommand{Command: v.Command, Description: v.Description})
	}

	_ = bot.SetMyCommands(&telego.SetMyCommandsParams{
		Commands: botCommands,
		// Scope:    &telego.BotCommandScopeDefault{},
	})
	// if err != nil {
	// 	fmt.Println(err)
	// }
	_, _ = bot.SendMessage(tu.Message(
		tu.ID(update.Message.Chat.ID),
		"Добро пожаловать!",
	))
}

func (handler *HandlerStruct) CatalogHandler(bot *telego.Bot, update telego.Update) {
	catalog, err := handler.store.GetCatalog()
	if err != nil {
		fmt.Println(err)
	}
	var inlineKeyboard *telego.InlineKeyboardMarkup
	keyboardRows := make([][]telego.InlineKeyboardButton, 0)
	for _, v := range catalog {
		keyboardRows = append(keyboardRows, tu.InlineKeyboardRow(
			tu.InlineKeyboardButton(fmt.Sprintf("%v - %.f рублей", v.Name, v.Price)).WithCallbackData(strconv.FormatInt(v.Id, 10)),
		))
	}
	inlineKeyboard = tu.InlineKeyboard(keyboardRows...)
	message := tu.Message(
		update.Message.Chat.ChatID(),
		"Каталог",
	).WithReplyMarkup(inlineKeyboard)

	// Sending message
	_, _ = bot.SendMessage(message)
}

func (handler *HandlerStruct) Close() (err error) {
	err = handler.store.Close()
	return
}

// var handler HandlerStruct
