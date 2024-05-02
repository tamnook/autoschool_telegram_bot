package store

import (
	"github.com/tamnook/autoschool_telegram_bot/internal/entity"
	"github.com/tamnook/autoschool_telegram_bot/internal/repository"
)

type RepositoryInterface interface {
	GetCatalog() ([]entity.Catalog, error)
	Close() error
	GetCommands() ([]entity.Command, error)
}
type StoreStruct struct {
	repository RepositoryInterface
}

func NewStore() (store *StoreStruct, err error) {
	repository, err := repository.NewRepository()
	if err != nil {
		return
	}
	store = &StoreStruct{
		repository: repository,
	}
	return
}

func (store *StoreStruct) GetCatalog() (catalog []entity.Catalog, err error) {
	catalog, err = store.repository.GetCatalog()
	return
}

func (store *StoreStruct) Close() (err error) {
	err = store.repository.Close()
	return
}

func (store *StoreStruct) GetCommands() (commands []entity.Command, err error) {
	commands, err = store.repository.GetCommands()
	return
}
