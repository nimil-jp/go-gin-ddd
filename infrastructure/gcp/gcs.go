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
	GetSignedUrl(dir string) (*SignedUrl, error)
	Delete(key string) error
}

type gcs struct{}

func NewGcs() IGcs {
	return gcs{}
}

type SignedUrl struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

func (gcs) GetSignedUrl(dir string) (*SignedUrl, error) {
	key, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	keyString := fmt.Sprintf("%s/%s", dir, key.String())

	url, err := storage.SignedURL(config.Env.Gcp.Bucket, keyString, &storage.SignedURLOptions{
		GoogleAccessID: gcp.Conf().Email,
		PrivateKey:     gcp.Conf().PrivateKey,
		Method:         http.MethodPut,
		Expires:        time.Now().Add(time.Minute * 1),
	})

	return &SignedUrl{
		Key: keyString,
		URL: url,
	}, nil
}

func (gcs) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	return gcp.GcsClient().Bucket(config.Env.Gcp.Bucket).Object(key).Delete(ctx)
}
