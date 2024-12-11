package services

import "github.com/itaraxa/turbo-waddle/internal/config"

type Storage struct{}

func NewStorage(c *config.GopherMartConfig) (*Storage, error) {
	return &Storage{}, nil
}
