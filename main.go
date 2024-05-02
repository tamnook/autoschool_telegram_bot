package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/tamnook/autoschool_telegram_bot/internal/app"
	"github.com/tamnook/autoschool_telegram_bot/internal/config"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	config.NewConfig()
}

type appInterface interface {
	Start() error
}

var appI appInterface

func main() {
	var err error
	appI, err = app.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	err = appI.Start()
	if err != nil {
		log.Fatal(err)
	}
}
