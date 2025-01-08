package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/entity"
)

type Repository interface {
	GetCatalog(ctx context.Context) (catalog []entity.Catalog, err error)
	GetCommands(ctx context.Context) (commands []entity.Command, err error)
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
		return catalog, err
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
