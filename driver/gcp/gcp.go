package gcp

import (
	"io/ioutil"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"

	"go-gin-ddd/config"
)

var (
	conf *jwt.Config
)

func init() {
	credBytes, err := ioutil.ReadFile(config.Env.Gcp.CredentialPath)
	if err != nil {
		panic(err)
	}

	conf, err = google.JWTConfigFromJSON(credBytes)
	if err != nil {
		panic(err)
	}
}

func Conf() *jwt.Config {
	return conf
}
