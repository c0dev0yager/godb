package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClientStruct struct {
	Config Config
	Master *mongo.Database
	Slave  *mongo.Database
}

func (client *mongoClientStruct) createClient(
	options *options.ClientOptions,
) *mongo.Database {
	mongoConnection, err := mongo.Connect(context.TODO(), options)
	if err != nil {
		panic(err)
	}
	err = mongoConnection.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
	return mongoConnection.Database(client.Config.DbName)
}

var client map[string]*mongoClientStruct

func NewClient(
	config Config,
) {
	config.manageConfig()
	mongoClient := &mongoClientStruct{
		Config: config,
	}

	mongoClient.Master = mongoClient.createClient(
		config.getPrimaryClientOptions(),
	)
	// logs.GetClient().Infof("MongoMasterSuccessful to uri : %s and database %s", config.URI, config.DbName)

	if config.SlaveEnabled {
		mongoClient.Slave = mongoClient.createClient(
			config.getSecondaryClientOptions(),
		)
		// logs.GetClient().Infof("MongoSlaveSuccessful to uri : %s and database %s", config.URI, config.DbName)
	}
	client[config.DbName] = mongoClient
}

func GetMaster(dbName string) *mongo.Database {
	if client == nil {
		panic("NoMongoClient")
		return nil
	}
	if client[dbName] == nil {
		panic("NoMongoClientForDbName:" + dbName)
		return nil
	}
	return client[dbName].Master
}

func GetSlave(dbName string) *mongo.Database {
	if client == nil {
		panic("NoMongoClient")
	}
	if client[dbName] == nil {
		panic("NoMongoClientForDbName-" + dbName)
	}
	if !client[dbName].Config.SlaveEnabled {
		return client[dbName].Master
	}
	return client[dbName].Slave
}
