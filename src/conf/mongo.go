package conf

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoCollectionName string

const (
	// LogCol Collection name for transfers events
	UserBalColName MongoCollectionName = "UsersBalance"
	UsersListCol   MongoCollectionName = "User"
)

var (
	MongoCl *mongo.Client
	MongoDB *mongo.Database
)

func ConnectMongo() {
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	opts := options.Client().ApplyURI(Config.MongoUrl.String())

	var err error
	MongoCl, err = mongo.Connect(ctx, opts)
	if err != nil {
		log.Panicf("Mongo: %s", err)
	}

	err = MongoCl.Ping(ctx, nil)
	//if err != nil {
	//	log.Fatalf("Mongo: %s", err)
	//}
	MongoDB = MongoCl.Database(Config.MongoDBName)
}

func GetMongoCol(chain int64, colName string) *mongo.Collection {
	return MongoCl.Database(fmt.Sprintf("%s_%d", Config.MongoDBName, chain)).Collection(colName)
}
