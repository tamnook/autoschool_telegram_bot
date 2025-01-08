package config

import (
	"os"
	"time"
)

var (
	Token           string
	DbURL           string
	WebHookHost     string
	HandlersTimeout time.Duration
)

func InitConfig() error {
	Token = os.Getenv("TOKEN")
	DbURL = os.Getenv("DATABASE_URL")
	WebHookHost = os.Getenv("WEB_HOOK_HOST")
	timeoutString := os.Getenv("HANDLERS_TIMEOUT")
	var err error
	HandlersTimeout, err = time.ParseDuration(timeoutString)
	return err
}
