package bot

import (
	"github.com/mymmrac/telego"
)

type BotStruct struct {
	bot telego.Bot
}

func NewBot(token string) *BotStruct {
	return &BotStruct{}
}

// func (bot *BotStruct) Start() (err error) {
// 	bot, err = telego.NewBot("7131272434:AAFDehT-ULcLWxA5VofSl3T7DNdKzHQDn8A", telego.WithDefaultDebugLogger())
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	return
// }
