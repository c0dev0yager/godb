package redismgr

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Connection struct {
	DB *redis.Client
}

func NewConnection(
	config Config,
) (*Connection, error) {
	config.manageConfig()
	options := redis.Options{
		Addr:         config.URI,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConn,
		DB:           config.Db,
		Username:     config.Username,
		Password:     config.Password,
	}
	client := redis.NewClient(&options)
	ctx := context.Background()

	result := client.Ping(ctx)

	if result.Err() != nil {
		return nil, result.Err()
	}
	conn := &Connection{
		DB: client,
	}
	return conn, nil
}
