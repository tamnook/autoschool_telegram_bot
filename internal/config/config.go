package config

import (
	"os"
	"time"
)

var (
	Token                string
	DbURL                string
	WebHookHost          string
	HandlersTimeout      time.Duration
	CacheRefreshDuration time.Duration
)

func InitConfig() error {
	Token = os.Getenv("TOKEN")
	DbURL = os.Getenv("DATABASE_URL")
	WebHookHost = os.Getenv("WEB_HOOK_HOST")
	timeoutString := os.Getenv("HANDLERS_TIMEOUT")
	var err error
	HandlersTimeout, err = time.ParseDuration(timeoutString)
	if err != nil {
		return err
	}
	cacheRefreshDuration := os.Getenv("CACHE_REFRESH_DURATION")
	CacheRefreshDuration, err = time.ParseDuration(cacheRefreshDuration)
	return err
}
