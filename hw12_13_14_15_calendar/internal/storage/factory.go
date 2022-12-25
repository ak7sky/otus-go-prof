package storage

import (
	"context"

	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/app"
	"github.com/ak7sky/otus-go-prof/hw12_13_14_15_calendar/internal/config"
)

const (
	memory = "memory"
	rdb    = "rdb"
)

func New(ctx context.Context, storageConf config.StorageConf) (app.Storage, error) {
	switch storageConf.Type {
	case memory:
		return newMemStorage(ctx)
	case rdb:
		return newRdbStorage(ctx, storageConf.DSN)
	default:
		return newMemStorage(ctx)
	}
}
