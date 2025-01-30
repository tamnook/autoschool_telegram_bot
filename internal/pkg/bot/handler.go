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
	"github.com/samber/lo"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/entity"
)

var (
	student = entity.Student{}
)

func MainMenu() *telego.ReplyKeyboardMarkup {
	return &telego.ReplyKeyboardMarkup{
		Keyboard: [][]telego.KeyboardButton{
			{{Text: "📚 Часто задаваемые вопросы"}},
			{{Text: "📞 Связь с менеджером"}, {Text: "🔔 Уведомления"}},
		},
		ResizeKeyboard: true,
	}
}
func (b *bot) startHandler(ctx context.Context, bot *telego.Bot, update telego.Update) {
	b.cache.InitStudentsCache(ctx, update.Message.Chat.ID)
	student = b.cache.GetStudentCache(update.Message.Chat.ID)
	fmt.Println(student)
	if lo.FromPtr(student.FullName) == "" || lo.FromPtr(student.Phone) == "" {
		text := "Введите номер телефона (в формате +7ХХХХХХХХХХ):"
		if lo.FromPtr(student.FullName) == "" {
			text = "Введите ФИО в формате:\nФамилия Имя Отчество (все обязательно с большой буквы, без лишних символов)"
		}
		_, _ = bot.SendMessage(telegoutil.MessageWithEntities(
			telegoutil.ID(update.Message.Chat.ID),
			telegoutil.Entity("Добро пожаловать в автошколу🚘"+"\n"+text),
		))
	} else {
		_, _ = bot.SendMessage(telegoutil.MessageWithEntities(
			telegoutil.ID(update.Message.Chat.ID),
			telegoutil.Entity("Добро пожаловать в автошколу🚘 \nВыберите действие из меню:"),
		).WithReplyMarkup(MainMenu()))
	}
}

func (b *bot) validateInput(ctx context.Context, bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ChatID()
	text := update.Message.Text
	ID := uuid.New()
	if lo.FromPtr(student.FullName) == "" {
		if isValidFullName(text) {
			err := b.repo.SaveFullNameStudent(ctx, entity.Student{ID: ID, FullName: lo.ToPtr(text), TelegramChatID: chatID.ID, TelegramUserName: lo.ToPtr(update.Message.Chat.Username)})
			if err != nil {
				fmt.Printf("Ошибка сохранения студента: %v", err)
				return
			}
			student.FullName = lo.ToPtr(text)
			b.cache.SetStudentCache(student)
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "Введите номер телефона (в формате +7ХХХХХХХХХХ):"))
		} else {
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "❌ ФИО введено некорректно. Попробуйте еще раз."))
		}
	} else if lo.FromPtr(student.Phone) == "" {
		if isValidPhoneNumber(text) {
			err := b.repo.SavePhoneStudent(ctx, entity.Student{ID: ID, Phone: lo.ToPtr(text), TelegramChatID: chatID.ID, TelegramUserName: lo.ToPtr(update.Message.Chat.Username)})
			if err != nil {
				fmt.Printf("Ошибка сохранения студента: %v", err)
				return
			}
			student.Phone = lo.ToPtr(text)
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "✅ Регистрация завершена!"))
			b.cache.InitStudentsCache(ctx, update.Message.Chat.ID)
			b.startHandler(ctx, bot, update)
		} else {
			_, _ = bot.SendMessage(telegoutil.Message(chatID, "❌ Номер телефона введен некорректно. Попробуйте еще раз."))
		}
	} else {
		_, _ = bot.SendMessage(telegoutil.MessageWithEntities(
			telegoutil.ID(update.Message.Chat.ID),
			telegoutil.Entity("Неизвестная команда. Используйте меню."),
		).WithReplyMarkup(MainMenu()))
	}
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

// Связь с менеджером
func (b *bot) sendManagerContact(ctx context.Context, bot *telego.Bot, update telego.Update) {
	chatID := update.Message.Chat.ChatID()
	_, _ = bot.SendMessage(telegoutil.Message(chatID, "☎️ Связаться с менеджером: +7 900 123 45 67"))
}

func (b *bot) handleMessage(ctx context.Context, bot *telego.Bot, update telego.Update) {
	student = b.cache.GetStudentCache(update.Message.Chat.ID)
	if student == (entity.Student{}) {
		chatID := update.Message.Chat.ChatID()
		ID := uuid.New()
		err := b.repo.CreateStudent(ctx, entity.Student{ID: ID, TelegramChatID: chatID.ID, TelegramUserName: lo.ToPtr(update.Message.Chat.Username)})
		if err != nil {
			fmt.Printf("Ошибка создания студента: %v", err)
			return
		}
	}

	switch update.Message.Text {
	case "/start":
		b.startHandler(ctx, bot, update)
	case "📚 Часто задаваемые вопросы":
		b.sendFAQMenuHandler(ctx, bot, update)
	case "📞 Связаться с менеджером":
		b.sendManagerContact(ctx, bot, update)
	// case "📅 Расписание":
	// 	sendSchedule(bot, update.Message.Chat.ID)
	// case "🔔 Уведомления":
	// 	sendNotifications(bot, update.Message.Chat.ID)
	default:
		b.validateInput(ctx, bot, update)
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
