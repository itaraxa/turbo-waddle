package services

import "context"

type userStorager interface{}

type logger interface{}

func Registration(ctx context.Context, l logger, us userStorager, login string, password string) error {

	return nil
}
