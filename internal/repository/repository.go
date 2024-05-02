package repository

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/entity"
)

type RepositoryStruct struct {
	conn *pgx.Conn
}

func NewRepository() (repository *RepositoryStruct, err error) {
	databaseUrl := config.Config.DbURL
	conn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		return
	}
	repository = &RepositoryStruct{
		conn: conn,
	}
	return
}

func (repository *RepositoryStruct) GetCatalog() (catalog []entity.Catalog, err error) {
	catalog = make([]entity.Catalog, 0)
	rows, err := repository.conn.Query(context.Background(), `SELECT * FROM services`)
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

func (repository *RepositoryStruct) Close() (err error) {
	err = repository.conn.Close(context.Background())
	return
}

func (repository *RepositoryStruct) GetCommands() (commands []entity.Command, err error) {
	commands = make([]entity.Command, 0)
	rows, err := repository.conn.Query(context.Background(), `SELECT id, command, description FROM commands WHERE NOT deleted`)
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
