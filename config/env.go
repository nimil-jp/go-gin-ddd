package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Env EnvType

type EnvType struct {
	App struct {
		Secret string
		URL    string
	}
	DB struct {
		Socket   string
		Host     string
		Port     uint
		User     string
		Password string
		Name     string
	}
	SMTP struct {
		Host     string
		Port     string
		User     string
		Password string
	}
	Mail struct {
		From string
		Name string
	}
}

func init() {
	dotenvPath := fmt.Sprintf("./%s.env", os.Getenv("GO_ENV"))
	if dotenvPathEnv := os.Getenv("DOTENV_PATH"); dotenvPathEnv != "" {
		dotenvPath = dotenvPathEnv
	}
	err := godotenv.Load(dotenvPath)
	if err != nil {
		panic(err)
	}

	err = envconfig.Process("", &Env)
	if err != nil {
		panic(err)
	}
}
