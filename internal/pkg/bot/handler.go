package bot

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegoutil"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/entity"
)

var (
	expectingFullName bool
	expectingPhone    bool
	userData          = make(map[int64]string)
)

func MainMenu() *telego.ReplyKeyboardMarkup {
	return &telego.ReplyKeyboardMarkup{
		Keyboard: [][]telego.KeyboardButton{
			{{Text: "📝 Регистрация"}, {Text: "📚 Часто задаваемые вопросы"}},
			{{Text: "📞 Связь с менеджером"}, {Text: "🔔 Уведомления"}},
		},
		ResizeKeyboard: true,
	}
}
func startHandler(_ context.Context, bot *telego.Bot, update telego.Update) {
	_, _ = bot.SendMessage(telegoutil.MessageWithEntities(
		telegoutil.ID(update.Message.Chat.ID),
		telegoutil.Entity("Добро пожаловать в автошколу🚘 \nВыберите действие из меню:"),
	).WithReplyMarkup(MainMenu()))
}

func (b *bot) validateInput(ctx context.Context, bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ChatID()
	text := update.Message.Text
	ID := uuid.New()
	fmt.Printf("expectingFullName: %v", update.Message.Chat.Username)
	if expectingFullName {
		if isValidFullName(text) {
			expectingFullName = false
			expectingPhone = true
			err := b.repo.SaveFullNameStudent(ctx, entity.Student{ID: ID, FullName: text, TelegramChatID: chatID.ID, TelegramUserName: update.Message.Chat.Username})
			if err != nil {
				fmt.Printf("Ошибка сохранения студента: %v", err)
			}
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "Введите номер телефона (в формате +7ХХХХХХХХХХ):"))
		} else {
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "❌ ФИО введено некорректно. Попробуйте еще раз."))
			expectingFullName = true
		}
	} else if expectingPhone {
		if isValidPhoneNumber(text) {
			expectingPhone = false
			err := b.repo.SavePhoneStudent(ctx, entity.Student{ID: ID, Phone: text, TelegramChatID: chatID.ID, TelegramUserName: update.Message.Chat.Username})
			if err != nil {
				fmt.Printf("Ошибка сохранения студента: %v", err)
			}
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "✅ Регистрация завершена!"))
		} else {
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "❌ Номер телефона введен некорректно. Попробуйте еще раз."))
		}
	}

	// _, _ = bot.SendMessage(telegoutil.MessageWithEntities(
	// 	telegoutil.ID(update.Message.Chat.ID),
	// 	telegoutil.Entity("Введите ваше ФИО (Фамилия Имя Отчество):"),
	// ).WithReplyMarkup(MainMenu()))
}
func (b *bot) sendFAQMenuHandler(ctx context.Context, bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ChatID()
	faqStruct := b.cache.GetAllFAQCache()

	var buttons [][]telego.InlineKeyboardButton
	for _, item := range faqStruct {
		button := telegoutil.InlineKeyboardRow(telegoutil.InlineKeyboardButton(item.Question).WithCallbackData(fmt.Sprintf("faq_%d", item.Id)))
		buttons = append(buttons, button)
	}
	keyboard := telegoutil.InlineKeyboardGrid(buttons)

	msg := telegoutil.Message(chatID, "📌 *Часто задаваемые вопросы:*\nВыберите интересующий вопрос:").
		WithParseMode(telego.ModeMarkdown).
		WithReplyMarkup(keyboard)

	_, _ = bot.SendMessage(msg)
}
func (b *bot) handleMessage(ctx context.Context, bot *telego.Bot, update telego.Update) {
	switch update.Message.Text {
	case "/start":
		startHandler(ctx, bot, update)
	case "📚 Часто задаваемые вопросы":
		b.sendFAQMenuHandler(ctx, bot, update)
	// case "📞 Связаться с менеджером":
	// 	sendManagerContact(bot, update.Message.Chat.ID)
	// case "📩 Обратный звонок":
	// 	requestCallback(bot, update.Message.Chat.ID)
	// case "📅 Расписание":
	// 	sendSchedule(bot, update.Message.Chat.ID)
	// case "🔔 Уведомления":
	// 	sendNotifications(bot, update.Message.Chat.ID)
	case "📝 Регистрация":
		bot.SendMessage(telegoutil.MessageWithEntities(
			telegoutil.ID(update.Message.Chat.ID),
			telegoutil.Entity("Введите ваше ФИО (Фамилия Имя Отчество):")))
		b.validateInput(ctx, bot, update)
		expectingFullName = true
		//b.registrationHandler(ctx, bot, update)
	default:
		b.validateInput(ctx, bot, update)
		// bot.SendMessage(telegoutil.MessageWithEntities(
		// 	telegoutil.ID(update.Message.Chat.ID),
		// 	telegoutil.Entity("Неизвестная команда. Используйте меню.")))
	}
}

func (b *bot) handleCallback(_ context.Context, bot *telego.Bot, callbackQuery *telego.CallbackQuery) {
	fmt.Println(callbackQuery.Message.GetChat())
	chatID := telego.ChatID{ID: callbackQuery.Message.GetChat().ID, Username: callbackQuery.Message.GetChat().Username}
	if strings.HasPrefix(callbackQuery.Data, "faq_") {
		id, err := strconv.Atoi(strings.TrimPrefix(callbackQuery.Data, "faq_"))
		if err != nil {
			bot.SendMessage(telegoutil.Message(callbackQuery.Message.GetChat().PersonalChat.ChatID(), "Ошибка обработки запроса."))
			return
		}
		answer, err := b.cache.GetFAQCache(int64(id)).Answer, nil
		if err != nil {
			bot.SendMessage(telegoutil.Message(chatID, "Ошибка загрузки ответа."))
			return
		}

		fmt.Println("Ответ: " + answer)
		editMsg := telegoutil.Message(chatID, answer).
			WithParseMode(telego.ModeMarkdown)

		_, _ = bot.EditMessageText(&telego.EditMessageTextParams{
			ChatID:    chatID,
			MessageID: callbackQuery.Message.GetMessageID(),
			Text:      editMsg.Text,
			ParseMode: editMsg.ParseMode,
		})
	}
}

func isValidFullName(fullName string) bool {
	fullName = strings.TrimSpace(fullName)
	matched, _ := regexp.MatchString(`^[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+\s[А-ЯЁ][а-яё]+$`, fullName)
	return matched
}

func isValidPhoneNumber(phone string) bool {
	matched, _ := regexp.MatchString(`^\+7\d{10}$`, phone)
	return matched
}
