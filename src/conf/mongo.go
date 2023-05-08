package conf

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	// LogColName Collection name for transfers events
	LogColName = "Logs"
)

var (
	MongoCl *mongo.Client
	MongoDB *mongo.Database
)

func ConnectMongo() error {

	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	opts := options.Client().ApplyURI(Config.MONGO_Url)

	var err error
	MongoCl, err = mongo.Connect(ctx, opts)
	if err != nil {
		//log.Fatalf("Mongo: %s", err)
		return err
	}

	err = MongoCl.Ping(ctx, nil)
	//if err != nil {
	//	log.Fatalf("Mongo: %s", err)
	//}
	MongoDB = MongoCl.Database(Config.MONGO_DBName)
	return err
}
