package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"
	"github.com/mymmrac/telego"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/bot"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/cache"
	"github.com/tamnook/autoschool_telegram_bot/internal/pkg/repository"
	"github.com/valyala/fasthttp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	done := make(chan struct{}, 1)
	go func() {
		<-sigs
		fmt.Println("Stopping...")
		time.Sleep(2 * time.Second)
		cancel()
		done <- struct{}{}
	}()

	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}
	err := config.InitConfig()
	if err != nil {
		log.Fatalf("error init config: %v", err)
	}

	telebot, err := telego.NewBot(config.Token, telego.WithDefaultDebugLogger())
	if err != nil {
		fmt.Printf("Ошибка создания бота: %s", err)
	}

	repo, err := repository.NewRepository(ctx)
	if err != nil {
		log.Fatalf("error creating repository: %v", err)
	}

	cacheMu := cache.NewCacheMu(repo)
	cacheMu.StartAutoCacheRefresh(ctx, config.CacheRefreshDuration)

	server := &fasthttp.Server{}

	b, err := bot.NewBot(ctx, telebot, server, repo, cacheMu)
	if err != nil {
		log.Fatalf("error creating bot: %v", err)
	}

	b.Start(ctx)

	select {
	case <-ctx.Done():
		fmt.Printf("context error: %v\n", ctx.Err())
	case <-done:
		fmt.Printf("done")
	}
}
