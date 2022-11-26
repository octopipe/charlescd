package database

import "context"

type Model struct {
}

type Database interface {
	Create(ctx context.Context, model interface{}) (interface{}, error)
}
