package gcp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/google/uuid"

	"go-gin-ddd/config"
	"go-gin-ddd/driver/gcp"
)

type IGcs interface {
	GetSignedURL(dir string, public bool) (*SignedURL, error)
	Delete(key string) error
}

type gcs struct{}

func NewGcs() IGcs {
	return gcs{}
}

type SignedURL struct {
	Key    string `json:"key"`
	URL    string `json:"url"`
	Public bool   `json:"public"`
}

func (gcs) GetSignedURL(dir string, public bool) (*SignedURL, error) {
	key, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	keyString := fmt.Sprintf("%s/%s", dir, key.String())

	var headers []string
	if public {
		headers = append(headers, "x-goog-acl:public-read")
	}

	url, err := storage.SignedURL(config.Env.Gcp.Bucket, keyString, &storage.SignedURLOptions{
		GoogleAccessID: gcp.Conf().Email,
		PrivateKey:     gcp.Conf().PrivateKey,
		Method:         http.MethodPut,
		Expires:        time.Now().Add(config.SignedURLDuration),
		Headers:        headers,
	})
	if err != nil {
		return nil, err
	}

	return &SignedURL{
		Key: keyString,
		URL: url,
	}, nil
}

func (gcs) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return gcp.GcsClient().Bucket(config.Env.Gcp.Bucket).Object(key).Delete(ctx)
}
