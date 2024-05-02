package config

import "os"

type ConfigStruct struct {
	Token string
	DbURL string
}

var Config ConfigStruct

func NewConfig() {
	Config = ConfigStruct{
		Token: os.Getenv("TOKEN"),
		DbURL: os.Getenv("DATABASE_URL"),
	}
}
