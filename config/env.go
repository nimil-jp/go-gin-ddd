package config

import (
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
	Smtp struct {
		Host     string
		Port     string
		User     string
		Password string
	}
	Mail struct {
		From string
		Name string
	}
	Gcp struct {
		CredentialPath string `split_words:"true"`
		Bucket         string
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

	err := envconfig.Process("", &Env)
	if err != nil {
		panic(err)
	}
}
