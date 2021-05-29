package model

import (
	"go-ddd/domain"
)

type Task struct {
	domain.SoftDeleteModel
	Name string `json:"name"`
	Done bool   `json:"done"`
}
