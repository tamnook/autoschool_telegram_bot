package cache

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/entity"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/repository"
)

type CacheMu interface {
	StartAutoCacheRefresh(ctx context.Context, interval time.Duration)
	GetFAQCache(id int64) entity.FAQStruct
	GetAllFAQCache() []entity.FAQStruct
}

type cacheMu struct {
	repo repository.Repository
}

func NewCacheMu(repo repository.Repository) CacheMu {
	return &cacheMu{
		repo: repo,
	}
}

var (
	faqCache  = make(map[int64]entity.FAQStruct) // Кэшируем вопросы и ответы
	cacheRWMu sync.RWMutex                       // Для безопасного доступа к кэшу
)

func (c *cacheMu) initFAQCache(ctx context.Context) {
	cacheRWMu.Lock()
	defer cacheRWMu.Unlock()

	faqStruct, err := c.repo.GetFAQQuestions(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Очистим кэш перед загрузкой
	faqCache = make(map[int64]entity.FAQStruct)

	for _, item := range faqStruct {
		// Добавляем в кэш
		faqCache[item.Id] = item
	}
}

func (c *cacheMu) GetFAQCache(id int64) entity.FAQStruct {
	return faqCache[id]
}
func (c *cacheMu) GetAllFAQCache() []entity.FAQStruct {
	return lo.Values(faqCache)
}
func (c *cacheMu) StartAutoCacheRefresh(ctx context.Context, interval time.Duration) {
	c.initFAQCache(ctx)
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.initFAQCache(ctx)
		}
	}()
}
