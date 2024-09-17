package mongo

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type ReadPreference string

const (
	Master ReadPreference = "master"
	Slave  ReadPreference = "slave"
)

type Config struct {
	Client             string
	AppName            string
	URI                string
	DbName             string
	Username           string
	Password           string
	MasterMinPoolSize  uint64
	MasterMaxPoolSize  uint64
	SlaveEnabled       bool
	SlaveMinPoolSize   uint64
	SlaveMaxPoolSize   uint64
	SetRetryWrites     bool
	SetRetryReads      bool
	ConnectionTimeout  time.Duration
	MaxConnIdleTime    time.Duration
	HeartbeatInterval  time.Duration
	SocketTimeout      time.Duration
	MaxStalenessSecond time.Duration
}

func (entity *Config) manageConfig() {
	if entity.MasterMinPoolSize <= 0 {
		entity.MasterMinPoolSize = 10
	}
	if entity.MasterMaxPoolSize <= 0 {
		entity.MasterMaxPoolSize = 40
	}
	if entity.SlaveMinPoolSize <= 0 {
		entity.SlaveMinPoolSize = 10
	}
	if entity.SlaveMaxPoolSize <= 0 {
		entity.SlaveMaxPoolSize = 40
	}
	if entity.ConnectionTimeout <= 0 {
		entity.ConnectionTimeout = 20 * time.Second
	}
	if entity.MaxConnIdleTime <= 0 {
		entity.MaxConnIdleTime = 5 * time.Minute
	}
	if entity.HeartbeatInterval <= 0 {
		entity.HeartbeatInterval = 5 * time.Second
	}
	if entity.SocketTimeout <= 0 {
		entity.SocketTimeout = 10 * time.Second
	}
	if entity.MaxStalenessSecond <= 0 {
		entity.MaxStalenessSecond = 120 * time.Second
	}
}

func (entity *Config) getPrimaryClientOptions() *options.ClientOptions {
	mongoUri := fmt.Sprintf("mongodb://%s/?w=majority", entity.URI)
	appName := fmt.Sprintf(entity.AppName, "_master")
	if entity.Client == "atlas" {
		mongoUri = fmt.Sprintf(
			"mongodb+srv://%s/?w=majority",
			entity.URI,
		)
	}

	clientOptions := options.Client().ApplyURI(
		mongoUri,
	)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetAppName(appName)
	clientOptions.SetMinPoolSize(entity.MasterMinPoolSize)
	clientOptions.SetMaxPoolSize(entity.MasterMaxPoolSize)
	clientOptions.SetRetryWrites(entity.SetRetryWrites)
	clientOptions.SetRetryReads(entity.SetRetryReads)
	clientOptions.SetMaxConnIdleTime(entity.MaxConnIdleTime)
	clientOptions.SetConnectTimeout(entity.ConnectionTimeout)
	clientOptions.SetHeartbeatInterval(entity.HeartbeatInterval)
	clientOptions.SetSocketTimeout(entity.SocketTimeout)
	if entity.Username != "" && entity.Password != "" {
		credential := options.Credential{
			Username: entity.Username,
			Password: entity.Password,
		}
		clientOptions.SetAuth(credential)
	}

	clientOptions.SetReadPreference(readpref.Primary())

	return clientOptions
}

func (entity *Config) getSecondaryClientOptions() *options.ClientOptions {
	mongoUri := fmt.Sprintf("mongodb://%s/?w=majority", entity.URI)
	appName := fmt.Sprintf(entity.AppName, "_slave")
	if entity.Client == "atlas" {
		mongoUri = fmt.Sprintf(
			"mongodb+srv://%s/?w=majority",
			entity.URI,
		)
	}

	clientOptions := options.Client().ApplyURI(
		mongoUri,
	)
	clientOptions.SetConnectTimeout(10 * time.Second)
	clientOptions.SetAppName(appName)
	clientOptions.SetMinPoolSize(entity.SlaveMinPoolSize)
	clientOptions.SetMaxPoolSize(entity.SlaveMaxPoolSize)
	clientOptions.SetRetryWrites(entity.SetRetryWrites)
	clientOptions.SetRetryReads(entity.SetRetryReads)
	clientOptions.SetMaxConnIdleTime(entity.MaxConnIdleTime)
	clientOptions.SetConnectTimeout(entity.ConnectionTimeout)
	clientOptions.SetHeartbeatInterval(entity.HeartbeatInterval)
	clientOptions.SetSocketTimeout(entity.SocketTimeout)
	if entity.Username != "" && entity.Password != "" {
		credential := options.Credential{
			Username: entity.Username,
			Password: entity.Password,
		}
		clientOptions.SetAuth(credential)
	}

	clientOptions.SetReadPreference(
		readpref.SecondaryPreferred(
			readpref.WithMaxStaleness(entity.MaxStalenessSecond),
		),
	)
	return clientOptions
}
