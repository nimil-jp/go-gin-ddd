package config

import (
	"encoding/json"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var Env EnvType

type EnvType struct {
	Port string `default:"8080"`
	App  struct {
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

func fileExists(path string) bool {
	if f, err := os.Stat(path); os.IsNotExist(err) || f.IsDir() {
		return false
	} else {
		return true
	}
}

func init() {
	dotenvPath := "./.env"
	if dotenvPathEnv := os.Getenv("DOTENV_PATH"); dotenvPathEnv != "" {
		dotenvPath = dotenvPathEnv
	}
	if fileExists(dotenvPath) {
		err := godotenv.Load(dotenvPath)
		if err != nil {
			panic(err)
		}
	}

	if env := os.Getenv("AWS_SECRET_ENV"); env != "" {
		var jsonEnv map[string]string
		err := json.Unmarshal([]byte(env), &jsonEnv)
		if err != nil {
			panic(err)
		}

		for k, v := range jsonEnv {
			err = os.Setenv(k, v)
			if err != nil {
				panic(err)
			}
		}
	}

	err := envconfig.Process("", &Env)
	if err != nil {
		panic(err)
	}
}
