package gcp

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	"go-gin-ddd/config"
)

var (
	client *storage.Client
)

func init() {
	var err error

	ctx := context.Background()
	client, err = storage.NewClient(ctx, option.WithCredentialsFile(config.Env.Gcp.CredentialPath))
	if err != nil {
		panic(err)
	}
}

func GcsClient() *storage.Client {
	return client
}
