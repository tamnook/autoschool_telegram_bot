package bot

import (
	"context"
	"fmt"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
)

func (b *bot) startHandler(ctx context.Context, bot *telego.Bot, update telego.Update) {
	ctx, cancel := context.WithTimeout(ctx, config.HandlersTimeout)
	defer cancel()
	botCommands := make([]telego.BotCommand, 0)
	commands, err := b.repo.GetCommands(ctx)
	if err != nil {
		fmt.Println(err)
	}
	for _, v := range commands {
		botCommands = append(botCommands, telego.BotCommand{Command: v.Command, Description: v.Description})
	}
	bot.SetMyCommands(&telego.SetMyCommandsParams{
		Commands: botCommands,
	})
	_, _ = bot.SendMessage(telegoutil.MessageWithEntities(
		telegoutil.ID(update.Message.Chat.ID),
		telegoutil.Entity("Добро пожаловать!\nНапишите, пожалуйста, вашу почту или телефон.\nНапример, "),
		telegoutil.Entity("example@email.com").Bold().Email(),
		telegoutil.Entity(" или "),
		telegoutil.Entity("+79123456789").PhoneNumber().Bold(),
	))
}
