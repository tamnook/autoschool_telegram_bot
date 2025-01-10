package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/entity"
)

type Repository interface {
	GetCatalog(ctx context.Context) (catalog []entity.Catalog, err error)
	GetCommands(ctx context.Context) (commands []entity.Command, err error)
	GetFAQQuestions(ctx context.Context) (commands []entity.FAQStruct, err error)
	SaveFullNameStudent(ctx context.Context, student entity.Student) (err error)
	CreateStudent(ctx context.Context, student entity.Student) (err error)
	SavePhoneStudent(ctx context.Context, student entity.Student) (err error)
}

type repositoryStruct struct {
	conn *pgx.Conn
}

func NewRepository(ctx context.Context) (repository *repositoryStruct, err error) {
	databaseUrl := config.DbURL
	conn, err := pgx.Connect(ctx, databaseUrl)
	if err != nil {
		return
	}

	go func() {
		for range ctx.Done() {
			_ = conn.Close(ctx)

		}
	}()

	repository = &repositoryStruct{
		conn: conn,
	}
	return
}

func (repository *repositoryStruct) GetCatalog(ctx context.Context) (catalog []entity.Catalog, err error) {
	catalog = make([]entity.Catalog, 0)
	rows, err := repository.conn.Query(ctx, `SELECT * FROM services`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.Catalog
		err = rows.Scan(
			&item.Name,
			&item.Price,
			&item.Id,
		)
		if err != nil {
			return catalog, err
		}
		catalog = append(catalog, item)
	}
	return
}

func (repository *repositoryStruct) GetCommands(ctx context.Context) (commands []entity.Command, err error) {
	commands = make([]entity.Command, 0)
	rows, err := repository.conn.Query(ctx, `SELECT id, command, description FROM commands WHERE NOT deleted`)
	if err != nil {
		return commands, err
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.Command
		err = rows.Scan(
			&item.Id,
			&item.Command,
			&item.Description,
		)
		if err != nil {
			return commands, err
		}
		commands = append(commands, item)
	}
	return
}

func (repository *repositoryStruct) CreateStudent(ctx context.Context, student entity.Student) (err error) {
	fmt.Printf("%v, %v, %v", student.ID, student.TelegramChatID, student.TelegramUserName)
	_, err = repository.conn.Exec(ctx, "INSERT INTO students (id, telegram_chat_id, telegram_user_name) VALUES ($1, $2, $3) ON CONFLICT (telegram_chat_id) DO NOTHING", student.ID, student.TelegramChatID, student.TelegramUserName)
	return err
}
func (repository *repositoryStruct) SaveFullNameStudent(ctx context.Context, student entity.Student) (err error) {
	_, err = repository.conn.Exec(ctx, "INSERT INTO students (id, full_name, telegram_chat_id, telegram_user_name) VALUES ($1, $2, $3, $4) ON CONFLICT (telegram_chat_id) DO UPDATE SET full_name = EXCLUDED.full_name, telegram_user_name = EXCLUDED.telegram_user_name", student.ID, student.FullName, student.TelegramChatID, student.TelegramUserName)
	return err
}
func (repository *repositoryStruct) SavePhoneStudent(ctx context.Context, student entity.Student) (err error) {
	_, err = repository.conn.Exec(ctx, "INSERT INTO students (id, phone, telegram_chat_id, telegram_user_name) VALUES ($1, $2, $3, $4) ON CONFLICT (telegram_chat_id) DO UPDATE SET phone = EXCLUDED.phone, telegram_user_name = EXCLUDED.telegram_user_name", student.ID, student.Phone, student.TelegramChatID, student.TelegramUserName)
	return err
}

func (repository *repositoryStruct) GetFAQQuestions(ctx context.Context) (commands []entity.FAQStruct, err error) {
	rows, err := repository.conn.Query(ctx, "SELECT id, question, answer FROM faq")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item entity.FAQStruct
		err = rows.Scan(
			&item.Id,
			&item.Question,
			&item.Answer,
		)
		if err != nil {
			return commands, err
		}
		commands = append(commands, item)
	}
	return
}
