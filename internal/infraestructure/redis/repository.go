package redis

import "github.com/renzonaitor/tweet-api/cmd/http/config"

type Repository struct {
	// add redis client
}

func NewRepository(cfg config.Config) *Repository {
	return &Repository{}
}
