package redismgr

import (
	"time"
)

type Config struct {
	URI         string
	Name        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
	PoolSize    int
	MinIdleConn int
	Db          int
	Username    string
	Password    string
}

func (entity *Config) manageConfig() {
	if entity.DialTimeout <= 0 {
		entity.DialTimeout = 3 * time.Second
	}
	if entity.ReadTimeout <= 0 {
		entity.ReadTimeout = 2 * time.Second
	}
	if entity.PoolSize <= 0 {
		entity.PoolSize = 20
	}
	if entity.MinIdleConn <= 0 {
		entity.MinIdleConn = 5
	}
}
